package warden

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-kratos/kratos/pkg/net/rpc/warden/resolver"
	"github.com/go-kratos/kratos/pkg/net/rpc/warden/resolver/direct"

	"github.com/go-kratos/kratos/pkg/conf/env"
	"github.com/go-kratos/kratos/pkg/conf/flagvar"
	"github.com/go-kratos/kratos/pkg/ecode"
	"github.com/go-kratos/kratos/pkg/naming"
	nmd "github.com/go-kratos/kratos/pkg/net/metadata"
	"github.com/go-kratos/kratos/pkg/net/netutil/breaker"
	"github.com/go-kratos/kratos/pkg/net/rpc/warden/balancer/p2c"
	"github.com/go-kratos/kratos/pkg/net/rpc/warden/internal/status"
	"github.com/go-kratos/kratos/pkg/net/trace"
	xtime "github.com/go-kratos/kratos/pkg/time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	gstatus "google.golang.org/grpc/status"
)

var _grpcTarget flagvar.StringVars

var (
	_once           sync.Once
	_defaultCliConf = &ClientConfig{
		Dial:              xtime.Duration(time.Second * 10),
		Timeout:           xtime.Duration(time.Millisecond * 250),
		Subset:            50,
		KeepAliveInterval: xtime.Duration(time.Second * 60),
		KeepAliveTimeout:  xtime.Duration(time.Second * 20),
	}
	_defaultClient *Client
)

func baseMetadata() metadata.MD {
	gmd := metadata.MD{nmd.Caller: []string{env.AppID}}
	if env.Color != "" {
		gmd[nmd.Color] = []string{env.Color}
	}
	return gmd
}

// Register direct resolver by default to handle direct:// scheme.
func init() {
	resolver.Register(direct.New())
}

// ClientConfig is rpc client conf.
type ClientConfig struct {
	Dial                   xtime.Duration
	Timeout                xtime.Duration
	Breaker                *breaker.Config
	Method                 map[string]*ClientConfig
	Clusters               []string
	Zone                   string
	Subset                 int
	NonBlock               bool
	KeepAliveInterval      xtime.Duration
	KeepAliveTimeout       xtime.Duration
	KeepAliveWithoutStream bool
}

// Client is the framework's client side instance, it contains the ctx, opt and interceptors.
// Create an instance of Client, by using NewClient().
type Client struct {
	conf    *ClientConfig
	breaker *breaker.Group
	mutex   sync.RWMutex

	opts     []grpc.DialOption
	handlers []grpc.UnaryClientInterceptor
}

// TimeoutCallOption timeout option.
type TimeoutCallOption struct {
	*grpc.EmptyCallOption
	Timeout time.Duration
}

// WithTimeoutCallOption can override the timeout in ctx and the timeout in the configuration file
func WithTimeoutCallOption(timeout time.Duration) *TimeoutCallOption {
	return &TimeoutCallOption{&grpc.EmptyCallOption{}, timeout}
}

// handle returns a new unary client interceptor for OpenTracing\Logging\LinkTimeout.
func (c *Client) handle() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		var (
			ok     bool
			t      trace.Trace
			gmd    metadata.MD
			conf   *ClientConfig
			cancel context.CancelFunc
			addr   string
			p      peer.Peer
		)
		var ec ecode.Codes = ecode.OK
		// apm tracing
		if t, ok = trace.FromContext(ctx); ok {
			t = t.Fork("", method)
			defer t.Finish(&err)
		}

		// setup metadata
		gmd = baseMetadata()
		trace.Inject(t, trace.GRPCFormat, gmd)
		c.mutex.RLock()
		if conf, ok = c.conf.Method[method]; !ok {
			conf = c.conf
		}
		c.mutex.RUnlock()
		brk := c.breaker.Get(method)
		if err = brk.Allow(); err != nil {
			_metricClientReqCodeTotal.Inc(method, "breaker")
			return
		}
		defer onBreaker(brk, &err)
		var timeOpt *TimeoutCallOption
		for _, opt := range opts {
			var tok bool
			timeOpt, tok = opt.(*TimeoutCallOption)
			if tok {
				break
			}
		}
		if timeOpt != nil && timeOpt.Timeout > 0 {
			ctx, cancel = context.WithTimeout(nmd.WithContext(ctx), timeOpt.Timeout)
		} else {
			_, ctx, cancel = conf.Timeout.Shrink(ctx)
		}

		defer cancel()
		nmd.Range(ctx,
			func(key string, value interface{}) {
				if valstr, ok := value.(string); ok {
					gmd[key] = []string{valstr}
				}
			},
			nmd.IsOutgoingKey)
		// merge with old matadata if exists
		if oldmd, ok := metadata.FromOutgoingContext(ctx); ok {
			gmd = metadata.Join(gmd, oldmd)
		}
		ctx = metadata.NewOutgoingContext(ctx, gmd)

		opts = append(opts, grpc.Peer(&p))
		if err = invoker(ctx, method, req, reply, cc, opts...); err != nil {
			gst, _ := gstatus.FromError(err)
			ec = status.ToEcode(gst)
			err = errors.WithMessage(ec, gst.Message())
		}
		if p.Addr != nil {
			addr = p.Addr.String()
		}
		if t != nil {
			t.SetTag(trace.String(trace.TagAddress, addr), trace.String(trace.TagComment, ""))
		}
		return
	}
}

func onBreaker(breaker breaker.Breaker, err *error) {
	if err != nil && *err != nil {
		if ecode.EqualError(ecode.ServerErr, *err) || ecode.EqualError(ecode.ServiceUnavailable, *err) || ecode.EqualError(ecode.Deadline, *err) || ecode.EqualError(ecode.LimitExceed, *err) {
			breaker.MarkFailed()
			return

		}
	}
	breaker.MarkSuccess()
}

// NewConn will create a grpc conn by default config.
func NewConn(target string, opt ...grpc.DialOption) (*grpc.ClientConn, error) {
	return DefaultClient().Dial(context.Background(), target, opt...)
}

// NewClient returns a new blank Client instance with a default client interceptor.
// opt can be used to add grpc dial options.
func NewClient(conf *ClientConfig, opt ...grpc.DialOption) *Client {
	c := new(Client)
	if err := c.SetConfig(conf); err != nil {
		panic(err)
	}
	c.UseOpt(grpc.WithBalancerName(p2c.Name))
	c.UseOpt(opt...)
	return c
}

// DefaultClient returns a new default Client instance with a default client interceptor and default dialoption.
// opt can be used to add grpc dial options.
func DefaultClient() *Client {
	_once.Do(func() {
		_defaultClient = NewClient(nil)
	})
	return _defaultClient
}

// SetConfig hot reloads client config
func (c *Client) SetConfig(conf *ClientConfig) (err error) {
	if conf == nil {
		conf = _defaultCliConf
	}
	if conf.Dial <= 0 {
		conf.Dial = xtime.Duration(time.Second * 10)
	}
	if conf.Timeout <= 0 {
		conf.Timeout = xtime.Duration(time.Millisecond * 250)
	}
	if conf.Subset <= 0 {
		conf.Subset = 50
	}
	if conf.KeepAliveInterval <= 0 {
		conf.KeepAliveInterval = xtime.Duration(time.Second * 60)
	}
	if conf.KeepAliveTimeout <= 0 {
		conf.KeepAliveTimeout = xtime.Duration(time.Second * 20)
	}

	// FIXME(maojian) check Method dial/timeout
	c.mutex.Lock()
	c.conf = conf
	if c.breaker == nil {
		c.breaker = breaker.NewGroup(conf.Breaker)
	} else {
		c.breaker.Reload(conf.Breaker)
	}
	c.mutex.Unlock()
	return nil
}

// Use attachs a global inteceptor to the Client.
// For example, this is the right place for a circuit breaker or error management inteceptor.
func (c *Client) Use(handlers ...grpc.UnaryClientInterceptor) *Client {
	finalSize := len(c.handlers) + len(handlers)
	if finalSize >= int(_abortIndex) {
		panic("warden: client use too many handlers")
	}
	mergedHandlers := make([]grpc.UnaryClientInterceptor, finalSize)
	copy(mergedHandlers, c.handlers)
	copy(mergedHandlers[len(c.handlers):], handlers)
	c.handlers = mergedHandlers
	return c
}

// UseOpt attachs a global grpc DialOption to the Client.
func (c *Client) UseOpt(opts ...grpc.DialOption) *Client {
	c.opts = append(c.opts, opts...)
	return c
}

func (c *Client) cloneOpts() []grpc.DialOption {
	dialOptions := make([]grpc.DialOption, len(c.opts))
	copy(dialOptions, c.opts)
	return dialOptions
}

func (c *Client) dial(ctx context.Context, target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	dialOptions := c.cloneOpts()
	if !c.conf.NonBlock {
		dialOptions = append(dialOptions, grpc.WithBlock())
	}
	dialOptions = append(dialOptions, grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                time.Duration(c.conf.KeepAliveInterval),
		Timeout:             time.Duration(c.conf.KeepAliveTimeout),
		PermitWithoutStream: !c.conf.KeepAliveWithoutStream,
	}))
	dialOptions = append(dialOptions, opts...)

	// init default handler
	var handlers []grpc.UnaryClientInterceptor
	handlers = append(handlers, c.recovery())
	handlers = append(handlers, clientLogging(dialOptions...))
	handlers = append(handlers, c.handlers...)
	// NOTE: c.handle must be a last interceptor.
	handlers = append(handlers, c.handle())

	dialOptions = append(dialOptions, grpc.WithUnaryInterceptor(chainUnaryClient(handlers)))
	c.mutex.RLock()
	conf := c.conf
	c.mutex.RUnlock()
	if conf.Dial > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(conf.Dial))
		defer cancel()
	}
	if u, e := url.Parse(target); e == nil {
		v := u.Query()
		for _, c := range c.conf.Clusters {
			v.Add(naming.MetaCluster, c)
		}
		if c.conf.Zone != "" {
			v.Add(naming.MetaZone, c.conf.Zone)
		}
		if v.Get("subset") == "" && c.conf.Subset > 0 {
			v.Add("subset", strconv.FormatInt(int64(c.conf.Subset), 10))
		}
		u.RawQuery = v.Encode()
		// 比较_grpcTarget中的appid是否等于u.path中的appid，并替换成mock的地址
		for _, t := range _grpcTarget {
			strs := strings.SplitN(t, "=", 2)
			if len(strs) == 2 && ("/"+strs[0]) == u.Path {
				u.Path = "/" + strs[1]
				u.Scheme = "passthrough"
				u.RawQuery = ""
				break
			}
		}
		target = u.String()
	}
	if conn, err = grpc.DialContext(ctx, target, dialOptions...); err != nil {
		fmt.Fprintf(os.Stderr, "warden client: dial %s error %v!", target, err)
	}
	err = errors.WithStack(err)
	return
}

// Dial creates a client connection to the given target.
// Target format is scheme://authority/endpoint?query_arg=value
// example: discovery://default/account.account.service?cluster=shfy01&cluster=shfy02
func (c *Client) Dial(ctx context.Context, target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	opts = append(opts, grpc.WithInsecure())
	return c.dial(ctx, target, opts...)
}

// DialTLS creates a client connection over tls transport to the given target.
func (c *Client) DialTLS(ctx context.Context, target string, file string, name string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	var creds credentials.TransportCredentials
	creds, err = credentials.NewClientTLSFromFile(file, name)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	opts = append(opts, grpc.WithTransportCredentials(creds))
	return c.dial(ctx, target, opts...)
}

// chainUnaryClient creates a single interceptor out of a chain of many interceptors.
//
// Execution is done in left-to-right order, including passing of context.
// For example ChainUnaryClient(one, two, three) will execute one before two before three.
func chainUnaryClient(handlers []grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	n := len(handlers)
	if n == 0 {
		return func(ctx context.Context, method string, req, reply interface{},
			cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
	}

	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var (
			i            int
			chainHandler grpc.UnaryInvoker
		)
		chainHandler = func(ictx context.Context, imethod string, ireq, ireply interface{}, ic *grpc.ClientConn, iopts ...grpc.CallOption) error {
			if i == n-1 {
				return invoker(ictx, imethod, ireq, ireply, ic, iopts...)
			}
			i++
			return handlers[i](ictx, imethod, ireq, ireply, ic, chainHandler, iopts...)
		}

		return handlers[0](ctx, method, req, reply, cc, chainHandler, opts...)
	}
}

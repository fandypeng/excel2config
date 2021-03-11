package auth

import (
	"context"
	"excel2config/internal/dao"
	"excel2config/internal/def"
	"excel2config/internal/server/sessions"
	"fmt"
	"github.com/go-kratos/kratos/pkg/ecode"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/prometheus/common/log"
	"net/http"
	"strings"
)

func NewAuthService(d dao.Dao) *Service {
	return &Service{dao: d}
}

type Service struct {
	dao dao.Dao
}

func (a *Service) NeedLogin() bm.HandlerFunc {
	return func(c *bm.Context) {
		sess := sessions.Default(c)
		var uid, token string
		if uidInterface := sess.Get("uid"); uidInterface != nil {
			uid = uidInterface.(string)
		}
		if tokenInterface := sess.Get("token"); tokenInterface != nil {
			token = tokenInterface.(string)
		}
		//return //TODO: check login status
		if uid == "" || token == "" {
			q := c.Request.URL.Query()
			uid = q.Get("uid")
			token = q.Get("token")
		}
		serverToken, err := a.dao.GetToken(context.TODO(), uid)
		if err != nil {
			c.JSON(nil, ecode.Int(int(def.ErrNeedLogin)))
			c.Abort()
			return
		}
		log.With("token", token).With("serverToken", serverToken).Debug("token compare")
		if token != serverToken {
			c.JSON(nil, ecode.Int(int(def.ErrNeedLogin)))
			c.Abort()
			return
		}
	}
}

func (a *Service) WsNeedLogin(uid, token string) (authed bool) {
	return
}

func (a *Service) NeedAuthorization(uid, gridKey string) (authed bool) {
	return
}

func (a *Service) CORS(allowOriginHosts []string) bm.HandlerFunc {
	return func(c *bm.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			var allow bool
			for _, allowOriginHost := range allowOriginHosts {
				if strings.HasSuffix(strings.ToLower(origin), allowOriginHost) {
					allow = true
					break
				}
			}
			if !allow {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}
		var headerKeys []string
		for k := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerKeys = append(headerKeys, "X-Token")
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("content-type, access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			header := c.Writer.Header()
			header.Set("Access-Control-Allow-Origin", origin)
			header.Set("Access-Control-Allow-Headers", headerStr)
			header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			header.Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			header.Set("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}

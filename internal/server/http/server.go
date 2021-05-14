package http

import (
	pb "excel2config/api"
	"excel2config/internal/model"
	"excel2config/internal/pkg/gzip"
	"excel2config/internal/server/sessions"
	"excel2config/internal/server/sessions/cookie"
	"excel2config/internal/service"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
)

// New new a bm server.
func New(s *service.Service) (engine *bm.Engine, err error) {
	var (
		cfg struct {
			bm.ServerConfig
			CrossDomains []string
		}
		ct paladin.TOML
	)
	if err = paladin.Get("http.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Server").UnmarshalTOML(&cfg); err != nil {
		return
	}
	engine = bm.DefaultServer(&cfg.ServerConfig)
	engine.Use(s.As.CORS(cfg.CrossDomains))
	engine.Use(gzip.Gzip(gzip.DefaultCompression))
	initRouter(engine, s)
	err = engine.Start()
	return
}

func initRouter(engine *bm.Engine, s *service.Service) {
	cookieStore := cookie.NewStore([]byte("secret_cookie"))
	cookieStore.Options(sessions.Options{
		Path:   "/",
		Domain: "http://localhost",
		MaxAge: 86400,
		//Secure:   true,
		//SameSite: http.SameSiteNoneMode,
	})
	engine.Use(sessions.Sessions("e2c_session", cookieStore))
	pb.RegisterSheet(engine, s.As, s)
	pb.RegisterUser(engine, s.As, s)
	pb.RegisterGroup(engine, s.As, s)
}

// example for http request handler.
func howToStart(c *bm.Context) {
	k := &model.Kratos{
		Hello: "Golang 大法好 !!!",
	}
	c.JSON(k, nil)
}

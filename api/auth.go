package api

import bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"

type Auth interface {
	NeedLogin() bm.HandlerFunc
	WsNeedLogin(uid, token string) (authed bool)
	NeedAuthorization(uid, gridKey string) (authed bool)
	CORS(allowOriginHosts []string) bm.HandlerFunc
}

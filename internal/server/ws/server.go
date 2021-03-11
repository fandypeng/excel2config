package ws

import (
	"context"
	"errors"
	"excel2config/internal/dao"
	"excel2config/internal/def"
	"excel2config/internal/model"
	"github.com/go-kratos/kratos/pkg/ecode"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

type Server struct {
	*http.Server
	d dao.Dao
}

var (
	mgr *ClientMgr

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func New(dao dao.Dao) *Server {
	wss := &Server{
		Server: &http.Server{
			Addr:    ":8001",
			Handler: nil,
		},
		d: dao,
	}
	wss.Handler = wss.defaultHandler()
	go func() {
		log.Info("websocket server listen at: 8001")
		if err := wss.Server.ListenAndServe(); err != nil {
			panic("websocket server start failed")
		}
	}()
	mgr = new(ClientMgr)
	mgr.Clients = make(map[string]*Client)
	return wss
}

func (wss *Server) defaultHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		userInfo, gridKey, err := wss.authConn(r)
		if err != nil || userInfo == nil {
			log.Errorw(context.TODO(), "ws conn auth failed: ", err)
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(err.Error()))
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Errorw(context.TODO(), "ws conn upgrade error: ", err)
			return
		}
		client := NewClient(conn, gridKey)
		mgr.AddClient(userInfo.Uid, userInfo.UserName, client)
		svr := newService(client, wss.d)
		go svr.readAndServe()
		go client.waitAndWrite()
	})
	return mux
}

func (wss *Server) authConn(r *http.Request) (userInfo *model.UserInfo, gridKey string, err error) {
	values := r.URL.Query()
	gridKey = values.Get("g")
	if gridKey == "" {
		log.Errorw(context.TODO(), "ws conn without gridKey")
		err = errors.New("gridKey not specified")
		return
	}
	//parse uid and token from url
	path := r.URL.Path
	pathInfo := strings.Split(path, "/")
	if len(pathInfo) < 5 {
		log.Errorw(context.TODO(), "ws conn auth param error, path: ", path)
		return
	}
	uid := pathInfo[2]
	token := pathInfo[4]
	ctx := context.TODO()
	excelInfo, err := wss.d.ExcelInfo(ctx, gridKey)
	if err != nil {
		log.Errorw(context.TODO(), "ws conn gridKey not exist, gridKey: ", gridKey)
		err = errors.New("excel is not exist")
		return
	}
	gid := excelInfo.GroupId
	groupInfo, err := wss.d.GroupInfo(ctx, gid)
	if err != nil {
		log.Errorw(context.TODO(), "ws conn group not exist, groupId: ", gid)
		err = errors.New("group is not exist")
		return
	}
	var havePermissions = false
	for _, m := range groupInfo.Members {
		if m.Uid == uid {
			havePermissions = true
		}
	}
	if !havePermissions {
		log.Errorw(context.TODO(), "ws conn user do not have promission, uid: ", uid, " gid: ", gid)
		err = errors.New("group is not exist")
		return
	}

	serverToken, err := wss.d.GetToken(context.TODO(), uid)
	if err != nil {
		return
	}
	if token != serverToken {
		err = ecode.Int(int(def.ErrNeedLogin))
		return
	}
	userInfo, err = wss.d.GetUserByUid(context.TODO(), uid)
	return
}

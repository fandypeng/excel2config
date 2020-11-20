package ws

import (
	"context"
	pb "excel2config/api"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/gorilla/websocket"
	"net/http"
)

type Server struct {
	*http.Server
	s pb.SheetBMServer
}

var (
	mgr         *ClientMgr
	uidAutoIncr int

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func New(svc pb.SheetBMServer) *Server {
	wss := &Server{
		Server: &http.Server{
			Addr:    ":8001",
			Handler: nil,
		},
	}
	wss.Handler = defaultHandler()
	wss.s = svc
	go func() {
		log.Info("websocket server listen at: 8001")
		if err := wss.Server.ListenAndServe(); err != nil {
			panic("websocket server start failed")
		}
	}()
	mgr = new(ClientMgr)
	mgr.Clients = make(map[int]*Client)
	return wss
}

func defaultHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Errorw(context.TODO(), "ws conn upgrade error: ", err)
		}
		uidAutoIncr += 1
		client := NewClient(conn)
		mgr.AddClient(uidAutoIncr, client)
		go client.readAndServe()
		go client.waitAndWrite()
	})
	return mux
}

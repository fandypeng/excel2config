package ws

import (
	"context"
	pb "excel2config/api"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
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
	wss.Handler = defaultHandler(wss)
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

func defaultHandler(wss *Server) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/excel", http.HandlerFunc(wss.LoadExcel))
	mux.Handle("/excel/sheet", http.HandlerFunc(wss.LoadExcelSheet))
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

func (wss *Server) LoadExcel(writer http.ResponseWriter, request *http.Request) {
	setupHeader(writer)
	resp, err := wss.s.LoadExcel(context.TODO(), &pb.LoadExcelReq{
		GridKey: "comment",
	})
	if err != nil {
		writer.Write([]byte("Internal server error, " + err.Error()))
		return
	}
	writer.Write([]byte(resp.Jsonstr))
}

func (wss *Server) LoadExcelSheet(writer http.ResponseWriter, request *http.Request) {
	setupHeader(writer)
	request.ParseForm()
	index := request.FormValue("index")
	gridKey := request.FormValue("gridKey")
	resp, err := wss.s.LoadExcelSheet(context.TODO(), &pb.LoadExcelSheetReq{
		GridKey: gridKey,
		Indexs:  strings.Split(index, ","),
	})
	if err != nil {
		writer.Write([]byte("Internal server error, " + err.Error()))
		return
	}
	writer.Write([]byte(resp.Jsonstr))
}

func setupHeader(writer http.ResponseWriter) {
	writer.Header().Add("Access-Control-Allow-Origin", "*")
	writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Add("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	writer.Header().Add("Access-Control-Allow-Credentials", "true")
}

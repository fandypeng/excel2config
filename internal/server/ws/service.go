package ws

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"excel2config/internal/dao"
	"excel2config/internal/model"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/gorilla/websocket"
	"golang.org/x/text/encoding/charmap"
	"io/ioutil"
	"net/url"
	"time"
)

type handler func(ctx context.Context, reqmsg []byte)

type service struct {
	*Client
	d        dao.Dao
	handlers map[string]handler
}

func newService(c *Client, d dao.Dao) *service {
	s := &service{
		Client: c,
		d:      d,
	}
	s.handlers = initHandlers(s)
	return s
}

func initHandlers(s *service) map[string]handler {
	return map[string]handler{
		"v":    s.updateGrid,
		"rv":   s.updateGridMulti,
		"cg":   s.updateGridConfig,
		"all":  s.updateGridCommon,
		"fc":   s.updateCalcChain,
		"drc":  s.updateRowColumn,
		"arc":  s.updateRowColumn,
		"fsc":  s.updateFilter,
		"fsr":  s.updateFilter,
		"sha":  s.addSheet,
		"shc":  s.copySheet,
		"shd":  s.deleteSheet,
		"shre": s.recoverSheet,
		"shr":  s.updateSheetOrder,
		"shs":  s.toggleSheet,
		"sh":   s.hideOrShowSheet,
	}
}

func (s *service) readAndServe() {
	defer func() {
		s.Close()
		s.remDeletedSheetWhenLeave()
		log.Infow(context.TODO(), "read groutine exited, uid", s.uid)
	}()
	s.setReadOpts()
	for {
		if s.isClosed {
			break
		}
		messageType, message, err := s.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Errorw(context.TODO(), "err", err, "msg", "ws error")
			}
			log.Errorw(context.TODO(), "ws recv msg err", err)
			s.remDeletedSheetWhenLeave()
			break
		}
		reqmsg, err := s.ungzip(message)
		if err != nil {
			log.Errorw(context.TODO(), "ungzip error, err: ", err)
			continue
		}
		log.Infow(context.TODO(), "uid: ", s.uid, "message_type", messageType, "msg", string(reqmsg))
		reqmsg = bytes.TrimSpace(bytes.Replace(reqmsg, newline, space, -1))
		s.handleRequest(context.Background(), reqmsg)
	}
}

func (s *service) handleRequest(ctx context.Context, reqmsg []byte) {
	if s.isClosed {
		return
	}
	var msg struct {
		T string `json:"t"`
	}
	var rsp struct {
		Type     int    `json:"type"`
		Id       string `json:"id,omitempty"`
		UserName string `json:"username,omitempty"`
		Data     string `json:"data"`
	}
	json.Unmarshal(reqmsg, &msg)
	uid := s.GetUid()
	switch msg.T {
	case "v", "rv", "rv_end", "cg", "all", "fc", "drc", "arc", "f", "fsc", "fsr", "sha", "shc", "shd", "shr", "shre", "sh", "c", "na":
		rsp.Type = 2
	case "mv":
		rsp.Type = 3
		rsp.Id = uid
		rsp.UserName = s.GetName()
	case "": //离线情况下把更新指令打包批量下发给客户端
		rsp.Type = 4
	default:
		rsp.Type = 1
	}

	handler, ok := s.handlers[msg.T]
	if ok {
		handler(ctx, reqmsg)
	}
	rsp.Data = string(reqmsg)
	jsonstr, _ := json.Marshal(rsp)
	mgr.Send2AllClients(s.Client, jsonstr)
}

// step1: encode bytes to IOS-8859-1, from https://stackoverflow.com/questions/40980402/compressed-with-pakozlib-in-javascript-decompressing-with-zlibpython-not-wo
// step2: gzip decode bytes
// step3: url decode
func (s *service) ungzip(gzipmsg []byte) (reqmsg []byte, err error) {
	if len(gzipmsg) == 0 {
		return
	}
	if string(gzipmsg) == "rub" {
		reqmsg = gzipmsg
		return
	}
	e := charmap.ISO8859_1.NewEncoder()
	encodeMsg, err := e.Bytes(gzipmsg)
	if err != nil {
		return
	}
	b := bytes.NewReader(encodeMsg)
	r, err := gzip.NewReader(b)
	if err != nil {
		return
	}
	defer r.Close()
	reqmsg, err = ioutil.ReadAll(r)
	if err != nil {
		return
	}
	reqstr, err := url.QueryUnescape(string(reqmsg))
	if err != nil {
		return
	}
	reqmsg = []byte(reqstr)
	return
}

func (s *service) updateGrid(ctx context.Context, reqmsg []byte) {
	req := new(model.UpdateV)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "jsonstr", string(reqmsg), "msg", "json unmarshal error")
		return
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	err = s.d.UpdateGridValue(ctx, s.gridKey, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", s.gridKey, "req", req, "msg", "update grid failed")
	}
}

func (s *service) updateGridMulti(ctx context.Context, reqmsg []byte) {
	req := new(model.UpdateRV)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "jsonstr", string(reqmsg), "msg", "json unmarshal error")
		return
	}
	if len(req.Range.Column) < 2 || len(req.Range.Row) < 2 {
		log.Errorw(ctx, "req", req, "msg", "invalid params")
		return
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	err = s.d.UpdateGridMulti(ctx, s.gridKey, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", s.gridKey, "req", req, "msg", "update grid failed")
	}
}

func (s *service) updateGridConfig(ctx context.Context, reqmsg []byte) {
	req := new(model.UpdateCG)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "jsonstr", string(reqmsg), "msg", "json unmarshal error")
		return
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	err = s.d.UpdateGridConfig(ctx, s.gridKey, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", s.gridKey, "req", req, "msg", "update grid failed")
	}
}

func (s *service) updateGridCommon(ctx context.Context, reqmsg []byte) {
	req := new(model.UpdateCommon)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "jsonstr", string(reqmsg), "msg", "json unmarshal error")
		return
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	err = s.d.UpdateGridCommon(ctx, s.gridKey, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", s.gridKey, "req", req, "msg", "update grid failed")
	}
}

func (s *service) updateCalcChain(ctx context.Context, reqmsg []byte) {
	req := new(model.UpdateCalcChain)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "jsonstr", string(reqmsg), "msg", "json unmarshal error")
		return
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	err = s.d.UpdateCalcChain(ctx, s.gridKey, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", s.gridKey, "req", req, "msg", "update calc chain failed")
	}
}

func (s *service) updateRowColumn(ctx context.Context, reqmsg []byte) {
	req := new(model.UpdateRowColumn)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "jsonstr", string(reqmsg), "msg", "json unmarshal error")
		return
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	err = s.d.UpdateRowColumn(ctx, s.gridKey, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", s.gridKey, "req", req, "msg", "update calc chain failed")
	}
}

func (s *service) updateFilter(ctx context.Context, reqmsg []byte) {
	req := new(model.UpdateFilter)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "jsonstr", string(reqmsg), "msg", "json unmarshal error")
		return
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	err = s.d.UpdateFilter(ctx, s.gridKey, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", s.gridKey, "req", req, "msg", "update calc chain failed")
	}
}

func (s *service) addSheet(ctx context.Context, reqmsg []byte) {
	req := new(model.AddSheet)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "jsonstr", string(reqmsg), "msg", "json unmarshal error")
		return
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	err = s.d.AddSheet(ctx, s.gridKey, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", s.gridKey, "req", req, "msg", "add sheet failed")
	}
}

func (s *service) copySheet(ctx context.Context, reqmsg []byte) {
	req := new(model.CopySheet)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "jsonstr", string(reqmsg), "msg", "json unmarshal error")
		return
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	err = s.d.CopySheet(ctx, s.gridKey, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", s.gridKey, "req", req, "msg", "add sheet failed")
	}
}

func (s *service) deleteSheet(ctx context.Context, reqmsg []byte) {
	req := new(model.DeleteSheet)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "jsonstr", string(reqmsg), "msg", "json unmarshal error")
		return
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	err = s.d.DeleteSheet(ctx, s.gridKey, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", s.gridKey, "req", req, "msg", "add sheet failed")
	}
}

func (s *service) recoverSheet(ctx context.Context, reqmsg []byte) {
	req := new(model.RecoverSheet)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "jsonstr", string(reqmsg), "msg", "json unmarshal error")
		return
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	err = s.d.RecoverSheet(ctx, s.gridKey, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", s.gridKey, "req", req, "msg", "add sheet failed")
	}
}

func (s *service) updateSheetOrder(ctx context.Context, reqmsg []byte) {
	req := new(model.UpdateSheetOrder)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "jsonstr", string(reqmsg), "msg", "json unmarshal error")
		return
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	err = s.d.UpdateSheetOrder(ctx, s.gridKey, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", s.gridKey, "req", req, "msg", "add sheet failed")
	}
}

func (s *service) toggleSheet(ctx context.Context, reqmsg []byte) {
	req := new(model.ToggleSheet)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "jsonstr", string(reqmsg), "msg", "json unmarshal error")
		return
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	err = s.d.ToggleSheet(ctx, s.gridKey, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", s.gridKey, "req", req, "msg", "add sheet failed")
	}
}

func (s *service) hideOrShowSheet(ctx context.Context, reqmsg []byte) {
	req := new(model.HideOrShowSheet)
	err := json.Unmarshal(reqmsg, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "jsonstr", string(reqmsg), "msg", "json unmarshal error")
		return
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	err = s.d.HideOrShowSheet(ctx, s.gridKey, req)
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", s.gridKey, "req", req, "msg", "add sheet failed")
	}
}

func (s *service) remDeletedSheetWhenLeave() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := s.d.RemDeletedSheet(ctx, s.gridKey)
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", s.gridKey, "msg", "rem deleted sheet failed")
	}
}

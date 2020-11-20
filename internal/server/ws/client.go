package ws

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/prometheus/common/log"
	"strconv"
	"sync"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 60 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 20 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 10240
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	*websocket.Conn
	msgs     chan []byte
	isClosed bool
	uid      int
	hander   func(msg string)

	sync.RWMutex
}

func NewClient(c *websocket.Conn) *Client {
	return &Client{
		Conn:     c,
		msgs:     make(chan []byte),
		isClosed: false,
	}
}

func (c *Client) readAndServe() {
	log.Debugln("new ws connect")
	defer func() {
		c.Close()
	}()
	c.setReadOpts()
	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Errorln("ws socket error: %v", err)
			}
			break
		}
		log.Infoln("uid: ", c.uid, " recv message_type: ", messageType, ", msg: ", string(message))
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.handleRequest(message)
	}
}

func (c *Client) setReadOpts() {
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(appData string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
}

func (c *Client) waitAndWrite() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()
	for {
		select {
		case message, ok := <-c.msgs:
			if !ok {
				log.Warnln("ws conn msg channel closed")
			}
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.With("message", string(message)).Errorf("ws conn write msg error: %v", err)
				return
			}
		case <-ticker.C:
			log.With("uid", c.uid).Debugln("ticker triggered")
			if c.isClosed {
				break
			}
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.Conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Errorf("ws conn write msg error: %v", err)
				return
			}
		}
	}
}

func (c *Client) Close() {
	c.Lock()
	defer c.Unlock()
	c.Conn.Close()
	c.isClosed = true
	c.msgs = nil
	mgr.DelClient(c.uid)
}

func (c *Client) GetUid() int {
	return c.uid
}

func (c *Client) handleRequest(reqmsg []byte) {
	if c.isClosed {
		return
	}
	var msg struct {
		T string `json:"t"`
	}
	var rsp struct {
		Type     int    `json:"type"`
		Id       int    `json:"id,omitempty"`
		UserName string `json:"username,omitempty"`
		Data     string `json:"data"`
	}
	json.Unmarshal(reqmsg, &msg)
	uid := c.GetUid()
	switch msg.T {
	case "v", "rv", "cg", "all", "fc", "drc", "arc", "f", "fsc", "fsr", "sha", "shc", "shd", "shr", "shre", "sh", "c", "na":
		rsp.Type = 2
	case "mv":
		rsp.Type = 3
		rsp.Id = uid
		rsp.UserName = "Guest" + strconv.Itoa(uid)
	case "rv_end": //离线情况下把更新指令打包批量下发给客户端
		rsp.Type = 4
	default:
		rsp.Type = 1
	}
	rsp.Data = string(reqmsg)
	//fmt.Println(fmt.Sprintf("uid: %d, receive msg: %v", uid, string(reqmsg)))
	jsonstr, _ := json.Marshal(rsp)
	mgr.Send2AllClients(c, jsonstr)
}

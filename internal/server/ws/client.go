package ws

import (
	"context"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/gorilla/websocket"
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
	maxMessageSize = 102400
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	*websocket.Conn
	msgs     chan []byte
	isClosed bool
	uid      string
	name     string
	hander   func(msg string)
	gridKey  string

	sync.RWMutex
}

func NewClient(c *websocket.Conn, gridKey string) *Client {
	return &Client{
		Conn:     c,
		msgs:     make(chan []byte),
		isClosed: false,
		gridKey:  gridKey,
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
				log.Warn("ws conn msg channel closed")
			}
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Errorw(context.TODO(), "message", string(message), "err", err, "msg", "ws conn write msg error")
				return
			}
		case <-ticker.C:
			log.Infow(context.TODO(), "uid", c.uid, "msg", "ticker triggered")
			if c.isClosed {
				break
			}
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.Conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Errorw(context.TODO(), "err", err, "msg", "ws conn write msg error")
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
}

func (c *Client) GetUid() string {
	return c.uid
}

func (c *Client) GetName() string {
	return c.name
}

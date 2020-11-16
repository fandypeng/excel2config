package ws

import (
	"context"
	"github.com/go-kratos/kratos/pkg/log"
	"sync"
)

type ClientMgr struct {
	Clients map[int]*Client
	sync.RWMutex
}

func (m *ClientMgr) AddClient(uid int, client *Client) {
	m.Lock()
	defer m.Unlock()
	client.uid = uid
	m.Clients[uid] = client
	log.Infoc(context.TODO(),"add client, uid: %d, len: %d", uid, len(m.Clients))
}

func (m *ClientMgr) DelClient(uid int) {
	m.Lock()
	defer m.Unlock()
	delete(m.Clients, uid)
}

func (m *ClientMgr) Send2AllClients(conn *Client, msg []byte) {
	m.RLock()
	defer m.RUnlock()
	for _, c := range m.Clients {
		if c == conn {
			continue
		}
		c.msgs <- msg
	}
}
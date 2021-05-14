package ws

import (
	"context"
	"github.com/go-kratos/kratos/pkg/log"
	"sync"
)

type ClientMgr struct {
	Clients map[string]*Client
	sync.RWMutex
}

func (m *ClientMgr) AddClient(uid, name string, client *Client) {
	m.Lock()
	defer m.Unlock()
	oldClient, ok := m.Clients[uid]
	if ok {
		oldClient.Close()
	}
	client.uid = uid
	client.name = name
	m.Clients[uid] = client
	log.Infoc(context.TODO(), "add client, uid: %v, len: %d", uid, len(m.Clients))
}

func (m *ClientMgr) DelClient(uid string) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.Clients[uid]; ok {
		delete(m.Clients, uid)
	}
	log.Info("del client: %s, client num: %d", uid, len(m.Clients))
}

func (m *ClientMgr) Send2AllClients(conn *Client, msg []byte) {
	m.RLock()
	defer m.RUnlock()
	for _, c := range m.Clients {
		if c.uid == conn.uid || c.gridKey != conn.gridKey {
			continue
		}
		c.msgs <- msg
	}
}

func (m *ClientMgr) checkConnAlive() {

}

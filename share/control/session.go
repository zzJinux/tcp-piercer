package control

import (
	"net"
	"sync"

	"github.com/zzJinux/tcp-piercer/share/message"
)

type Session struct {
	Conn            net.Conn
	MsgChan         *message.MessageChan
	PrivateEndpoint *net.TCPAddr
	PublicEndpoint  *net.TCPAddr
}

type Sessions struct {
	sync.RWMutex
	m map[string]*Session
}

func (u *Sessions) Get(key string) (*Session, bool) {
	u.RLock()
	v, found := u.m[key]
	u.RUnlock()
	return v, found
}

func (u *Sessions) Set(key string, value *Session) {
	u.Lock()
	u.m[key] = value
	u.Unlock()
}

func (u *Sessions) Del(key string) {
	u.Lock()
	delete(u.m, key)
	u.Unlock()
}

package session

import (
	"net"
	"sync"
)

// Session结构体
type Session struct {
	sync.Mutex
	ID   string
	Conn net.Conn
}

// Session管理器
type SessionManager struct {
	sessions map[string]*Session
}

// 获取Session
func (sm *SessionManager) GetSession(id string) *Session {
	return sm.sessions[id]
}

// 发送Session命令
func (sm *SessionManager) SessionCommand(id, cmd string) {
	session := sm.GetSession(id)
	session.Lock()
	defer session.Unlock()

	//TODO
}

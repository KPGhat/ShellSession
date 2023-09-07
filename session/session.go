package session

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

// Session结构体
type Session struct {
	Conn      net.Conn
	isAlive   bool
	readLock  *sync.Mutex
	writeLock *sync.Mutex
	//Buffer []byte
}

// Send data to session
func (session *Session) Send(data []byte) {
	session.writeLock.Lock()
	defer session.writeLock.Unlock()

	_, err := session.Conn.Write(data)
	if err != nil {
		log.Printf("[-]Send data to sessioin error: %v\n", err)
		session.isAlive = false
		return
	}
}

// Read data from session
func (session *Session) Read() []byte {
	session.readLock.Lock()
	defer session.readLock.Unlock()

	// fixme maybe have len bug
	data := make([]byte, 0x100)
	readLen, err := session.Conn.Read(data)
	if err != nil {
		log.Printf("[-]Read data to sessioin error: %v\n", err)
		session.isAlive = false
		return nil
	}

	if readLen == 0 {
		return nil
	}
	return data
}

// Session Manager
type SessionManager struct {
	sessions []*Session
}

var globalSessionManager SessionManager

func GetSessionManager() *SessionManager {
	return &globalSessionManager
}

// GET a Session
func (sm *SessionManager) GetSession(id int) *Session {
	if id < len(sm.sessions) {
		return sm.sessions[id]
	}

	return nil
}

// ADD a Session
func (sm *SessionManager) AddSession(conn net.Conn) {
	sm.sessions = append(sm.sessions, &Session{
		Conn:      conn,
		isAlive:   true,
		readLock:  &sync.Mutex{},
		writeLock: &sync.Mutex{},
	})
}

func (sm *SessionManager) ListAllSession(output io.Writer) {
	if len(sm.sessions) == 0 {
		output.Write([]byte("[-]No session established\n"))
		return
	}
	for i, session := range sm.sessions {
		remoteAddr := session.Conn.RemoteAddr().String()
		isAlive := "true"
		if !session.isAlive {
			isAlive = "false"
		}
		sessionInfo := fmt.Sprintf("id:%d\thost:%s\talive:%s\n", i, remoteAddr, isAlive)
		_, err := output.Write([]byte(sessionInfo))
		if err != nil {
			log.Printf("[-]Session list: %v\n", err)
			return
		}
	}
}

func (sm *SessionManager) DestroySessionManager() {
	for _, session := range sm.sessions {
		err := session.Conn.Close()
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
}

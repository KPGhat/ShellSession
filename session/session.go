package session

import (
	"fmt"
	"github.com/KPGhat/ShellSession/cmd"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// Session结构体
type Session struct {
	Conn      net.Conn
	isAlive   bool
	Buffer    []byte
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
func (session *Session) Read(data []byte) int {
	session.readLock.Lock()
	defer session.readLock.Unlock()

	readLen, err := session.Conn.Read(data)
	if err != nil {
		log.Printf("[-]Read data to sessioin error: %v\n", err)
		session.isAlive = false
		return 0
	}

	return readLen
}

func (session *Session) ReadListener(running *bool, callback func([]byte)) {
	for *running {
		session.Conn.SetReadDeadline(time.Time{})
		// fixme maybe have len bug
		data := make([]byte, 0x100)
		n := session.Read(data)

		if n > 0 {
			callback(data[:n])
		}
	}
}

func (session *Session) SessionInfo() string {
	remoteAddr := session.Conn.RemoteAddr().String()
	isAlive := "true"
	if !session.isAlive {
		isAlive = "false"
	}
	return fmt.Sprintf("host: %s\talive: %s", remoteAddr, isAlive)
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
	newSession := &Session{
		Conn:      conn,
		isAlive:   true,
		readLock:  &sync.Mutex{},
		writeLock: &sync.Mutex{},
	}
	sessLen := len(sm.sessions)
	sm.sessions = append(sm.sessions, newSession)

	log.Println(fmt.Sprintf("[+]Add Session %d:\t", sessLen) + newSession.SessionInfo())
}

func (sm *SessionManager) ListAllSession(output io.Writer) {
	if len(sm.sessions) == 0 {
		output.Write([]byte("[-]No session established\n"))
		return
	}
	for i, session := range sm.sessions {
		sessionInfo := fmt.Sprintf("id: %d\t", i) + session.SessionInfo()
		_, err := output.Write([]byte(sessionInfo))
		if err != nil {
			log.Printf("[-]Session list: %v\n", err)
			return
		}
	}
}

func (sm *SessionManager) ExecCmdForAll(command string, output io.Writer) {
	limiter := make(chan struct{}, cmd.Config.MaxConn/2+1)
	wg := sync.WaitGroup{}
	for _, session := range sm.sessions {
		limiter <- struct{}{}
		wg.Add(1)

		// TODO add get result and store the result
		go func(sess *Session) {
			sess.Send([]byte(command + "\n"))
			<-limiter
			wg.Done()
		}(session)
	}

	wg.Wait()
}

func (sm *SessionManager) DestroySessionManager() {
	for _, session := range sm.sessions {
		err := session.Conn.Close()
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
}

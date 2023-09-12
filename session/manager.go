package session

import (
	"errors"
	"fmt"
	"github.com/KPGhat/ShellSession/cmd"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

// Session Manager
type SessionManager struct {
	sessions []*Session
	context  map[int]struct{}
	mu       sync.Mutex
}

var globalSessionManager SessionManager

func init() {
	globalSessionManager.context = make(map[int]struct{})
}

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
		_, err := output.Write([]byte(sessionInfo + "\n"))
		if err != nil {
			log.Printf("[-]Session list: %v\n", err)
			return
		}
	}
}

func (sm *SessionManager) AddContext(id int) error {
	if _, ok := sm.context[id]; ok || id >= len(sm.sessions) {
		return errors.New(fmt.Sprintf("[-]Session Manage Context <%d> has already added or not exist\n", id))
	}
	sm.mu.Lock()
	sm.context[id] = struct{}{}
	sm.mu.Unlock()
	return nil
}

func (sm *SessionManager) KeepAliveConn() {
	backup := make([]*Session, len(sm.sessions))
	copy(backup, sm.sessions)
	var aliveHost []*Session
	for _, session := range backup {
		if session.isAlive {
			aliveHost = append(aliveHost, session)
		}
	}

	if len(backup) == len(aliveHost) {
		return
	}
	sm.sessions = make([]*Session, len(aliveHost))
	copy(sm.sessions, aliveHost)
}

func (sm *SessionManager) AddAllContext() {
	id := 0
	for GetSessionManager().AddContext(id) == nil {
		id++
	}
}

func (sm *SessionManager) DelContext(id int) error {
	if _, ok := sm.context[id]; !ok {
		return errors.New(fmt.Sprintf("[-]Session Manage Context <%d> not exist\n", id))
	}
	sm.mu.Lock()
	delete(sm.context, id)
	sm.mu.Unlock()
	return nil
}

func (sm *SessionManager) DelAllContext() {
	for id, _ := range sm.context {
		sm.mu.Lock()
		delete(sm.context, id)
		sm.mu.Unlock()
	}
}

func (sm *SessionManager) GetAllContext() string {
	var result []string
	for id, _ := range sm.context {
		result = append(result, strconv.Itoa(id))
	}
	return strings.Join(result, ",")
}

func (sm *SessionManager) HandleAllContext(callback func(session *Session)) {
	for id, _ := range sm.context {
		callback(sm.GetSession(id))
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
			result := sess.ExecCmd([]byte(command))
			output.Write(result)
			<-limiter
			defer wg.Done()
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

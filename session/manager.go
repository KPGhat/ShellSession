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
	sessions map[int]*Session
	context  map[int]struct{}
	lastID   int
	mu       sync.Mutex
}

var globalSessionManager SessionManager

func init() {
	globalSessionManager.sessions = make(map[int]*Session)
	globalSessionManager.context = make(map[int]struct{})
	globalSessionManager.lastID = 0
}

func GetSessionManager() *SessionManager {
	return &globalSessionManager
}

// GET a Session
func (sm *SessionManager) GetSession(id int) *Session {
	if session, err := sm.sessions[id]; !err {
		return session
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
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.sessions[sm.lastID] = newSession
	log.Println(fmt.Sprintf("[+]Add Session %d:\t", sm.lastID) + newSession.SessionInfo())
	sm.lastID++
}

func (sm *SessionManager) ListAllSession(output io.Writer, onlyAlive bool) {
	if len(sm.sessions) == 0 {
		output.Write([]byte("[-]No session established\n"))
		return
	}
	for i, session := range sm.sessions {
		if onlyAlive && !session.isAlive {
			continue
		}

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

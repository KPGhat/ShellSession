package session

import (
	"fmt"
	"github.com/KPGhat/ShellSession/utils"
	"github.com/fatih/color"
	"io"
	"log"
	"net"
	"sync"
)

type Manager struct {
	sessionManager map[int]*Session
	contextManager map[int]*Context
	lastSessionID  int
	lastContextID  int
	mu             sync.Mutex
}

var globalManager Manager

func init() {
	globalManager.sessionManager = make(map[int]*Session)
	globalManager.contextManager = make(map[int]*Context)
	globalManager.lastSessionID = -1
	globalManager.lastContextID = -1
}

func GetManager() *Manager {
	return &globalManager
}

// GET a Session
func (manager *Manager) GetSession(id int) *Session {
	if session, ok := manager.sessionManager[id]; ok {
		return session
	}

	return nil
}

// ADD a Session
func (manager *Manager) AddSession(conn net.Conn) {
	newSession := &Session{
		Conn:      conn,
		isAlive:   true,
		readLock:  &sync.Mutex{},
		writeLock: &sync.Mutex{},
	}
	manager.mu.Lock()
	defer manager.mu.Unlock()
	manager.lastSessionID++
	manager.sessionManager[manager.lastSessionID] = newSession
	utils.Congrats(fmt.Sprintf("Add Session %d:\t", manager.lastSessionID) + newSession.SessionInfo())
}

func (manager *Manager) ListAllSession(output io.Writer, onlyAlive bool) {
	if len(manager.sessionManager) == 0 {
		red := color.New(color.FgHiRed).SprintfFunc()
		output.Write([]byte(red("[-]No session established\n")))
		return
	}
	for i, session := range manager.sessionManager {
		if onlyAlive && !session.isAlive {
			continue
		}

		sessionInfo := fmt.Sprintf("id: %d\t", i) + session.SessionInfo()
		_, err := output.Write([]byte(sessionInfo + "\n"))
		if err != nil {
			utils.Warning(fmt.Sprintf("Session list: %v", err))
			return
		}
	}
}

func (manager *Manager) ExecCmdForAll(command string, output io.Writer) {
	limiter := make(chan struct{}, 100)
	wg := sync.WaitGroup{}
	for _, session := range manager.sessionManager {
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

func (manager *Manager) DestroySessionManager() {
	for _, session := range manager.sessionManager {
		err := session.Conn.Close()
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
}

func (manager *Manager) CreateContext() int {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	manager.lastContextID++
	manager.contextManager[manager.lastContextID] = InitContext()

	return manager.lastContextID
}

func (manager *Manager) ListAllContext(output io.Writer) {
	if len(manager.sessionManager) == 0 {
		red := color.New(color.FgHiRed).SprintfFunc()
		output.Write([]byte(red("[-]No context created\n")))
		return
	}
	for i, context := range manager.contextManager {
		contextInfo := fmt.Sprintf("id: %d\t", i) + context.ContextInfo()
		_, err := output.Write([]byte(contextInfo + "\n"))
		if err != nil {
			utils.Warning(fmt.Sprintf("Context list: %v", err))
			return
		}
	}
}

func (manager *Manager) GetContext(id int) *Context {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	if context, ok := manager.contextManager[id]; ok {
		return context
	}

	return nil
}

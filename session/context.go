package session

import (
	"errors"
	"fmt"
	"github.com/KPGhat/ShellSession/utils"
	"io"
	"maps"
	"strconv"
	"strings"
	"sync"
)

type Context struct {
	context map[int]struct{}
	ipMap   map[string]struct{}
	mu      sync.Mutex
}

func InitContext() *Context {
	var context Context
	context.context = make(map[int]struct{})
	context.ipMap = make(map[string]struct{})
	return &context
}

func (context *Context) Add(id int) error {
	if _, ok := context.context[id]; ok || id >= len(GetManager().sessionManager) {
		return errors.New(fmt.Sprintf("Session <%d> has already added or not exist\n", id))
	}
	context.mu.Lock()
	ip := strings.Split(GetManager().GetSession(id).Conn.RemoteAddr().String(), ":")[0]
	context.ipMap[ip] = struct{}{}
	context.context[id] = struct{}{}
	context.mu.Unlock()
	return nil
}

func (context *Context) AddAll() {
	for id, session := range GetManager().sessionManager {
		if session.IsAlive {
			context.Add(id)
		}
	}
}

func (context *Context) AddAllIP() {
	var ip string
	for id, session := range GetManager().sessionManager {
		ip = strings.Split(session.Conn.RemoteAddr().String(), ":")[0]
		if _, ok := context.ipMap[ip]; !ok && session.IsAlive {
			context.Add(id)
		}
	}
}

func (context *Context) Delete(id int) error {
	context.mu.Lock()
	defer context.mu.Unlock()
	if _, ok := context.context[id]; !ok {
		return errors.New(fmt.Sprintf("Session Manage Context <%d> not exist\n", id))
	}

	ip := strings.Split(GetManager().GetSession(id).Conn.RemoteAddr().String(), ":")[0]
	delete(context.ipMap, ip)
	delete(context.context, id)
	return nil
}

func (context *Context) DeleteAll() {
	var sessionToDel []int
	for id, _ := range context.context {
		sessionToDel = append(sessionToDel, id)
	}
	for _, id := range sessionToDel {
		context.Delete(id)
	}
}

func (context *Context) ContextInfo() string {
	var result []string
	for id, _ := range context.context {
		result = append(result, strconv.Itoa(id))
	}
	return strings.Join(result, ",")
}

func (context *Context) Size() string {
	return fmt.Sprintf("%d", len(context.context))
}

func (context *Context) List(output io.Writer) {
	if len(context.context) == 0 {
		utils.Error("[-]No context created")
		return
	}
	for id, _ := range context.context {
		sessionInfo := fmt.Sprintf("id: %d\t", id) + GetManager().GetSession(id).SessionInfo()
		_, err := output.Write([]byte(sessionInfo + "\n"))
		if err != nil {
			utils.Error(fmt.Sprintf("Context list: %v", err))
			return
		}
	}
}

func (context *Context) HandleAllContext(callback func(session *Session)) {
	copyContext := maps.Clone(context.context)
	for id, _ := range copyContext {
		session := GetManager().GetSession(id)
		if session != nil {
			callback(session)
		}
	}
}

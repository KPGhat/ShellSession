package session

import (
	"errors"
	"fmt"
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
			context.ipMap[ip] = struct{}{}
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

	delete(context.context, id)
	return nil
}

func (context *Context) DeleteAll() {
	for id, _ := range context.context {
		context.Delete(id)
	}
}

func (context *Context) GetAllContext() string {
	var result []string
	for id, _ := range context.context {
		result = append(result, strconv.Itoa(id))
	}
	return strings.Join(result, ",")
}

func (context *Context) ContextInfo() string {
	return fmt.Sprintf("managing session< %s >", context.GetAllContext())
}

func (context *Context) HandleAllContext(callback func(session *Session)) {
	for id, _ := range context.context {
		session := GetManager().GetSession(id)
		if session != nil {
			callback(session)
		}
	}
}

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
	mu      sync.Mutex
}

func InitContext() *Context {
	var context Context
	context.context = make(map[int]struct{})
	return &context
}

func (context *Context) AddContext(id int) error {
	if _, ok := context.context[id]; ok || id >= len(GetManager().sessionManager) {
		return errors.New(fmt.Sprintf("Session <%d> has already added or not exist\n", id))
	}
	context.mu.Lock()
	context.context[id] = struct{}{}
	context.mu.Unlock()
	return nil
}

func (context *Context) AddAllContext() {
	for id, session := range GetManager().sessionManager {
		if session.isAlive {
			context.AddContext(id)
		}
	}
}

func (context *Context) DelContext(id int) error {
	if _, ok := context.context[id]; !ok {
		return errors.New(fmt.Sprintf("Session Manage Context <%d> not exist\n", id))
	}
	context.mu.Lock()
	delete(context.context, id)
	context.mu.Unlock()
	return nil
}

func (context *Context) DelAllContext() {
	for id, _ := range context.context {
		context.mu.Lock()
		delete(context.context, id)
		context.mu.Unlock()
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

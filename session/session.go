package session

import (
	"bytes"
	"fmt"
	"github.com/KPGhat/ShellSession/utils"
	"net"
	"strings"
	"sync"
	"time"
)

// Session结构体
type Session struct {
	Conn      net.Conn
	IsAlive   bool
	Id        int
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
		utils.Warning(fmt.Sprintf("Send data to sessioin error: %v", err))
		session.IsAlive = false
		return
	}
	session.IsAlive = true
}

// Read data from session
func (session *Session) Read(data []byte) (int, error) {
	session.readLock.Lock()
	defer session.readLock.Unlock()

	readLen, err := session.Conn.Read(data)
	return readLen, err
}

func (session *Session) ReadUntil(suffix []byte) ([]byte, bool) {
	buffer := make([]byte, 1)
	var isTimeout bool
	var data bytes.Buffer

	for {
		session.Conn.SetReadDeadline(time.Now().Add(time.Second))
		n, err := session.Read(buffer)

		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				utils.Warning(fmt.Sprintf("Read data timeout: %v", err))
				isTimeout = true
			} else {
				session.IsAlive = false
				isTimeout = false
			}
			break
		}
		data.Write(buffer[:n])

		if bytes.HasSuffix(data.Bytes(), suffix) {
			break
		}
	}
	session.Conn.SetReadDeadline(time.Time{})
	return data.Bytes(), isTimeout
}

func (session *Session) ReadListener(running *bool, callback func([]byte)) {
	for *running {
		session.Conn.SetReadDeadline(time.Time{})
		// fixme maybe have len bug
		data := make([]byte, 1024)
		n, err := session.Read(data)
		if err != nil {
			utils.Warning(fmt.Sprintf("Read data to session error: %v", err))
			*running = false
			session.IsAlive = false
		}

		if n > 0 {
			callback(data[:n])
		}
	}
}

func (session *Session) ExecCmd(command []byte) []byte {
	prefix := utils.RandString(8)
	suffix := utils.RandString(8)
	newCommand := " echo " + prefix + " && " + string(command) + "; echo " + suffix + "\n"
	session.Send([]byte(newCommand))

	var execResult []byte

	for execResult == nil || strings.EqualFold(" && "+string(command)+"; echo ", string(execResult)) {
		_, isTimeout := session.ReadUntil([]byte(prefix))
		if isTimeout {
			return []byte{}
		}

		result, _ := session.ReadUntil([]byte(suffix))
		var found bool
		execResult, found = bytes.CutSuffix(result, []byte(suffix))
		if !found {
			return []byte{}
		}
	}
	return bytes.TrimLeft(execResult, "\r\n ")
}

func (session *Session) SessionInfo() string {
	remoteAddr := session.Conn.RemoteAddr().String()
	isAlive := "true"
	if !session.IsAlive {
		isAlive = "false"
	}
	return fmt.Sprintf("host: %s\talive: %s", remoteAddr, isAlive)
}

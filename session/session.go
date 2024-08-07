package session

import (
	"bytes"
	"fmt"
	"github.com/KPGhat/ShellSession/utils"
	"net"
	"sync"
	"time"
)

// Session结构体
type Session struct {
	Conn      net.Conn
	IsAlive   bool
	Id        int
	IsEcho    bool
	readLock  *sync.Mutex
	writeLock *sync.Mutex
	//Buffer []byte
}

// Send data to session
func (session *Session) Send(data string) {
	session.writeLock.Lock()
	defer session.writeLock.Unlock()

	_, err := session.Conn.Write([]byte(data))
	if err != nil {
		utils.Error(fmt.Sprintf("Send data to session error: %v", err))
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

func (session *Session) ReadUntil(suffix string) ([]byte, bool) {
	buffer := make([]byte, 1)
	var isTimeout bool
	var data bytes.Buffer

	for {
		session.Conn.SetReadDeadline(time.Now().Add(time.Second))
		n, err := session.Read(buffer)

		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				//utils.Error(fmt.Sprintf("Read data timeout: %v", err))
				isTimeout = true
			} else {
				session.IsAlive = false
				isTimeout = false
			}
			break
		}
		data.Write(buffer[:n])

		if bytes.HasSuffix(data.Bytes(), []byte(suffix)) {
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
			utils.Error(fmt.Sprintf("Read data to session error: %v", err))
			*running = false
			session.IsAlive = false
		}

		if n > 0 {
			callback(data[:n])
		}
	}
}

func (session *Session) ExecCmd(command string) []byte {
	prefix := utils.RandString(8)
	suffix := utils.RandString(8)
	newCommand := " echo " + prefix + " && " + command + "; echo " + suffix + "\n"
	session.Send(newCommand)

	result, isTimeOut := session.ReadUntil(suffix)
	if isTimeOut {
		return []byte{}
	}

	if session.IsEcho {
		_, nonexistent := session.ReadUntil(prefix)
		if !nonexistent {
			result, _ = session.ReadUntil(suffix)
		}
	}

	splitPrefix := bytes.Split(result, []byte(prefix))
	result = splitPrefix[len(splitPrefix)-1]
	result = bytes.Split(result, []byte(suffix))[0]
	return bytes.TrimLeft(result, "\r\n")
}

func (session *Session) SessionInfo() string {
	remoteAddr := session.Conn.RemoteAddr().String()
	isAlive := "true"
	if !session.IsAlive {
		isAlive = "false"
	}
	return fmt.Sprintf("host: %s\talive: %s", remoteAddr, isAlive)
}

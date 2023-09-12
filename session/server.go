package session

import (
	"fmt"
	"github.com/KPGhat/ShellSession/cmd"
	"log"
	"net"
	"time"
)

func handleSession(conn net.Conn) {
	sessionManager := GetSessionManager()
	sessionManager.AddSession(conn)
}

func StarServer() {
	address := fmt.Sprintf("%s:%d", cmd.Config.Host, cmd.Config.Port)
	shellListener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer StopServer(shellListener)

	go func() {
		for {
			time.Sleep(5 * time.Second)
			GetSessionManager().KeepAliveConn()
		}
	}()

	sem := make(chan struct{}, cmd.Config.MaxConn)
	for {
		conn, _ := shellListener.Accept()

		sem <- struct{}{}
		go func() {
			handleSession(conn)
			<-sem
		}()
	}
}

func StopServer(listener net.Listener) {
	err := listener.Close()
	if err != nil {
		log.Fatalf("%v", err)
	}
	sessionManager := GetSessionManager()
	sessionManager.DestroySessionManager()
}

package session

import (
	"fmt"
	"log"
	"net"
)

func handleSession(conn net.Conn) {
	//TODO
}

func StarServer(host string, port int, max int) error {
	address := fmt.Sprintf("%s:%d", host, port)
	shellListener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer func(shellListener net.Listener) {
		err = shellListener.Close()
		if err != nil {
			log.Fatalf("%v", err)
		}
	}(shellListener)

	sem := make(chan struct{}, max)
	for {
		conn, _ := shellListener.Accept()

		sem <- struct{}{}
		go func() {
			handleSession(conn)
			<-sem
		}()
	}

}

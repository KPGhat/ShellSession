package cli

import (
	"bufio"
	"github.com/KPGhat/ShellSession/session"
	"io"
	"os"
	"strings"
	"time"
)

func interact(session *session.Session, output io.Writer) {
	isInteract := true
	go func() {
		for isInteract {
			session.Conn.SetReadDeadline(time.Time{})
			data := session.Read()
			_, err := output.Write(data)
			if err != nil {
				isInteract = false
				return
			}
		}
	}()

	for isInteract {
		inputReader := bufio.NewReader(os.Stdin)
		command, _ := inputReader.ReadString('\n')
		command = strings.TrimSpace(command)
		if command == "interact exit" {
			isInteract = false
			break
		}
		session.Send([]byte(command + "\n"))
	}
}

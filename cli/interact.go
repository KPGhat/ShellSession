package cli

import (
	"bufio"
	"github.com/KPGhat/ShellSession/session"
	"io"
	"strings"
)

func interact(session *session.Session, input io.Reader, output io.Writer) {
	isInteract := true
	go session.ReadListener(&isInteract, func(data []byte) {
		_, err := output.Write(data)
		if err != nil {
			isInteract = false
			return
		}
	})

	for isInteract {
		inputReader := bufio.NewReader(input)
		command, _ := inputReader.ReadString('\n')
		command = strings.TrimSpace(command)
		if command == "bg sh" {
			isInteract = false
			break
		}
		session.Send([]byte(" " + command + "\n"))
	}
}

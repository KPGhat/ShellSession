package cli

import (
	"github.com/KPGhat/ShellSession/session"
	"log"
	"os"
	"strconv"
	"strings"
)

type cliType int

const (
	SESSION cliType = iota
	EXIT
	NOTEXIST
)

func dispatch(cmd string) cliType {
	cmd = strings.TrimSpace(cmd)
	cmdSplit := strings.Split(cmd, " ")
	sessionManager := session.GetSessionManager()
	if cmdSplit[0] == "session" {
		switch cmdSplit[1] {
		case "-l":
			sessionManager.ListAllSession(os.Stdout)
		case "-i":
			handleInteract(cmdSplit)
		}
		return SESSION
	} else if cmdSplit[0] == "exit" {
		log.Println("[+]Exiting program")
		return EXIT
	}

	return NOTEXIST
}

func handleInteract(cmd []string) {
	if len(cmd) != 3 {
		log.Println("[-]Session interact error.\nExample:\tsession -i id")
		return
	}

	sessionid, err := strconv.Atoi(cmd[2])
	if err != nil {
		log.Println("[-]Session interact error.\nExample:\tsession -i id\n[-]id is not a number")
		return
	}

	sess := session.GetSessionManager().GetSession(sessionid)
	interact(sess, os.Stdout)

}

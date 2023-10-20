package cli

import (
	"fmt"
	"github.com/KPGhat/ShellSession/session"
	"github.com/KPGhat/ShellSession/utils"
	"os"
	"strconv"
	"strings"
)

type cliType int

const (
	SESSION cliType = iota
	CONTEXT
	EXIT
	NOTEXIST
)

func dispatch(cmd string) cliType {
	cmd = strings.TrimSpace(cmd)
	cmdSplit := strings.Split(cmd, " ")
	sessionManager := session.GetManager()
	if cmdSplit[0] == "session" || cmdSplit[0] == "sess" {
		switch cmdSplit[1] {
		case "-l":
			if len(cmdSplit) == 3 && cmdSplit[2] == "all" {
				sessionManager.ListAllSession(os.Stdout, false)
			} else {
				sessionManager.ListAllSession(os.Stdout, true)
			}
		case "-i":
			handleInteract(cmdSplit)
		case "-a":
			handleForAllSession(cmdSplit)
		}
		return SESSION
	} else if cmdSplit[0] == "context" || cmdSplit[0] == "ctx" {
		switch cmdSplit[1] {
		case "-c":
			contextID := session.GetManager().CreateContext()
			handleContext(contextID)
		case "-i":
			dispatchContext(cmdSplit)
		case "-l":
			session.GetManager().ListAllContext(os.Stdout)
		}
		return CONTEXT
	} else if cmdSplit[0] == "exit" {
		utils.Congrats(fmt.Sprintf("Exiting program"))
		return EXIT
	}

	return NOTEXIST
}

func handleInteract(cmd []string) {
	if len(cmd) != 3 {
		utils.Warning("Session interact error.\nExample:\tsession -i id")
		return
	}

	sessionID, err := strconv.Atoi(cmd[2])
	if err != nil {
		utils.Warning("id is not a number")
		return
	}

	sess := session.GetManager().GetSession(sessionID)
	interact(sess, os.Stdin, os.Stdout)

}

func handleForAllSession(command []string) {
	execCmd := strings.Join(command[2:], " ")
	session.GetManager().ExecCmdForAll(execCmd, os.Stdout)
}

func dispatchContext(cmdSplit []string) {
	if len(cmdSplit) != 3 {
		utils.Warning("Wrong context command format\nExample:\tcontext -i [id]")
		return
	}
	id, err := strconv.Atoi(cmdSplit[2])
	if err != nil {
		utils.Warning("Error enter context: id must be a int")
		return
	}
	err = handleContext(id)
	if err != nil {
		utils.Warning(fmt.Sprintf("Error enter context: %v", err))
		return
	}
}

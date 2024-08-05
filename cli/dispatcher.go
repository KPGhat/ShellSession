package cli

import (
	"fmt"
	"github.com/KPGhat/ShellSession/cmd"
	"github.com/KPGhat/ShellSession/session"
	"github.com/KPGhat/ShellSession/utils"
	"os"
	"strconv"
	"strings"
)

type cliType int

const (
	HANDLEABLE cliType = iota
	EXIT
	NOTEXIST
)

var (
	cmdHandleMap map[string]func([]string)
)

func init() {
	cmdHandleMap = map[string]func([]string){
		"session": handleSession,
		"context": handleContext,
		"clear":   handleClear,
		"log":     handleLog,
	}
}

func dispatch(cmdSplit []string) cliType {
	// Use short name for convenience
	if cmdSplit[0] == "sess" || cmdSplit[0] == "s" {
		cmdSplit[0] = "session"
	} else if cmdSplit[0] == "ctx" || cmdSplit[0] == "c" {
		cmdSplit[0] = "context"
	} else if cmdSplit[0] == "exit" {
		utils.Congrats("Exiting program")
		return EXIT
	}

	if handleFunc, ok := cmdHandleMap[cmdSplit[0]]; ok {
		handleFunc(cmdSplit[1:])
		return HANDLEABLE
	}

	return NOTEXIST
}

func handleSession(args []string) {
	if len(args) < 1 {
		utils.Error("Missing session args")
		return
	}
	sessionManager := session.GetManager()
	switch args[0] {
	case "-l":
		if len(args) == 2 && args[1] == "all" {
			sessionManager.ListAllSession(os.Stdout, false)
		} else {
			sessionManager.ListAllSession(os.Stdout, true)
		}
	case "-i":
		handleInteract(args[1:])
	case "-a":
		execCmd := strings.Join(args[1:], " ")
		session.GetManager().ExecCmdForAll(execCmd, os.Stdout)
	}
}

func handleContext(args []string) {
	if len(args) < 1 {
		utils.Error("Missing context args")
		return
	}

	switch args[0] {
	case "-c":
		contextID := session.GetManager().CreateContext()
		enterContext(contextID)
	case "-i":
		dispatchContext(args[1:])
	case "-l":
		session.GetManager().ListAllContext(os.Stdout)
	}
}

func handleClear(args []string) {
	option := "default"
	if len(args) >= 1 {
		option = args[0]
	}
	switch option {
	case "-a":
		session.GetManager().HandleAllSession(func(s *session.Session) {
			if !s.IsAlive {
				session.GetManager().DelSession(s.Id)
				utils.Congrats("Clear not alive session <" + strconv.Itoa(s.Id) + ">")
				return
			}
			randstr := utils.RandString(16)
			result := string(s.ExecCmd("echo " + randstr))
			result = strings.Trim(result, "\n\r")
			if !(result == randstr) {
				session.GetManager().DelSession(s.Id)
				utils.Congrats("Clear closed session <" + strconv.Itoa(s.Id) + ">")
			}
		})
	default:
		session.GetManager().HandleAllSession(func(s *session.Session) {
			if !s.IsAlive {
				session.GetManager().DelSession(s.Id)
				utils.Congrats("Clear not alive session <" + strconv.Itoa(s.Id) + ">")
			}
		})
	}

}

func handleInteract(args []string) {
	if len(args) != 1 {
		utils.Error("Session interact error.\nExample:\tsession -i id")
		return
	}

	sessionID, err := strconv.Atoi(args[0])
	if err != nil {
		utils.Error("id is not a number")
		return
	}

	sess := session.GetManager().GetSession(sessionID)
	interact(sess, os.Stdin, os.Stdout)

}

func handleLog(args []string) {
	if len(args) != 1 {
		utils.Error("Missing log arg")
		return
	}

	switch args[0] {
	case "on":
		cmd.Config.LogOff = false
		utils.Congrats("Log ON")
	case "off":
		cmd.Config.LogOff = true
		utils.Congrats("Log OFF")
	}
}

func dispatchContext(args []string) {
	if len(args) != 1 {
		utils.Error("Wrong context command format\nExample:\tcontext -i [id]")
		return
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		utils.Error("Error enter context: id must be a int")
		return
	}
	err = enterContext(id)
	if err != nil {
		utils.Error(fmt.Sprintf("Error enter context: %v", err))
		return
	}
}

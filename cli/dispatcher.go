package cli

import (
	"bufio"
	"fmt"
	"github.com/KPGhat/ShellSession/plugin"
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
		case "-m":
			handleManager()
		case "-a":
			handleForAllSession(cmdSplit)
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
		fmt.Println("[-]Session interact error.\nExample:\tsession -i id")
		return
	}

	sessionid, err := strconv.Atoi(cmd[2])
	if err != nil {
		fmt.Println("[-]Session interact error.\nExample:\tsession -i id\n[-]id is not a number")
		return
	}

	sess := session.GetSessionManager().GetSession(sessionid)
	interact(sess, os.Stdin, os.Stdout)

}

func handleForAllSession(command []string) {
	execCmd := strings.Join(command[2:], " ")
	session.GetSessionManager().ExecCmdForAll(execCmd, os.Stdout)
}

func handleManager() {
	for {
		prompt := "context[" + session.GetSessionManager().GetAllContext() + "]>"
		os.Stdout.Write([]byte(prompt))

		reader := bufio.NewReader(os.Stdin)
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)
		cmdSplit := strings.Split(cmd, " ")
		if cmdSplit[0] == "add" {
			if len(cmdSplit) < 2 {
				fmt.Println("[-]Session manage add error.\nExample:\tadd id [id...]\n\tadd all")
			}

			if cmdSplit[1] == "all" {
				session.GetSessionManager().AddAllContext()
			} else {
				for _, idStr := range cmdSplit[1:] {
					id, err := strconv.Atoi(idStr)
					if err != nil {
						fmt.Println("[-]Session manage add error: id is a number")
						break
					}
					err = session.GetSessionManager().AddContext(id)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		} else if cmdSplit[0] == "del" {
			if len(cmdSplit) < 2 {
				fmt.Println("[-]Session manage delete error.\nExample:\tdel id [id...]\n\tdel all")
			}

			if cmdSplit[1] == "all" {
				session.GetSessionManager().DelAllContext()
			} else {
				for _, idStr := range cmdSplit[1:] {
					id, err := strconv.Atoi(idStr)
					if err != nil {
						fmt.Println("[-]Session manage add error: id is a number")
						break
					}
					err = session.GetSessionManager().DelContext(id)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		} else if cmdSplit[0] == "sh" {
			if len(cmdSplit) < 2 {
				fmt.Println("[-]Session manage execute shell error.\nExample:\tsh cmd")
			}
			session.GetSessionManager().HandleAllContext(func(session *session.Session) {
				result := session.ExecCmd([]byte(strings.Join(cmdSplit[1:], " ")))
				os.Stdout.Write(result)
			})
		} else if cmdSplit[0] == "upload" {
			session.GetSessionManager().HandleAllContext(func(session *session.Session) {
				execResult := plugin.Plugin.Upload(session, cmdSplit[1:])
				os.Stdout.Write([]byte(execResult))
			})
		} else if cmdSplit[0] == "exit" {
			fmt.Println("[+]Exiting context manage...")
			session.GetSessionManager().DelAllContext()
			break
		}
	}
}

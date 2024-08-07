package cli

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/KPGhat/ShellSession/utils"
	"github.com/fatih/color"
	"os"
	"strconv"
	"strings"

	"github.com/KPGhat/ShellSession/plugin"
	"github.com/KPGhat/ShellSession/session"
)

func enterContext(id int) error {
	context := session.GetManager().GetContext(id)
	if context == nil {
		return errors.New("No such context")
	}
	for {
		blue := color.New(color.FgBlue).SprintFunc()
		cyan := color.New(color.FgHiCyan).SprintFunc()
		prompt := blue("context[") + cyan("size:"+context.Size()) + blue("]>")
		os.Stdout.Write([]byte(prompt))

		reader := bufio.NewReader(os.Stdin)
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)
		cmdSplit := strings.Split(cmd, " ")
		if cmdSplit[0] == "add" || cmdSplit[0] == "del" {
			operateContext(context, cmdSplit)
		} else if cmdSplit[0] == "list" {
			context.List(os.Stdout)
		} else if cmdSplit[0] == "sh" {
			if len(cmdSplit) < 2 {
				utils.Error("Session manage execute shell error.\nExample:\tsh cmd")
				continue
			}
			context.HandleAllContext(func(session *session.Session) {
				result := session.ExecCmd(strings.Join(cmdSplit[1:], " "))
				os.Stdout.Write(result)
			})
		} else if cmdSplit[0] == "upload" {
			context.HandleAllContext(func(session *session.Session) {
				execResult := plugin.Plugin.Upload(session, cmdSplit[1:])
				os.Stdout.Write([]byte(execResult))
			})
		} else if cmdSplit[0] == "exit" {
			utils.Congrats("Exiting context manage...")
			break
		}
	}

	return nil
}

func operateContext(context *session.Context, cmdSplit []string) {
	if len(cmdSplit) < 2 {
		utils.Error(fmt.Sprintf("Session manage %s error.\nExample:\t%s id [id...]\n\t%s all\n", cmdSplit[0], cmdSplit[0], cmdSplit[0]))
		return
	}

	var actionFunc func(id int) error
	var actionForAll func()
	if cmdSplit[0] == "add" {
		actionFunc = context.Add
		actionForAll = context.AddAll
		if cmdSplit[1] == "allip" {
			actionForAll = context.AddAllIP
		}
	} else if cmdSplit[0] == "del" {
		actionFunc = context.Delete
		actionForAll = context.DeleteAll
	}

	if strings.Contains(cmdSplit[1], "all") {
		actionForAll()
	} else {
		for _, idStr := range cmdSplit[1:] {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				utils.Error("Session manage add error: id is a number")
				break
			}
			err = actionFunc(id)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

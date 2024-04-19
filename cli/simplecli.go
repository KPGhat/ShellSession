package cli

import (
	"bufio"
	"github.com/KPGhat/ShellSession/utils"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

func CliControl() {
	running := true
	exitStatue := false
	defer func() {
		if err := recover(); err != nil {
			log.Println("Panic: ", err)
			CliControl()
		}
	}()
	for running {
		prompt := "gsh>"
		color.New(color.FgMagenta).Fprint(os.Stdout, prompt)

		reader := bufio.NewReader(os.Stdin)
		cmd, _ := reader.ReadString('\n')
		cmd = strings.Trim(cmd, "\r\x20\n")

		if cmd == "" {
			continue
		}
		cmdSplit := strings.Split(cmd, " ")
		cmdType := dispatch(cmdSplit)

		switch cmdType {
		case EXIT:
			// enter twice to exit
			if exitStatue {
				running = false
			} else {
				exitStatue = true
				utils.Congrats("Please enter exit again")
			}
		case NOTEXIST:
			utils.Warning("gsh: " + cmdSplit[0] + ": no such command")
		}
	}

}

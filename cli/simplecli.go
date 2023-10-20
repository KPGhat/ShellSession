package cli

import (
	"bufio"
	"log"
	"os"

	"github.com/fatih/color"
)

func CliControl() {
	running := true
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
		cmdType := dispatch(cmd)
		if cmdType == EXIT {
			running = false
		}
	}

}

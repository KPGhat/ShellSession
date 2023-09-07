package cli

import (
	"bufio"
	"os"
)

func CliControl() {
	running := true
	for running {
		prompt := "->"
		os.Stdout.Write([]byte(prompt))

		reader := bufio.NewReader(os.Stdin)
		cmd, _ := reader.ReadString('\n')
		cmdType := dispatch(cmd)
		if cmdType == EXIT {
			running = false
		}
	}

}

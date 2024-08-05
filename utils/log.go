package utils

import (
	"fmt"
	"github.com/KPGhat/ShellSession/cmd"
	"github.com/fatih/color"
)

func Congrats(message string) {
	if !cmd.Config.LogOff {
		green := color.New(color.FgHiGreen).SprintFunc()
		fmt.Println("\r" + green("[+]") + message)
	}
}

func Error(message string) {
	if !cmd.Config.LogOff {
		red := color.New(color.FgHiRed).SprintFunc()
		fmt.Println("\r" + red("[-]") + message)
	}
}

func Warning(message string) {
	if !cmd.Config.LogOff {
		yellow := color.New(color.FgHiYellow).SprintFunc()
		fmt.Println("\r" + yellow("[*]") + message)
	}
}

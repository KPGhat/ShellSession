package utils

import (
	"fmt"
	"github.com/fatih/color"
)

func Congrats(message string) {
	green := color.New(color.FgHiGreen).SprintFunc()
	fmt.Println(green("[+]") + message)
}

func Warning(message string) {
	red := color.New(color.FgHiRed).SprintFunc()
	fmt.Println(red("[-]") + message)
}

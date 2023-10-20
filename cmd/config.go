package cmd

import (
	"fmt"
	"github.com/fatih/color"
)

type config struct {
	Host   string
	Port   int
	LogOff bool
}

var Config config

func PrintConfig() {
	green := color.New(color.FgHiGreen).SprintFunc()
	fmt.Printf("%sListening at %s:%d\n", green("[+]"), Config.Host, Config.Port)
}

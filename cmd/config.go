package cmd

import (
	"fmt"
	"github.com/KPGhat/ShellSession/utils"
)

type config struct {
	Host string
	Port int
}

var Config config

func PrintConfig() {
	utils.Congrats(fmt.Sprintf("Listening at %s:%d", Config.Host, Config.Port))
}

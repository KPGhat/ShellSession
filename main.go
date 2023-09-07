package main

import (
	"github.com/KPGhat/ShellSession/cli"
	"github.com/KPGhat/ShellSession/cmd"
	"github.com/KPGhat/ShellSession/session"
	"os"
	"os/signal"
)

func main() {
	cmd.Flag()

	// Capture Ctrl-C Signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go session.StarServer()
	cli.CliControl()
}

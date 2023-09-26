package cmd

import "fmt"

type config struct {
	Host    string
	Port    int
	MaxConn int
}

var Config config

func PrintConfig() {
	fmt.Printf("[+]Listening at %s:%d and the maximum number of connections is %d\n", Config.Host, Config.Port, Config.MaxConn)
}

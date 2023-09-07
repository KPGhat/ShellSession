package cmd

type config struct {
	Host    string
	Port    int
	MaxConn int
}

var Config config

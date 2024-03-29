package cmd

import (
	"flag"
	"os"
)

func Flag() {
	flag.StringVar(&Config.Host, "host", "0.0.0.0", "The listen host")
	flag.IntVar(&Config.Port, "port", 2333, "The listen port")
	flag.BoolVar(&Config.LogOff, "nolog", false, "Turn off log")

	var help bool
	flag.BoolVar(&help, "h", false, "Print this help info")

	flag.Parse()
	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}
}

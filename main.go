package main

import (
	"flag"
	"os"
	"regexp"

	"github.com/teambition/gear/logging"
	"github.com/zensh/ipservice/src"
)

var (
	portReg  = regexp.MustCompile(`^\d+$`)
	port     = flag.String("port", "8080", `Server port.`)
	dataPath = flag.String("data", "", "IP data file path.")
)

func main() {
	flag.Parse()
	if portReg.MatchString(*port) {
		*port = ":" + *port
	}
	if *port == "" || *dataPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	app := src.New(*dataPath)
	// start app
	logging.Info("IP Service start " + *port)
	app.Listen(*port)
}

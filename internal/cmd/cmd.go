package cmd

import (
	"flag"
)

var (
	// run as server?
	serverMode = false
)

// parseCommandLine parses the command line arguments
func parseCommandLine() {
	// configure command line arguments
	flag.BoolVar(&serverMode, "server", serverMode, "run as server")

	// parse command line arguments
	flag.Parse()

}

// Run is the main entry point
func Run() {
	parseCommandLine()

	if serverMode {
		newServer().run()
		return
	}
}

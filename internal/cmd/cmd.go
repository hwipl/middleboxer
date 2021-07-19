package cmd

import (
	"flag"
)

var (
	// serverMode determines if we run as a server or a client
	serverMode = false
)

// parseCommandLine parses the command line arguments
func parseCommandLine() {
	// configure command line arguments
	flag.BoolVar(&serverMode, "server", serverMode,
		"run as server (default: run as client)")

	// parse command line arguments
	flag.Parse()

}

// Run is the main entry point
func Run() {
	parseCommandLine()

	// run as server?
	if serverMode {
		newServer("").run()
		return
	}

	// run as client
	newClient("").run()
}

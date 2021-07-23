package cmd

import (
	"flag"
)

var (
	// serverMode determines if we run as a server or a client
	serverMode = false

	// serverAddress is the address of the server
	serverAddress = ""

	// clientId is the id of the client
	clientId uint8 = 0
)

// parseCommandLine parses the command line arguments
func parseCommandLine() {
	// configure command line arguments
	flag.BoolVar(&serverMode, "server", serverMode,
		"run as server (default: run as client)")
	flag.StringVar(&serverAddress, "address", serverAddress,
		"set address to connect to (client mode) or listen on (server mode)")

	// parse command line arguments
	flag.Parse()

}

// Run is the main entry point
func Run() {
	parseCommandLine()

	// run as server?
	if serverMode {
		newServer(serverAddress).run()
		return
	}

	// run as client
	newClient(serverAddress, clientId).run()
}

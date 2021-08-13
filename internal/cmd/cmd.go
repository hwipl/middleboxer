package cmd

import (
	"flag"
	"log"
	"math"
)

var (
	// serverMode determines if we run as a server or a client
	serverMode = false

	// serverAddress is the address of the server
	serverAddress = ""

	// clientID is the id of the client
	clientID uint8 = 0
)

// parseCommandLine parses the command line arguments
func parseCommandLine() {
	// configure command line arguments
	flag.BoolVar(&serverMode, "server", serverMode,
		"run as server (default: run as client)")
	flag.StringVar(&serverAddress, "address", serverAddress,
		"set address to connect to (client mode) or listen on (server mode)")
	cid := flag.Uint("id", 0, "set id of the client")

	// parse command line arguments
	flag.Parse()

	// set client id
	if *cid > math.MaxUint8 {
		log.Fatal("invalid client id")
	}
	clientID = uint8(*cid)
}

// Run is the main entry point
func Run() {
	parseCommandLine()

	// run as server?
	if serverMode {
		// TODO: get sender and receiver ID from command line
		plan := newPlan(0, 1)
		newServer(serverAddress, plan).run()
		return
	}

	// run as client
	newClient(serverAddress, clientID).run()
}

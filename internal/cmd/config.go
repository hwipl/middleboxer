package cmd

import (
	"flag"
	"log"
	"math"
)

// Config contains the configuration
type Config struct {
	// ServerMode determines if we run as a server or a client
	ServerMode bool

	// ServerAddress is the address of the server
	ServerAddress string

	// ClientID is the id of the client
	ClientID uint8
}

// ParseCommandLine fills the config from command line arguments
func (c *Config) ParseCommandLine() {
	// configure command line arguments
	flag.BoolVar(&c.ServerMode, "server", c.ServerMode,
		"run as server (default: run as client)")
	flag.StringVar(&c.ServerAddress, "address", c.ServerAddress,
		"set address to connect to (client mode) or listen on (server mode)")
	cid := flag.Uint("id", 0, "set id of the client")

	// parse command line arguments
	flag.Parse()

	// set client id
	if *cid > math.MaxUint8 {
		log.Fatal("invalid client id")
	}
	c.ClientID = uint8(*cid)
}

// NewConfig creates a new Config
func NewConfig() *Config {
	return &Config{}
}

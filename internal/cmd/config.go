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

	// SenderID is the id of the sending client
	SenderID uint8
}

// ParseCommandLine fills the config from command line arguments
func (c *Config) ParseCommandLine() {
	// configure command line arguments
	flag.BoolVar(&c.ServerMode, "server", c.ServerMode,
		"run as server (default: run as client)")
	flag.StringVar(&c.ServerAddress, "address", c.ServerAddress,
		"set address to connect to (client mode) or listen on (server mode)")
	cid := flag.Uint("id", 0, "set id of the client")
	sid := flag.Uint("sid", 0, "set id of the sending client")

	// parse command line arguments
	flag.Parse()

	// set client id, sender id
	for _, i := range []*uint{cid, sid} {
		if *i > math.MaxUint8 {
			log.Fatal("invalid client id: ", *i)
		}
	}
	c.ClientID = uint8(*cid)
	c.SenderID = uint8(*sid)
}

// NewConfig creates a new Config
func NewConfig() *Config {
	return &Config{}
}

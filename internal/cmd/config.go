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

	// ReceiverID is the id of the receiving client
	ReceiverID uint8

	// SenderDevice is the name of the sender's network interface
	SenderDevice string

	// ReceiverDevice is the name of the receiver's network interface
	ReceiverDevice string

	// SenderSrcMAC is the sender's source MAC address
	SenderSrcMAC string

	// SenderDstMAC is the sender's destination MAC address
	SenderDstMAC string

	// ReceiverSrcMAC is the receiver's source MAC address
	ReceiverSrcMAC string

	// ReceiverDstMAC is the receiver's destination MAC address
	ReceiverDstMAC string

	// SenderSrcIP is the sender's source IP address
	SenderSrcIP string

	// SenderDstIP is the sender's destination IP address
	SenderDstIP string

	// ReceiverSrcIP is the receiver's source IP address
	ReceiverSrcIP string

	// ReceiverDstIP is the receiver's destination IP address
	ReceiverDstIP string

	// Protocol is the layer 4 protocol
	Protocol uint16

	// SenderSrcPort is the sender's source port
	SenderSrcPort uint16
}

// ParseCommandLine fills the config from command line arguments
func (c *Config) ParseCommandLine() {
	// configure command line arguments
	flag.BoolVar(&c.ServerMode, "server", c.ServerMode,
		"run as server (default: run as client)")
	flag.StringVar(&c.ServerAddress, "address", c.ServerAddress,
		"set address to connect to (client mode) or listen on (server mode)")
	cid := flag.Uint("id", uint(c.ClientID), "set id of the client")
	sid := flag.Uint("sid", uint(c.SenderID), "set id of the sending client")
	rid := flag.Uint("rid", uint(c.ReceiverID), "set id of the receiving client")

	// parse command line arguments
	flag.Parse()

	// set client id, sender id, receiver id
	for _, i := range []*uint{cid, sid, rid} {
		if *i > math.MaxUint8 {
			log.Fatal("invalid client id: ", *i)
		}
	}
	c.ClientID = uint8(*cid)
	c.SenderID = uint8(*sid)
	c.ReceiverID = uint8(*rid)
}

// NewConfig creates a new Config
func NewConfig() *Config {
	return &Config{
		ClientID:   1,
		SenderID:   1,
		ReceiverID: 2,
	}
}

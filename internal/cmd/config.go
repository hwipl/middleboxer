package cmd

import (
	"flag"
	"log"
	"math"
	"net"
	"strconv"
	"strings"
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

	// PortRange is the tested port range
	PortRange string
}

// getMACFromString converts a string to a hardware (MAC) address
func getMACFromString(mac string) net.HardwareAddr {
	m, err := net.ParseMAC(mac)
	if err != nil {
		return nil
	}
	return m
}

// GetSenderSrcMAC returns the sender's source MAC address
func (c *Config) GetSenderSrcMAC() net.HardwareAddr {
	return getMACFromString(c.SenderSrcMAC)
}

// GetSenderDstMAC returns the sender's destination MAC address
func (c *Config) GetSenderDstMAC() net.HardwareAddr {
	return getMACFromString(c.SenderDstMAC)
}

// GetReceiverSrcMAC returns the receiver's source MAC address
func (c *Config) GetReceiverSrcMAC() net.HardwareAddr {
	return getMACFromString(c.ReceiverSrcMAC)
}

// GetReceiverDstMAC returns the receiver's destination MAC address
func (c *Config) GetReceiverDstMAC() net.HardwareAddr {
	return getMACFromString(c.ReceiverDstMAC)
}

// getIPFromString converts a string to an IP address
func getIPFromString(ip string) net.IP {
	return net.ParseIP(ip)
}

// GetSenderSrcIP returns the sender's source IP address
func (c *Config) GetSenderSrcIP() net.IP {
	return getIPFromString(c.SenderSrcIP)
}

// GetSenderDstIP returns the sender's destination IP address
func (c *Config) GetSenderDstIP() net.IP {
	return getIPFromString(c.SenderDstIP)
}

// GetReceiverSrcIP returns the receiver's source IP address
func (c *Config) GetReceiverSrcIP() net.IP {
	return getIPFromString(c.ReceiverSrcIP)
}

// GetReceiverDstIP returns the receiver's destination IP address
func (c *Config) GetReceiverDstIP() net.IP {
	return getIPFromString(c.ReceiverDstIP)
}

// GetPortRange returns the first and last port of the port range
func (c *Config) GetPortRange() (first uint16, last uint16) {
	// get first and last port as string
	fs, ls := "", ""
	s := strings.Split(c.PortRange, ":")
	switch len(s) {
	case 1:
		fs = s[0]
		ls = s[0]
	case 2:
		fs = s[0]
		ls = s[1]
	default:
		return
	}

	// parse first port string
	f, err := strconv.ParseUint(fs, 10, 16)
	if err != nil {
		return
	}

	// parse last port string
	l, err := strconv.ParseUint(ls, 10, 16)
	if err != nil {
		return
	}

	// return first and last port
	first = uint16(f)
	last = uint16(l)
	return
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
	flag.StringVar(&c.SenderDevice, "sdev", c.SenderDevice,
		"set device of the sending client")
	flag.StringVar(&c.ReceiverDevice, "rdev", c.ReceiverDevice,
		"set device of the receiving client")
	flag.StringVar(&c.SenderSrcMAC, "ssmac", c.SenderSrcMAC,
		"set source MAC of the sending client")
	flag.StringVar(&c.SenderDstMAC, "sdmac", c.SenderDstMAC,
		"set destination MAC of the sending client")
	flag.StringVar(&c.ReceiverSrcMAC, "rsmac", c.ReceiverSrcMAC,
		"set source MAC of the receiving client")
	flag.StringVar(&c.ReceiverDstMAC, "rdmac", c.ReceiverDstMAC,
		"set destination MAC of the receiving client")
	flag.StringVar(&c.SenderSrcIP, "ssip", c.SenderSrcIP,
		"set source IP of the sending client")
	flag.StringVar(&c.SenderDstIP, "sdip", c.SenderDstIP,
		"set destination IP of the sending client")
	flag.StringVar(&c.ReceiverSrcIP, "rsip", c.ReceiverSrcIP,
		"set source IP of the receiving client")
	flag.StringVar(&c.ReceiverDstIP, "rdip", c.ReceiverDstIP,
		"set destination IP of the receiving client")

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

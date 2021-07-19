package cmd

import (
	"log"
	"net"
)

// client stores information about a client
type client struct {
	conn net.Conn
}

// run runs this client
func (c *client) run() {
	defer func() {
		_ = c.conn.Close()
	}()
}

// newClient connects to serverAddress and creates a new client
func newClient(serverAddress string) *client {
	// create connection to server
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.Fatal(err)
	}

	// return client
	return &client{
		conn,
	}
}

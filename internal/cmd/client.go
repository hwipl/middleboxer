package cmd

import (
	"log"
	"net"
)

// client stores information about a client
type client struct {
	conn net.Conn
	id   uint8
}

// run runs this client
func (c *client) run() {
	defer func() {
		_ = c.conn.Close()
	}()

	log.Println("Client connected to:", c.conn.RemoteAddr())
}

// newClient connects to serverAddress and creates a new client
func newClient(serverAddress string, id uint8) *client {
	// create connection to server
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.Fatal(err)
	}

	// return client
	return &client{
		conn,
		id,
	}
}

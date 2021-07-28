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

// registerClient registers this client on the server
func (c *client) registerClient() bool {
	reg := MessageRegister{c.id}
	return writeMessage(c.conn, &reg)
}

// sendNop sends a nop message to the server
func (c *client) sendNop() bool {
	nop := MessageNop{}
	return writeMessage(c.conn, &nop)
}

// run runs this client
func (c *client) run() {
	defer func() {
		_ = c.conn.Close()
	}()

	log.Println("Client connected to:", c.conn.RemoteAddr())
	if !c.registerClient() {
		return
	}
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

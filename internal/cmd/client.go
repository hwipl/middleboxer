package cmd

import (
	"log"
	"net"
	"time"
)

const (
	// NopInterval specifies the seconds between sending nop messages
	NopInterval = 15
)

// client stores information about a client
type client struct {
	conn    net.Conn
	id      uint8
	tests   chan *MessageTest
	results chan *MessageResult
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

// sendResult returns the result message to the server
func (c *client) sendResult(result *MessageResult) bool {
	return writeMessage(c.conn, result)
}

// receive gets messages from the server
func (c *client) receive() {
	defer close(c.tests)
	for {
		// read message from server
		msg := readMessage(c.conn)
		if msg == nil {
			return
		}

		// handle test command messages
		if msg.GetType() != MessageTypeTest {
			continue
		}
		test, ok := msg.(*MessageTest)
		if !ok {
			log.Println("Received invalid test message from server")
			continue
		}
		c.tests <- test
	}
}

// runTest runs the test requested by the server in the test message
func (c *client) runTest(test *MessageTest) {
	if test.Initiate {
		go newSender(test, c.results).run()
	} else {
		go newReceiver(test, c.results).run()
	}
}

// run runs this client
func (c *client) run() {
	defer func() {
		_ = c.conn.Close()
	}()

	log.Println("Client connected to:", c.conn.RemoteAddr())

	// register client
	if !c.registerClient() {
		return
	}
	log.Println("Client registered on server")

	// create ticker for nop messages
	ticker := time.NewTicker(time.Second * NopInterval)
	defer ticker.Stop()

	// start receiving messages
	go c.receive()

	log.Println("Client ready and waiting for test commands")
	for {
		select {
		case <-ticker.C:
			if !c.sendNop() {
				return
			}
		case test, more := <-c.tests:
			if !more {
				return
			}
			c.runTest(test)
		case result := <-c.results:
			if !c.sendResult(result) {
				return
			}
		}
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
		make(chan *MessageTest),
		make(chan *MessageResult),
	}
}

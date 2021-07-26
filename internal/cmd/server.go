package cmd

import (
	"log"
	"net"
)

// clientHandler handles a client connected to the server
type clientHandler struct {
	conn net.Conn
}

// newClientHandler creates a new client handler with conn
func newClientHandler(conn net.Conn) *clientHandler {
	return &clientHandler{
		conn,
	}
}

// registerClient reads a client registration from client and returns
// client id and ok
func (c *clientHandler) registerClient() (clientId uint8, ok bool) {
	// read message from client
	msg := readMessage(c.conn)
	if msg == nil {
		return
	}

	// handle register message
	if msg.GetType() != MessageTypeRegister {
		return
	}
	clientId = msg.(*MessageRegister).Client
	ok = true
	return
}

// handleClient handles a client connection
func (c *clientHandler) run() {
	defer func() {
		_ = c.conn.Close()
	}()

	log.Printf("Client %s connected", c.conn.RemoteAddr())

	// await client registration
	clientId, ok := c.registerClient()
	if !ok {
		return
	}
	log.Printf("Client %s registered with id %d", c.conn.RemoteAddr(),
		clientId)

	// enter main loop
	for {
		// read message from client
		msg := readMessage(c.conn)
		if msg == nil {
			break
		}

		// handle message based on type
		switch msg.GetType() {
		case MessageTypeNop:
		default:
			// invalid client message; disconnect client
			log.Println("Invalid message from client",
				c.conn.RemoteAddr())
			break
		}
	}
}

// server stores information about a server
type server struct {
	listener net.Listener
}

// run runs this server
func (s *server) run() {
	defer func() {
		_ = s.listener.Close()
	}()

	log.Println("Server listening on:", s.listener.Addr())
	for {
		client, err := s.listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go newClientHandler(client).run()
	}
}

// newServer creates an new server that listens on address
func newServer(address string) *server {
	// create listener
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	// return server
	return &server{
		listener,
	}
}

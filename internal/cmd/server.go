package cmd

import (
	"log"
	"net"
)

// clientResult is a result sent by a client
type clientResult struct {
	clientID uint8
	result   *MessageResult
}

// clientHandler handles a client connected to the server
type clientHandler struct {
	conn       net.Conn
	id         uint8
	clientRegs chan *clientHandler
}

// newClientHandler creates a new client handler with conn
func newClientHandler(conn net.Conn, clientRegs chan *clientHandler) *clientHandler {
	return &clientHandler{
		conn,
		0,
		clientRegs,
	}
}

// registerClient reads a client registration from client and returns ok
func (c *clientHandler) registerClient() bool {
	// read message from client
	msg := readMessage(c.conn)
	if msg == nil {
		return false
	}

	// handle register message
	if msg.GetType() != MessageTypeRegister {
		return false
	}
	c.id = msg.(*MessageRegister).Client
	return true
}

// handleClient handles a client connection
func (c *clientHandler) run() {
	defer func() {
		_ = c.conn.Close()
	}()

	log.Printf("Client %s connected", c.conn.RemoteAddr())

	// await client registration
	if ok := c.registerClient(); !ok {
		return
	}
	c.clientRegs <- c
	log.Printf("Client %s registered with id %d", c.conn.RemoteAddr(),
		c.id)

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
	listener   net.Listener
	clientRegs chan *clientHandler
	clients    map[uint8]*clientHandler
}

// listen waits for new connections from clients
func (s *server) listen() {
	defer func() {
		_ = s.listener.Close()
	}()

	log.Println("Server listening on:", s.listener.Addr())
	for {
		client, err := s.listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go newClientHandler(client, s.clientRegs).run()
	}
}

// run runs this server
func (s *server) run() {
	go s.listen()

	for {
		select {
		case c := <-s.clientRegs:
			s.clients[c.id] = c
		}
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
		make(chan *clientHandler),
		make(map[uint8]*clientHandler),
	}
}

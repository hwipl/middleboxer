package cmd

import (
	"log"
	"net"
)

// server stores information about a server
type server struct {
	listener net.Listener
}

// registerClient reads a client registration from client and returns
// client id and ok
func (s *server) registerClient(client net.Conn) (clientId uint8, ok bool) {
	// read message from client
	msg := readMessage(client)
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
func (s *server) handleClient(client net.Conn) {
	defer func() {
		_ = client.Close()
	}()

	log.Println("Client connected:", client.RemoteAddr())

	// await client registration
	clientId, ok := s.registerClient(client)
	if !ok {
		return
	}
	log.Println("Client registered with id", clientId)

	// enter main loop
	for {
		// read message from client
		msg := readMessage(client)
		if msg == nil {
			break
		}

		// handle message based on type
		switch msg.GetType() {
		case MessageTypeNop:
		default:
			// invalid client message; disconnect client
			log.Println("Invalid message from client")
			break
		}
	}
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
		go s.handleClient(client)
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

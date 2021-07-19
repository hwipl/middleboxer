package cmd

import (
	"log"
	"net"
)

// server stores information about a server
type server struct {
	listener net.Listener
}

// handleClient handles a client connection
func (s *server) handleClient(client net.Conn) {
}

// run runs this server
func (s *server) run() {
	for {
		client, err := s.listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go s.handleClient(client)
	}
}

// newServer creates an new server
func newServer() *server {
	// create listener
	listener, err := net.Listen("tcp", "")
	if err != nil {
		log.Fatal(err)
	}

	// return server
	return &server{
		listener,
	}
}

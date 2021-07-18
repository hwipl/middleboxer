package cmd

// server stores information about a server
type server struct {
}

// run this server
func (s *server) run() {
}

// newServer creates an new server
func newServer() *server {
	return &server{}
}

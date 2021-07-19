package cmd

// client stores information about a client
type client struct {
}

// run runs this client
func (c *client) run() {
}

// newClient creates a new client
func newClient() *client {
	return &client{}
}

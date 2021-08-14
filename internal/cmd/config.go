package cmd

// Config contains the configuration
type Config struct {
	// ServerMode determines if we run as a server or a client
	ServerMode bool

	// ServerAddress is the address of the server
	ServerAddress string

	// ClientID is the id of the client
	ClientID uint8
}

// NewConfig creates a new Config
func NewConfig() *Config {
	return &Config{}
}

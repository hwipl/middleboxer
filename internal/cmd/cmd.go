package cmd

// Run is the main entry point
func Run() {
	// create config
	config := NewConfig()
	config.ParseCommandLine()

	// run as server?
	if config.ServerMode {
		// TODO: get sender and receiver ID from command line
		plan := newPlan(0, 1)
		newServer(config.ServerAddress, plan).run()
		return
	}

	// run as client
	newClient(config.ServerAddress, config.ClientID).run()
}

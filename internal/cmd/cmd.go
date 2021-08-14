package cmd

// Run is the main entry point
func Run() {
	// create config
	config := NewConfig()
	config.ParseCommandLine()

	// run as server?
	if config.ServerMode {
		// TODO: get receiver ID from command line
		plan := newPlan(config.SenderID, 1)
		newServer(config.ServerAddress, plan).run()
		return
	}

	// run as client
	newClient(config.ServerAddress, config.ClientID).run()
}

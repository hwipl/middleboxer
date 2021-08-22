package cmd

// Run is the main entry point
func Run() {
	// create config
	config := NewConfig()
	config.ParseCommandLine()

	// run as server?
	if config.ServerMode {
		plan := newPlan(config)
		newServer(config.ServerAddress, plan).run()
		plan.printResults()
		return
	}

	// run as client
	newClient(config.ServerAddress, config.ClientID).run()
}

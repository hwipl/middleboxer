package cmd

// Run is the main entry point
func Run() {
	// create config
	config := NewConfig()
	config.ParseCommandLine()

	// run as server?
	if config.ServerMode {
		plan := newPlan(config.SenderID, config.ReceiverID)
		newServer(config.ServerAddress, plan).run()
		return
	}

	// run as client
	newClient(config.ServerAddress, config.ClientID).run()
}

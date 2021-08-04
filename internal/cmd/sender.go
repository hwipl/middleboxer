package cmd

// sender is a test in sender mode
type sender struct {
	test    *MessageTest
	results chan *MessageResult
}

// run runs the sender
func (s *sender) run() {
	// TODO: do something
}

// newSender creates a new test in sender mode
func newSender(test *MessageTest, results chan *MessageResult) *sender {
	return &sender{
		test,
		results,
	}
}

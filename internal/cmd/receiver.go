package cmd

// receiver is a test in receiver mode
type receiver struct {
	test    *MessageTest
	results chan *MessageResult
}

// run runs the receiver
func (r *receiver) run() {
	// TODO: do something
}

// newReceiver creates a new test in receiver mode
func newReceiver(test *MessageTest, results chan *MessageResult) *receiver {
	return &receiver{
		test,
		results,
	}
}

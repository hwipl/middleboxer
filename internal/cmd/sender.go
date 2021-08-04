package cmd

import (
	"log"
)

// sender is a test in sender mode
type sender struct {
	test     *MessageTest
	results  chan *MessageResult
	listener *packetListener
}

// sendPacket sends packet
func (s *sender) sendPacket(packet []byte) {
	if err := s.listener.send(packet); err != nil {
		s.results <- &MessageResult{s.test.ID, ResultError}
		log.Println(err)
	}
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
		packetListeners.get(test.Device),
	}
}

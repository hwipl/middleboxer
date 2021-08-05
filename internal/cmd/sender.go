package cmd

import (
	"log"

	"github.com/google/gopacket"
)

// senderPacket is a test packet sent by the sender
type senderPacket struct {
	test   *MessageTest
	layers []gopacket.SerializableLayer
	b      []byte
}

// bytes returns the packet as bytes
func (s *senderPacket) bytes() []byte {
	return s.b
}

// newSenderPacket creates a new packet to send
func newSenderPacket(test *MessageTest) *senderPacket {
	s := senderPacket{
		test,
		[]gopacket.SerializableLayer{},
		[]byte{},
	}
	return &s
}

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

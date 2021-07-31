package cmd

import (
	"github.com/google/gopacket"
	"github.com/hwipl/packet-go/pkg/pcap"
)

// receiver is a test in receiver mode
type receiver struct {
	test    *MessageTest
	results chan *MessageResult
}

// HandlePacket handles a packet received via pcap
func (r *receiver) HandlePacket(packet gopacket.Packet) {
	// TODO: do something
}

// run runs the receiver
func (r *receiver) run() {
	// create listener
	listener := pcap.Listener{
		PacketHandler: r,
		Device:        r.test.Device,
	}

	// prepare listener
	listener.Prepare()

	// TODO: start loop, inform server
}

// newReceiver creates a new test in receiver mode
func newReceiver(test *MessageTest, results chan *MessageResult) *receiver {
	return &receiver{
		test,
		results,
	}
}

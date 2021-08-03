package cmd

import (
	"github.com/google/gopacket"
	"github.com/hwipl/packet-go/pkg/pcap"
)

// packetListener is a pcap packet listener
type packetListener struct {
	listener pcap.Listener
	handlers []pcap.PacketHandler
	packets  chan gopacket.Packet
}

// HandlePacket implements the PacketHandler interface and moves all packets
// into the packets channel
func (p *packetListener) HandlePacket(packet gopacket.Packet) {
	p.packets <- packet
}

// loop is the main loop of this packet listener
func (p *packetListener) loop() {
	for {
		select {
		case packet, more := <-p.packets:
			// handle packets coming from pcap
			if !more {
				return
			}
			for _, h := range p.handlers {
				h.HandlePacket(packet)
			}
		}
	}
}

// newPacketListener creates a new packet listener
func newPacketListener(device string) *packetListener {
	// create packet listener
	p := &packetListener{
		packets: make(chan gopacket.Packet),
	}

	// create pcap listener
	p.listener = pcap.Listener{
		PacketHandler: p,
		Device:        device,
	}

	// prepare pcap listener
	p.listener.Prepare()

	// start pcap listener loop
	go p.listener.Loop()

	// start packet listener loop
	go p.loop()

	return p
}

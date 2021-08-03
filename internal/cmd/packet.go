package cmd

import (
	"github.com/google/gopacket"
	"github.com/hwipl/packet-go/pkg/pcap"
)

// packetListenerReg is a (de)register message for the packet listener
type packetListenerReg struct {
	add     bool
	handler pcap.PacketHandler
}

// packetListener is a pcap packet listener
type packetListener struct {
	listener pcap.Listener
	handlers []pcap.PacketHandler
	regs     chan packetListenerReg
	packets  chan gopacket.Packet
}

// HandlePacket implements the PacketHandler interface and moves all packets
// into the packets channel
func (p *packetListener) HandlePacket(packet gopacket.Packet) {
	p.packets <- packet
}

// removeHandler removes the packet handler from the packet listener
func (p *packetListener) removeHandler(handler pcap.PacketHandler) {
	for i, h := range p.handlers {
		if h == handler {
			p.handlers = append(p.handlers[:i],
				p.handlers[i+1:]...)
			return
		}
	}
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

		case reg, more := <-p.regs:
			// handle packet handler (de)registrations
			if !more {
				return
			}
			if reg.add {
				// TODO: add handler
			} else {
				p.removeHandler(reg.handler)
			}
		}
	}
}

// newPacketListener creates a new packet listener
func newPacketListener(device string) *packetListener {
	// create packet listener
	p := &packetListener{
		regs:    make(chan packetListenerReg),
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

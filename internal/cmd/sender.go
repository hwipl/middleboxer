package cmd

import (
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
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

// createPacketEthernet creates the ethernet header of the packet
func (s *senderPacket) createPacketEthernet() {
	eth := layers.Ethernet{
		SrcMAC: s.test.SrcMAC,
		DstMAC: s.test.DstMAC,
	}

	if s.test.SrcIP.To4() != nil {
		eth.EthernetType = layers.EthernetTypeIPv4
	} else {
		eth.EthernetType = layers.EthernetTypeIPv6
	}

	s.layers = append(s.layers, &eth)
}

// createPacketIPv4 creates the ipv4 header of the packet
func (s *senderPacket) createPacketIPv4() {
	ip := layers.IPv4{
		Version: 4,
		Flags:   layers.IPv4DontFragment,
		TTL:     64,
		SrcIP:   s.test.SrcIP,
		DstIP:   s.test.DstIP,
	}

	switch s.test.Protocol {
	case ProtocolUDP:
		ip.Protocol = layers.IPProtocolUDP
	case ProtocolTCP:
		ip.Protocol = layers.IPProtocolTCP
	}

	s.layers = append(s.layers, &ip)
}

// createPacketIP creates the ip header of the packet
func (s *senderPacket) createPacketIP() {
	if s.test.SrcIP.To4() != nil {
		s.createPacketIPv4()
	} else {
		// TODO: add ipv6
	}
}

// createPacketUDP creates the udp header of the packet
func (s *senderPacket) createPacketUDP() {
	udp := layers.UDP{
		SrcPort: layers.UDPPort(s.test.SrcPort),
		DstPort: layers.UDPPort(s.test.DstPort),
	}
	layer3 := s.layers[1].(gopacket.NetworkLayer)
	udp.SetNetworkLayerForChecksum(layer3)

	s.layers = append(s.layers, &udp)
}

// createPacketL4 creates the layer 4 header of the packet
func (s *senderPacket) createPacketL4() {
	switch s.test.Protocol {
	case ProtocolUDP:
		s.createPacketUDP()
	case ProtocolTCP:
		// TODO: add tcp
	}
}

// createPacket creates the packet to send
func (s *senderPacket) createPacket() {
	// create packet layers
	s.createPacketEthernet()
	s.createPacketIP()
	s.createPacketL4()

	// serialize packet to bytes
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	buf := gopacket.NewSerializeBuffer()
	err := gopacket.SerializeLayers(buf, opts, s.layers...)
	if err != nil {
		log.Fatal(err)
	}
	s.b = buf.Bytes()
}

// newSenderPacket creates a new packet to send
func newSenderPacket(test *MessageTest) *senderPacket {
	s := senderPacket{
		test,
		[]gopacket.SerializableLayer{},
		[]byte{},
	}
	s.createPacket()
	return &s
}

// sender is a test in sender mode
type sender struct {
	test     *MessageTest
	results  chan *MessageResult
	listener *packetListener
	packet   []byte
}

// sendPacket sends packet
func (s *sender) sendPacket() {
	if err := s.listener.send(s.packet); err != nil {
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
		newSenderPacket(test).bytes(),
	}
}

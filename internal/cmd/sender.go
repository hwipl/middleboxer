package cmd

import (
	"log"
	"time"

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

// createPacketIPv6 creates the ipv header of the packet
func (s *senderPacket) createPacketIPv6() {
	ip := layers.IPv6{
		Version:  6,
		HopLimit: 64,
		SrcIP:    s.test.SrcIP,
		DstIP:    s.test.DstIP,
	}

	switch s.test.Protocol {
	case ProtocolUDP:
		ip.NextHeader = layers.IPProtocolUDP
	case ProtocolTCP:
		ip.NextHeader = layers.IPProtocolTCP
	}

	s.layers = append(s.layers, &ip)
}

// createPacketIP creates the ip header of the packet
func (s *senderPacket) createPacketIP() {
	if s.test.SrcIP.To4() != nil {
		s.createPacketIPv4()
	} else {
		s.createPacketIPv6()
	}
}

// createPacketTCP creates the tcp header of the packet
func (s *senderPacket) createPacketTCP() {
	tcp := layers.TCP{
		SrcPort: layers.TCPPort(s.test.SrcPort),
		DstPort: layers.TCPPort(s.test.DstPort),
		SYN:     true,
		Window:  64000,
	}
	layer3 := s.layers[1].(gopacket.NetworkLayer)
	if err := tcp.SetNetworkLayerForChecksum(layer3); err != nil {
		log.Fatal(err)
	}

	s.layers = append(s.layers, &tcp)
}

// createPacketUDP creates the udp header of the packet
func (s *senderPacket) createPacketUDP() {
	udp := layers.UDP{
		SrcPort: layers.UDPPort(s.test.SrcPort),
		DstPort: layers.UDPPort(s.test.DstPort),
	}
	layer3 := s.layers[1].(gopacket.NetworkLayer)
	if err := udp.SetNetworkLayerForChecksum(layer3); err != nil {
		log.Fatal(err)
	}

	s.layers = append(s.layers, &udp)
}

// createPacketL4 creates the layer 4 header of the packet
func (s *senderPacket) createPacketL4() {
	switch s.test.Protocol {
	case ProtocolUDP:
		s.createPacketUDP()
	case ProtocolTCP:
		s.createPacketTCP()
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

// handleIPv4 checks if ip addresses match
func (s *sender) handleIPv4(packet gopacket.Packet) bool {
	// get ip layer
	ipv4Layer := packet.Layer(layers.LayerTypeIPv4)
	if ipv4Layer == nil {
		return false
	}
	ipv4, _ := ipv4Layer.(*layers.IPv4)

	// check ips
	if !ipv4.SrcIP.Equal(s.test.DstIP) || !ipv4.DstIP.Equal(s.test.SrcIP) {
		return false
	}
	return true
}

// handleIPv6 checks if ip addresses match
func (s *sender) handleIPv6(packet gopacket.Packet) bool {
	// get ip layer
	ipv6Layer := packet.Layer(layers.LayerTypeIPv6)
	if ipv6Layer == nil {
		return false
	}
	ipv6, _ := ipv6Layer.(*layers.IPv6)

	// check ips
	if !ipv6.SrcIP.Equal(s.test.DstIP) || !ipv6.DstIP.Equal(s.test.SrcIP) {
		return false
	}
	return true
}

// handleIP checks if ip addresses match
func (s *sender) handleIP(packet gopacket.Packet) bool {
	// TODO: also check if only one ip address matches
	if s.test.SrcIP.To4() != nil {
		return s.handleIPv4(packet)
	}
	return s.handleIPv6(packet)
}

// handleICMPv4 handles ICMPv4 destination unreachable messages
func (s *sender) handleICMPv4(packet gopacket.Packet) {
	// handle icmp messages only
	icmpv4Layer := packet.Layer(layers.LayerTypeICMPv4)
	if icmpv4Layer == nil {
		return
	}
	icmpv4, _ := icmpv4Layer.(*layers.ICMPv4)

	// handle destination unreachable messages only
	if icmpv4.TypeCode.Type() != layers.ICMPv4TypeDestinationUnreachable {
		return
	}

	// get encapsulated packet headers
	encap := gopacket.NewPacket(icmpv4.Payload, layers.LayerTypeIPv4,
		gopacket.Default)

	// get encapsulated ip header
	ipv4Layer := encap.Layer(layers.LayerTypeIPv4)
	if ipv4Layer == nil {
		return
	}
	ipv4, _ := ipv4Layer.(*layers.IPv4)

	// check ip addresses
	if !ipv4.SrcIP.Equal(s.test.SrcIP) || !ipv4.DstIP.Equal(s.test.DstIP) {
		return
	}

	// get encapsulated udp header
	udpLayer := encap.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		return
	}
	udp, _ := udpLayer.(*layers.UDP)

	// check ports
	if udp.SrcPort != layers.UDPPort(s.test.SrcPort) ||
		udp.DstPort != layers.UDPPort(s.test.DstPort) {
		return
	}

	// create result based on icmp code
	result := &MessageResult{
		ID: s.test.ID,
	}
	switch code := icmpv4.TypeCode.Code(); code {
	case layers.ICMPv4CodePort:
		result.Result = ResultICMPv4PortUnreachable
	default:
		log.Println("unexpected icmpv4 type code:", code)
	}

	// send result back to server
	s.results <- result
}

// handleICMPv6 handles ICMPv6 destination unreachable messages
func (s *sender) handleICMPv6(packet gopacket.Packet) {
	// handle icmp messages only
	icmpv6Layer := packet.Layer(layers.LayerTypeICMPv6)
	if icmpv6Layer == nil {
		return
	}
	icmpv6, _ := icmpv6Layer.(*layers.ICMPv6)

	// handle destination unreachable messages only
	if icmpv6.TypeCode.Type() != layers.ICMPv6TypeDestinationUnreachable {
		return
	}

	// get encapsulated packet headers, skipping first 4 bytes (unused)
	encap := gopacket.NewPacket(icmpv6.Payload[4:], layers.LayerTypeIPv6,
		gopacket.Default)

	// get encapsulated ip header
	ipv6Layer := encap.Layer(layers.LayerTypeIPv6)
	if ipv6Layer == nil {
		return
	}
	ipv6, _ := ipv6Layer.(*layers.IPv6)

	// check ip addresses
	if !ipv6.SrcIP.Equal(s.test.SrcIP) || !ipv6.DstIP.Equal(s.test.DstIP) {
		return
	}

	// get encapsulated udp header
	udpLayer := encap.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		return
	}
	udp, _ := udpLayer.(*layers.UDP)

	// check ports
	if udp.SrcPort != layers.UDPPort(s.test.SrcPort) ||
		udp.DstPort != layers.UDPPort(s.test.DstPort) {
		return
	}

	// create result based on icmp code
	result := &MessageResult{
		ID: s.test.ID,
	}
	switch code := icmpv6.TypeCode.Code(); code {
	case layers.ICMPv6CodePortUnreachable:
		result.Result = ResultICMPv6PortUnreachable
	default:
		log.Println("unexpected icmpv6 type code:", code)
	}

	// send result back to server
	s.results <- result
}

// handleTCPReset handles TCP reset messages
func (s *sender) handleTCPReset(packet gopacket.Packet) {
	// handle tcp messages only
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer == nil {
		return
	}
	tcp, _ := tcpLayer.(*layers.TCP)

	// handle reset messages only
	if !tcp.RST {
		return
	}

	// check ports
	if tcp.SrcPort != layers.TCPPort(s.test.DstPort) ||
		tcp.DstPort != layers.TCPPort(s.test.SrcPort) {
		return
	}

	// send result back to server
	s.results <- &MessageResult{
		s.test.ID,
		ResultTCPReset,
	}
}

// HandlePacket handles a packet received via the listener
func (s *sender) HandlePacket(packet gopacket.Packet) {
	if !s.handleIP(packet) {
		return
	}
	s.handleICMPv4(packet)
	s.handleICMPv6(packet)
	s.handleTCPReset(packet)
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
	// register handler to receive icmp error messages
	s.listener.register(s)

	// send test packet three times
	for i := 0; i < 3; i++ {
		s.sendPacket()
		time.Sleep(time.Millisecond)
	}

	// wait a second for icmp errors and stop
	time.Sleep(time.Second)
	s.listener.deregister(s)
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

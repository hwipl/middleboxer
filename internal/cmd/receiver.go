package cmd

import (
	"bytes"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// receiver is a test in receiver mode
type receiver struct {
	test    *MessageTest
	results chan *MessageResult
}

// handleEthernet checks if ethernet values in packet match the current test
func (r *receiver) handleEthernet(packet gopacket.Packet) bool {
	// if we do not care about mac addresses, skip the following checks
	if r.test.SrcMAC == nil && r.test.DstMAC == nil {
		return true
	}

	// get ethernet header
	ethLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethLayer == nil {
		return false
	}
	eth, _ := ethLayer.(*layers.Ethernet)

	// check mac addresses
	if r.test.SrcMAC != nil && !bytes.Equal(eth.SrcMAC, r.test.SrcMAC) {
		return false
	}
	if r.test.DstMAC != nil && !bytes.Equal(eth.DstMAC, r.test.DstMAC) {
		return false
	}

	return true
}

// checkIPs checks if src and dst ip addresses match current test
func (r *receiver) checkIPs(src, dst net.IP) bool {
	if r.test.SrcIP != nil && !r.test.SrcIP.Equal(src) {
		return false
	}
	if r.test.DstIP != nil && !r.test.DstIP.Equal(dst) {
		return false
	}

	return true
}

// handleIPv4 checks if IPv4 values in packet match the current test
func (r *receiver) handleIPv4(packet gopacket.Packet) bool {
	// get ip header
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		return false
	}
	ip, _ := ipLayer.(*layers.IPv4)

	// check ip addresses
	return r.checkIPs(ip.SrcIP, ip.DstIP)
}

// handleIPv6 checks if IPv6 values in packet match the current test
func (r *receiver) handleIPv6(packet gopacket.Packet) bool {
	// get ip header
	ipLayer := packet.Layer(layers.LayerTypeIPv6)
	if ipLayer == nil {
		return false
	}
	ip, _ := ipLayer.(*layers.IPv6)

	// check ip addresses
	return r.checkIPs(ip.SrcIP, ip.DstIP)
}

// handleIP checks if IPv4 or IPv6 values in packet match the current test
func (r *receiver) handleIP(packet gopacket.Packet) bool {
	// if we do not care about ip addresses, skip the following checks
	if r.test.SrcIP == nil && r.test.DstIP == nil {
		return true
	}

	// check ipv4 or ipv6 addresses
	if r.test.SrcIP.To4() != nil {
		return r.handleIPv4(packet)
	}
	return r.handleIPv6(packet)
}

// checkPorts checks if src and dst ports match current test
func (r *receiver) checkPorts(src, dst uint16) bool {
	if r.test.SrcPort != 0 && r.test.SrcPort != src {
		return false
	}
	if r.test.DstPort != 0 && r.test.DstPort != dst {
		return false
	}

	return true
}

// handleTCP checks if tcp values in packet match the current test
func (r *receiver) handleTCP(packet gopacket.Packet) bool {
	// if we do not care about ports, skip the following checks
	if r.test.SrcPort == 0 && r.test.DstPort == 0 {
		return true
	}

	// get tcp header
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer == nil {
		return false
	}
	tcp, _ := tcpLayer.(*layers.TCP)

	// check ports
	return r.checkPorts(uint16(tcp.SrcPort), uint16(tcp.DstPort))
}

// handleUDP checks if udp values in packet match the current test
func (r *receiver) handleUDP(packet gopacket.Packet) bool {
	// if we do not care about ports, skip the following checks
	if r.test.SrcPort == 0 && r.test.DstPort == 0 {
		return true
	}

	// get udp header
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		return false
	}
	udp, _ := udpLayer.(*layers.UDP)

	// check ports
	return r.checkPorts(uint16(udp.SrcPort), uint16(udp.DstPort))
}

// handleL4 checks if L4 values in packet match the current test
func (r *receiver) handleL4(packet gopacket.Packet) bool {
	// if we do not care about l4, skip the following checks
	if r.test.Protocol == ProtocolNone {
		return true
	}

	// check tcp
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		if r.test.Protocol != ProtocolTCP {
			return false
		}
		return r.handleTCP(packet)
	}

	// check udp
	if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		if r.test.Protocol != ProtocolUDP {
			return false
		}
		return r.handleUDP(packet)
	}

	return false
}

// HandlePacket handles a packet received via pcap
func (r *receiver) HandlePacket(packet gopacket.Packet) {
	// check ethernet
	if !r.handleEthernet(packet) {
		return
	}

	// check ip
	if !r.handleIP(packet) {
		return
	}

	// check layer 4
	if !r.handleL4(packet) {
		return
	}

	// send result back to server
	r.results <- &MessageResult{
		ID:     r.test.ID,
		Result: ResultPass,
		Packet: packet.Data(),
	}
}

// run runs the receiver
func (r *receiver) run() {
	// register receiver as packet handler on device and
	// tell server we are are ready
	packetListeners.get(r.test.Device).register(r)
	r.results <- &MessageResult{
		ID:     r.test.ID,
		Result: ResultReady,
	}

	// wait two seconds and stop in case we do not get the packet
	// from the sender
	time.Sleep(2 * time.Second)
	packetListeners.get(r.test.Device).deregister(r)
}

// newReceiver creates a new test in receiver mode
func newReceiver(test *MessageTest, results chan *MessageResult) *receiver {
	return &receiver{
		test,
		results,
	}
}

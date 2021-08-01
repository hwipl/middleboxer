package cmd

import (
	"bytes"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/hwipl/packet-go/pkg/pcap"
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
	if r.test.DstIP != nil && !r.test.SrcIP.Equal(dst) {
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

	// TODO: l4 protocol, ports, send result back, stop listener
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

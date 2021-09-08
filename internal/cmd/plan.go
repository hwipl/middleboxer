package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// plan result constants
const (
	planResultPass = iota
	planResultReject
	planResultDrop
)

// planResultRange is a port range with the same plan result
type planResultRange struct {
	result    uint8
	firstPort uint16
	lastPort  uint16
}

// planResults is a collection of results of a completed plan for printing
type planResults struct {
	ranges []*planResultRange
}

// String converts planResults to a string
func (p *planResults) String() string {
	s := ""
	for _, r := range p.ranges {
		if r.firstPort == r.lastPort {
			s += fmt.Sprintf("%d\t", r.firstPort)
		} else {
			s += fmt.Sprintf("%d:%d\t", r.firstPort, r.lastPort)
		}
		switch r.result {
		case planResultPass:
			s += fmt.Sprintf("pass\n")
		case planResultReject:
			s += fmt.Sprintf("reject\n")
		case planResultDrop:
			s += fmt.Sprintf("drop\n")
		}
	}
	return s
}

// add adds r to the collection of results; expects results added with
// increasing port numbers, without gaps
func (p *planResults) add(port uint16, result uint8) {
	if length := len(p.ranges); length > 0 &&
		p.ranges[length-1].result == result &&
		p.ranges[length-1].lastPort == port-1 {
		p.ranges[length-1].lastPort = port
	} else {
		newRange := &planResultRange{
			result:    result,
			firstPort: port,
			lastPort:  port,
		}
		p.ranges = append(p.ranges, newRange)
	}
}

// planItem is a specific test in a test execution plan
type planItem struct {
	ID              uint32
	Port            uint16
	SenderMsg       *MessageTest
	ReceiverMsg     *MessageTest
	receiverReady   bool
	SenderResults   []*MessageResult
	ReceiverResults []*MessageResult
}

// containsPass checks if plan item contains a passing result
func (p *planItem) containsPass() bool {
	for _, r := range p.ReceiverResults {
		if r.Result == ResultPass {
			return true
		}
		log.Println("other result:", r)
	}
	return false
}

// containsReject checks if plan item contains a rejected result
func (p *planItem) containsReject() bool {
	for _, r := range p.SenderResults {
		switch r.Result {
		case ResultICMPv4NetworkUnreachable,
			ResultICMPv4HostUnreachable,
			ResultICMPv4ProtocolUnreachable,
			ResultICMPv4PortUnreachable,
			ResultICMPv4FragmentationNeeded,
			ResultICMPv4SourceRoutingFailed,
			ResultICMPv4NetworkUnknown,
			ResultICMPv4HostUnknown,
			ResultICMPv4SourceIsolated,
			ResultICMPv4NetworkProhibited,
			ResultICMPv4HostProhibited,
			ResultICMPv4NetworkTOS,
			ResultICMPv4HostTOS,
			ResultICMPv4CommProhibited,
			ResultICMPv4HostPrecedence,
			ResultICMPv4PrecedenceCutoff,
			ResultICMPv6NoRouteToDst,
			ResultICMPv6AdminProhibited,
			ResultICMPv6BeyondScopeOfSrc,
			ResultICMPv6AddressUnreachable,
			ResultICMPv6PortUnreachable,
			ResultICMPv6SrcAddressFailed,
			ResultICMPv6RejectRouteToDst,
			ResultICMPv6SrcRoutingHeader,
			ResultICMPv6HeadersTooLong,
			ResultTCPReset:
			return true
		default:
			log.Println("other result:", r)
		}
	}
	return false
}

// containsDrop checks if plan item contains a dropped result
func (p *planItem) containsDrop() bool {
	if len(p.ReceiverResults) == 0 && len(p.SenderResults) == 0 {
		return true
	}
	return false
}

// getEthernetDiffs returns differences in ethernet fields as string
func (p *planItem) getEthernetDiffs(packet gopacket.Packet) string {
	// get ethernet header
	ethLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethLayer == nil {
		return ""
	}
	eth, _ := ethLayer.(*layers.Ethernet)

	// check mac addresses
	s := ""
	if p.SenderMsg.SrcMAC != nil && !bytes.Equal(eth.SrcMAC, p.SenderMsg.SrcMAC) {
		s += fmt.Sprintf("src mac: %s -> %s\n", p.SenderMsg.SrcMAC,
			eth.SrcMAC)
	}
	if p.SenderMsg.DstMAC != nil && !bytes.Equal(eth.DstMAC, p.SenderMsg.DstMAC) {
		s += fmt.Sprintf("dst mac: %s -> %s\n", p.SenderMsg.DstMAC,
			eth.DstMAC)
	}
	return s
}

// getIPAddrDiffs returns differences in ip addresses as string
func (p *planItem) getIPAddrDiffs(src, dst net.IP) string {
	s := ""
	if p.SenderMsg.SrcIP != nil && !p.SenderMsg.SrcIP.Equal(src) {
		s += fmt.Sprintf("src ip: %s -> %s\n", p.SenderMsg.SrcIP, src)
	}
	if p.SenderMsg.DstIP != nil && !p.SenderMsg.DstIP.Equal(dst) {
		s += fmt.Sprintf("dst ip: %s -> %s\n", p.SenderMsg.DstIP, dst)
	}
	return s
}

// getIPv4Diffs returns differences in ipv4 fields as string
func (p *planItem) getIPv4Diffs(packet gopacket.Packet) string {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		return ""
	}
	ip, _ := ipLayer.(*layers.IPv4)
	return p.getIPAddrDiffs(ip.SrcIP, ip.DstIP)
}

// getIPv6fDiffs returns differences in ipv6 fields as string
func (p *planItem) getIPv6Diffs(packet gopacket.Packet) string {
	ipLayer := packet.Layer(layers.LayerTypeIPv6)
	if ipLayer == nil {
		return ""
	}
	ip, _ := ipLayer.(*layers.IPv6)
	return p.getIPAddrDiffs(ip.SrcIP, ip.DstIP)
}

// getIPDiffs returns differences in ip fields as string
func (p *planItem) getIPDiffs(packet gopacket.Packet) string {
	ip4Layer := packet.Layer(layers.LayerTypeIPv4)
	if ip4Layer != nil {
		return p.getIPv4Diffs(packet)
	}

	ip6Layer := packet.Layer(layers.LayerTypeIPv6)
	if ip6Layer != nil {
		return p.getIPv6Diffs(packet)
	}

	log.Println("packet does not contain ip header")
	return ""
}

// getPortDiffs returns differences in port numbers as string
func (p *planItem) getPortDiffs(src, dst uint16) string {
	s := ""
	if p.SenderMsg.SrcPort != src {
		s += fmt.Sprintf("src port: %d -> %d\n", p.SenderMsg.SrcPort,
			src)
	}
	if p.SenderMsg.DstPort != dst {
		s += fmt.Sprintf("dst port: %d -> %d\n", p.SenderMsg.DstPort,
			dst)
	}
	return s
}

// getTCPDiffs returns differences in tcp fields as string
func (p *planItem) getTCPDiffs(packet gopacket.Packet) string {
	// get tcp header
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer == nil {
		return ""
	}
	tcp, _ := tcpLayer.(*layers.TCP)

	// check ports
	return p.getPortDiffs(uint16(tcp.SrcPort), uint16(tcp.DstPort))
}

// getUDPDiffs returns differences in udp fields as string
func (p *planItem) getUDPDiffs(packet gopacket.Packet) string {
	// get udp header
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		return ""
	}
	udp, _ := udpLayer.(*layers.UDP)

	// check ports
	return p.getPortDiffs(uint16(udp.SrcPort), uint16(udp.DstPort))
}

// getL4Diffs returns differences in l4 fields as string
func (p *planItem) getL4Diffs(packet gopacket.Packet) string {
	// check tcp
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		return p.getTCPDiffs(packet)
	}

	// check udp
	if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		return p.getUDPDiffs(packet)
	}

	log.Println("packet does not contain expected l4 header")
	return ""
}

// printPacketDiffs prints differences in packet fields
func (p *planItem) printPacketDiffs() {
	last := ""
	for _, r := range p.ReceiverResults {
		if r.Result != ResultPass {
			continue
		}
		packet := gopacket.NewPacket(r.Packet,
			layers.LayerTypeEthernet, gopacket.Default)
		s := ""
		s += p.getEthernetDiffs(packet)
		s += p.getIPDiffs(packet)
		s += p.getL4Diffs(packet)
		if s != "" && s != last {
			log.Printf("Port %d diffs:\n%s", p.Port, s)
			last = s
		}
	}
}

// newPlanItem creates a new planItem
func newPlanItem(id uint32, port uint16, senderMsg, receiverMsg *MessageTest) *planItem {
	return &planItem{
		ID:          id,
		Port:        port,
		SenderMsg:   senderMsg,
		ReceiverMsg: receiverMsg,
	}
}

// plan is a test execution plan
type plan struct {
	senderID       uint8
	receiverID     uint8
	items          map[uint32]*planItem
	currentItem    uint32
	senderActive   bool
	receiverActive bool
}

// isSender checks if clientID is in the senders list
func (p *plan) isSender(clientID uint8) bool {
	if p.senderID == clientID {
		return true
	}
	return false
}

// isReceiver checks if clientID is in the receivers list
func (p *plan) isReceiver(clientID uint8) bool {
	if p.receiverID == clientID {
		return true
	}
	return false
}

// handleResult handles result coming from clientID
func (p *plan) handleResult(clientID uint8, result *MessageResult) {
	// check if client is a sender or receiver
	isSender := p.isSender(clientID)
	if !isSender {
		if !p.isReceiver(clientID) {
			log.Println("Received result from invalid client")
			return
		}
	}

	// get plan item
	item := p.items[result.ID]
	if item == nil {
		log.Println("Received result with invalid ID")
		return
	}

	// add result to result list
	if isSender {
		item.SenderResults = append(item.SenderResults, result)
	} else {
		if result.Result == ResultReady {
			// handle "ready" results
			if item.receiverReady {
				log.Println("Double ready from client", clientID)
				return
			}
			item.receiverReady = true
			return
		}

		// handle other results
		item.ReceiverResults = append(item.ReceiverResults, result)
	}
}

// handleClient handles a new client
func (p *plan) handleClient(clientID uint8) {
	// is client the sender?
	if clientID == p.senderID {
		if p.senderActive {
			log.Println("Sender client already active")
		}
		p.senderActive = true
		return
	}

	// is client the receiver?
	if clientID == p.receiverID {
		if p.receiverActive {
			log.Println("Receiver client already active")
		}
		p.receiverActive = true
		return
	}

	// invalid client
	log.Println("Invalid client")
}

// clientsActive checks if all clients are active
func (p *plan) clientsActive() bool {
	if p.senderActive && p.receiverActive {
		return true
	}
	return false
}

// getCurrentItem returns the current plan item
func (p *plan) getCurrentItem() *planItem {
	return p.items[p.currentItem]
}

// getNextItem returns the next plan item
func (p *plan) getNextItem() *planItem {
	p.currentItem++
	return p.items[p.currentItem]
}

// printResults prints results of this plan to the console
func (p *plan) printResults() {
	i := uint32(0)
	results := planResults{}
	for {
		item := p.items[i]
		if item == nil {
			break
		}

		switch {
		case item.containsPass():
			results.add(item.Port, planResultPass)
		case item.containsReject():
			results.add(item.Port, planResultReject)
		case item.containsDrop():
			results.add(item.Port, planResultDrop)
		}

		i++
	}
	log.Printf("Printing results:\n%s", &results)
}

func (p *plan) printPacketDiffs() {
	i := uint32(0)
	for {
		item := p.items[i]
		if item == nil {
			break
		}
		item.printPacketDiffs()
		i++
	}
}

// writeFile writes all plan items including results to file
func (p *plan) writeFile(file string) {
	log.Println("Writing plan to file", file)
	j, err := json.MarshalIndent(p.items, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(file, j, 0600); err != nil {
		log.Fatal(err)
	}
}

// newSenderMessage creates a new sender message for a plan
func newSenderMessage(id uint32, port uint16, config *Config) *MessageTest {
	return &MessageTest{
		ID:       id,
		Initiate: true,
		Device:   config.SenderDevice,
		SrcMAC:   config.GetSenderSrcMAC(),
		DstMAC:   config.GetSenderDstMAC(),
		SrcIP:    config.GetSenderSrcIP(),
		DstIP:    config.GetSenderDstIP(),
		Protocol: config.Protocol,
		SrcPort:  config.SenderSrcPort,
		DstPort:  port,
	}
}

// newReceiverMessage creates a new receiver message for a plan
func newReceiverMessage(id uint32, port uint16, config *Config) *MessageTest {
	return &MessageTest{
		ID:       id,
		Initiate: false,
		Device:   config.ReceiverDevice,
		SrcMAC:   config.GetReceiverSrcMAC(),
		DstMAC:   config.GetReceiverDstMAC(),
		SrcIP:    config.GetReceiverSrcIP(),
		DstIP:    config.GetReceiverDstIP(),
		Protocol: config.Protocol,
		SrcPort:  config.SenderSrcPort,
		DstPort:  port,
	}
}

// newPlan creates a new plan
func newPlan(config *Config) *plan {
	// initialize plan
	items := make(map[uint32]*planItem)

	// fill plan with plan items
	id := uint32(0)
	first, last := config.GetPortRange()
	for i := first; i <= last && i != 0; i++ {
		senderMsg := newSenderMessage(id, i, config)
		receiverMsg := newReceiverMessage(id, i, config)
		item := newPlanItem(id, i, senderMsg, receiverMsg)
		items[id] = item
		id++
	}

	return &plan{
		senderID:   config.SenderID,
		receiverID: config.ReceiverID,
		items:      items,
	}
}

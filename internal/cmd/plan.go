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

// planPacketDiff is a difference in packet fields
type planPacketDiff struct {
	Field    string
	Sender   string
	Receiver string
}

// String converts a packet difference to a string
func (p *planPacketDiff) String() string {
	return fmt.Sprintf("%s: %s -> %s", p.Field, p.Sender, p.Receiver)
}

type planPacketDiffs []*planPacketDiff

func (p *planPacketDiffs) contains(diff *planPacketDiff) bool {
	for _, d := range *p {
		if *d == *diff {
			return true
		}
	}
	return false
}

func (p *planPacketDiffs) add(field, sender, receiver string) {
	d := &planPacketDiff{field, sender, receiver}
	if p.contains(d) {
		return
	}
	*p = append(*p, d)
}

func (p *planPacketDiffs) String() string {
	s := ""
	for i, d := range *p {
		if i == 0 {
			s += fmt.Sprintf("%s", d)
		} else {
			s += fmt.Sprintf("\n%s", d)
		}
	}
	return s
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
	PacketDiffs     planPacketDiffs
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

// getEthernetDiffs gets differences in ethernet fields
func (p *planItem) getEthernetDiffs(packet gopacket.Packet) {
	// get ethernet header
	ethLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethLayer == nil {
		return
	}
	eth, _ := ethLayer.(*layers.Ethernet)

	// check mac addresses
	if p.SenderMsg.SrcMAC != nil && !bytes.Equal(eth.SrcMAC, p.SenderMsg.SrcMAC) {
		p.PacketDiffs.add(
			"SrcMAC",
			fmt.Sprintf("%s", p.SenderMsg.SrcMAC),
			fmt.Sprintf("%s", eth.SrcMAC),
		)
	}
	if p.SenderMsg.DstMAC != nil && !bytes.Equal(eth.DstMAC, p.SenderMsg.DstMAC) {
		p.PacketDiffs.add(
			"DstMAC",
			fmt.Sprintf("%s", p.SenderMsg.DstMAC),
			fmt.Sprintf("%s", eth.DstMAC),
		)
	}
}

// getIPAddrDiffs gets differences in ip addresses
func (p *planItem) getIPAddrDiffs(src, dst net.IP) {
	if p.SenderMsg.SrcIP != nil && !p.SenderMsg.SrcIP.Equal(src) {
		p.PacketDiffs.add(
			"SrcIP",
			fmt.Sprintf("%s", p.SenderMsg.SrcIP),
			fmt.Sprintf("%s", src),
		)
	}
	if p.SenderMsg.DstIP != nil && !p.SenderMsg.DstIP.Equal(dst) {
		p.PacketDiffs.add(
			"DstIP",
			fmt.Sprintf("%s", p.SenderMsg.DstIP),
			fmt.Sprintf("%s", dst),
		)
	}
}

// getIPv4Diffs gets differences in ipv4 fields
func (p *planItem) getIPv4Diffs(packet gopacket.Packet) {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		return
	}
	ip, _ := ipLayer.(*layers.IPv4)
	p.getIPAddrDiffs(ip.SrcIP, ip.DstIP)
}

// getIPv6fDiffs getss differences in ipv6 fields
func (p *planItem) getIPv6Diffs(packet gopacket.Packet) {
	ipLayer := packet.Layer(layers.LayerTypeIPv6)
	if ipLayer == nil {
		return
	}
	ip, _ := ipLayer.(*layers.IPv6)
	p.getIPAddrDiffs(ip.SrcIP, ip.DstIP)
}

// getIPDiffs gets differences in ip fields
func (p *planItem) getIPDiffs(packet gopacket.Packet) {
	ip4Layer := packet.Layer(layers.LayerTypeIPv4)
	if ip4Layer != nil {
		p.getIPv4Diffs(packet)
		return
	}

	ip6Layer := packet.Layer(layers.LayerTypeIPv6)
	if ip6Layer != nil {
		p.getIPv6Diffs(packet)
		return
	}

	log.Println("packet does not contain ip header")
}

// getPortDiffs gets differences in port numbers
func (p *planItem) getPortDiffs(src, dst uint16) {
	if p.SenderMsg.SrcPort != src {
		p.PacketDiffs.add(
			"SrcPort",
			fmt.Sprintf("%d", p.SenderMsg.SrcPort),
			fmt.Sprintf("%d", src),
		)
	}
	if p.SenderMsg.DstPort != dst {
		p.PacketDiffs.add(
			"DstPort",
			fmt.Sprintf("%d", p.SenderMsg.DstPort),
			fmt.Sprintf("%d", dst),
		)
	}
}

// getTCPDiffs gets differences in tcp fields
func (p *planItem) getTCPDiffs(packet gopacket.Packet) {
	// get tcp header
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer == nil {
		return
	}
	tcp, _ := tcpLayer.(*layers.TCP)

	// check ports
	p.getPortDiffs(uint16(tcp.SrcPort), uint16(tcp.DstPort))
}

// getUDPDiffs gets differences in udp fields
func (p *planItem) getUDPDiffs(packet gopacket.Packet) {
	// get udp header
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		return
	}
	udp, _ := udpLayer.(*layers.UDP)

	// check ports
	p.getPortDiffs(uint16(udp.SrcPort), uint16(udp.DstPort))
}

// getL4Diffs gets differences in l4 fields
func (p *planItem) getL4Diffs(packet gopacket.Packet) {
	// check tcp
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		p.getTCPDiffs(packet)
		return
	}

	// check udp
	if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		p.getUDPDiffs(packet)
		return
	}

	log.Println("packet does not contain expected l4 header")
}

func (p *planItem) getPacketDiffs(packet []byte) {
	pkt := gopacket.NewPacket(packet, layers.LayerTypeEthernet,
		gopacket.Default)

	p.getEthernetDiffs(pkt)
	p.getIPDiffs(pkt)
	p.getL4Diffs(pkt)
}

// printPacketDiffs prints differences in packet fields
func (p *planItem) printPacketDiffs() {
	if len(p.PacketDiffs) > 0 {
		log.Println(fmt.Sprintf("Port %d packet differences:\n%s",
			p.Port, &p.PacketDiffs))
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

		if result.Result == ResultPass {
			// handle "pass" results
			item.getPacketDiffs(result.Packet)
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

// printPacketDiffs prints packet differences to the console
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

package cmd

import (
	"fmt"
	"log"
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

// planResult is a result of a completed plan for printing
type planResult struct {
	port               uint16
	numPass            int
	numPortUnreachable int
	numReset           int
	numOther           int
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

// addRange adds port and result to the port/result ranges
func (p *planResults) addRange(port uint16, result uint8) {
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

// add adds r to the collection of results; expects results added with
// increasing port numbers, without gaps
func (p *planResults) add(r *planResult) {
	// passing packets
	if r.numPass > 0 {
		p.addRange(r.port, planResultPass)
		return
	}

	// rejected packets
	if r.numPortUnreachable > 0 || r.numReset > 0 {
		p.addRange(r.port, planResultReject)
		return
	}

	// dropped packets
	if r.numPass == 0 &&
		r.numPortUnreachable == 0 &&
		r.numReset == 0 &&
		r.numOther == 0 {
		p.addRange(r.port, planResultDrop)
		return
	}

	// other packets
	log.Println("other result:", r)
}

// planItem is a specific test in a test execution plan
type planItem struct {
	id              uint32
	port            uint16
	senderMsg       *MessageTest
	receiverMsg     *MessageTest
	receiverReady   bool
	senderResults   []*MessageResult
	receiverResults []*MessageResult
}

// newPlanItem creates a new planItem
func newPlanItem(id uint32, port uint16, senderMsg, receiverMsg *MessageTest) *planItem {
	return &planItem{
		id:          id,
		port:        port,
		senderMsg:   senderMsg,
		receiverMsg: receiverMsg,
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
		item.senderResults = append(item.senderResults, result)
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
		item.receiverResults = append(item.receiverResults, result)
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

		result := planResult{port: item.port}
		for _, r := range item.senderResults {
			// handle other results
			// TODO: do this properly, check for reject messages
			// TODO: add all unreachable codes?
			switch r.Result {
			case ResultICMPv4PortUnreachable, ResultICMPv6PortUnreachable:
				result.numPortUnreachable++
			case ResultTCPReset:
				result.numReset++
			default:
				result.numOther++
			}
		}
		for _, r := range item.receiverResults {
			if r.Result == ResultPass {
				result.numPass++
				continue
			}
			result.numOther++
		}
		results.add(&result)

		i++
	}
	log.Printf("Printing results:\n%s", &results)
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

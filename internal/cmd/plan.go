package cmd

import "log"

// planItem is a specific test in a test execution plan
type planItem struct {
	id              uint32
	senderMsg       *MessageTest
	receiverMsg     *MessageTest
	senderResults   []*MessageResult
	receiverResults []*MessageResult
}

// newPlanItem creates a new planItem
func newPlanItem(id uint32, senderMsg, receiverMsg *MessageTest) *planItem {
	return &planItem{
		id:          id,
		senderMsg:   senderMsg,
		receiverMsg: receiverMsg,
	}
}

// plan is a test execution plan
type plan struct {
	senderIDs   []uint8
	receiverIDs []uint8
	items       map[uint32]*planItem
	clients     []uint8
}

// listContainsID checks if list contains id
func listContainsID(list []uint8, id uint8) bool {
	for _, i := range list {
		if i == id {
			return true
		}
	}
	return false
}

// isSender checks if clientID is in the senders list
func (p *plan) isSender(clientID uint8) bool {
	if listContainsID(p.senderIDs, clientID) {
		return true
	}
	return false
}

// isReceiver checks if clientID is in the receivers list
func (p *plan) isReceiver(clientID uint8) bool {
	if listContainsID(p.receiverIDs, clientID) {
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
		item.receiverResults = append(item.receiverResults, result)
	}
}

// handleClient handles a new client
func (p *plan) handleClient(clientID uint8) {
	// check if client is valid
	if !listContainsID(p.senderIDs, clientID) {
		if !listContainsID(p.receiverIDs, clientID) {
			log.Println("Invalid client")
			return
		}
	}

	// add client to active clients
	if listContainsID(p.clients, clientID) {
		log.Println("Client already active")
		return
	}
	p.clients = append(p.clients, clientID)
}

// clientsActive checks if all clients are active
func (p *plan) clientsActive() bool {
	// check if all senders are present
	for _, i := range p.senderIDs {
		if !listContainsID(p.clients, i) {
			return false
		}
	}

	// check if all receivers are present
	for _, i := range p.receiverIDs {
		if !listContainsID(p.clients, i) {
			return false
		}
	}

	// senders and receivers are present
	return true
}

// newPlan creates a new plan
func newPlan(senderIDs, receiverIDs []uint8) *plan {
	items := make(map[uint32]*planItem)
	return &plan{
		senderIDs:   senderIDs,
		receiverIDs: receiverIDs,
		items:       items,
	}
}

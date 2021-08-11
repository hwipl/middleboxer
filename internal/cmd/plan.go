package cmd

import "log"

// planItem is a specific test in a test execution plan
type planItem struct {
	id              uint32
	senderMsg       *MessageTest
	receiverMsg     *MessageTest
	receiversReady  []uint8
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
	senderID   uint8
	receiverID uint8
	items      map[uint32]*planItem
	clients    []uint8
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
			if listContainsID(item.receiversReady, clientID) {
				log.Println("Double ready from client", clientID)
				return
			}
			item.receiversReady = append(item.receiversReady, clientID)
			return
		}

		// handle other results
		item.receiverResults = append(item.receiverResults, result)
	}
}

// handleClient handles a new client
func (p *plan) handleClient(clientID uint8) {
	// check if client is valid
	if p.senderID != clientID {
		if p.receiverID != clientID {
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
	if !listContainsID(p.clients, p.senderID) {
		return false
	}

	// check if all receivers are present
	if !listContainsID(p.clients, p.receiverID) {
		return false
	}

	// senders and receivers are present
	return true
}

// newPlan creates a new plan
func newPlan(senderID, receiverID uint8) *plan {
	items := make(map[uint32]*planItem)
	return &plan{
		senderID:   senderID,
		receiverID: receiverID,
		items:      items,
	}
}

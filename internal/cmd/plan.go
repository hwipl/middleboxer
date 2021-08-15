package cmd

import "log"

// planItem is a specific test in a test execution plan
type planItem struct {
	id              uint32
	senderMsg       *MessageTest
	receiverMsg     *MessageTest
	receiverReady   bool
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
	p.currentItem += 1
	return p.items[p.currentItem]
}

// newPlan creates a new plan
func newPlan(config *Config) *plan {
	items := make(map[uint32]*planItem)
	return &plan{
		senderID:   config.SenderID,
		receiverID: config.ReceiverID,
		items:      items,
	}
}

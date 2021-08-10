package cmd

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

// newPlan creates a new plan
func newPlan(senderIDs, receiverIDs []uint8) *plan {
	items := make(map[uint32]*planItem)
	return &plan{
		senderIDs,
		receiverIDs,
		items,
	}
}

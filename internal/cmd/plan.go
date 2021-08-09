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

// newPlan creates a new plan
func newPlan(senderIDs, receiverIDs []uint8) *plan {
	items := make(map[uint32]*planItem)
	return &plan{
		senderIDs,
		receiverIDs,
		items,
	}
}

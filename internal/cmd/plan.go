package cmd

// planItem is a specific test in a test execution plan
type planItem struct {
	id              uint32
	senderMsg       *MessageTest
	receiverMsg     *MessageTest
	senderResults   []*MessageResult
	receiverResults []*MessageResult
}

// plan is a test execution plan
type plan struct {
	senderIDs   []uint8
	receiverIDs []uint8
	items       map[uint32]*planItem
}

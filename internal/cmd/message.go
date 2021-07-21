package cmd

const (
	// MessageHeaderLength is the length of the Type and Length fields
	// of a message
	MessageHeaderLength = 3

	// MessageMaxLength is the maximum length of a message in bytes
	MessageMaxLength = 4096
)

// Message is a TLV message
type Message struct {
	Type   uint8
	Length uint16
	Data   []byte
}


// newMessage creates a new Message with type and data
func newMessage(typ uint8, data []byte) *Message {
	if len(data) > MessageMaxLength-MessageHeaderLength {
		return nil
	}

	return &Message{
		typ,
		uint16(len(data)),
		data,
	}
}

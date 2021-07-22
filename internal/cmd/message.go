package cmd

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
)

const (
	// MessageHeaderLength is the length of the Type and Length fields
	// of a message
	MessageHeaderLength = 3

	// MessageMaxLength is the maximum length of a message in bytes
	MessageMaxLength = 4096
)

// Message types
const (
	MessageTypeNone = iota
	MessageTypeNop
	MessageTypeRegister
	MessageTypeInvalid
)

// Message is an interface for all messages
type Message interface {
	GetType() uint8
}

// MessageNop is a no operation message
type MessageNop struct{}

// GetType returns the type of the message
func (m *MessageNop) GetType() uint8 {
	return MessageTypeNop
}

// MessageRegister is a register message
type MessageRegister struct {
	Client uint8
}

// GetType returns the type of the message
func (m *MessageRegister) GetType() uint8 {
	return MessageTypeRegister
}

// TLVMessage is a TLV message
type TLVMessage struct {
	Type   uint8
	Length uint16
	Data   []byte
}

// serialize encodes a message as bytes
func (m *TLVMessage) serialize() []byte {
	var buf bytes.Buffer

	err := binary.Write(&buf, binary.BigEndian, m)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

// readBytes reads length bytes from conn
func readBytes(conn net.Conn, length uint16) []byte {
	buf := make([]byte, length)
	count := 0
	for count < MessageHeaderLength {
		n, err := conn.Read(buf[count:])
		if err != nil {
			log.Printf("Connection to %s: %s\n",
				conn.RemoteAddr(), err)
			return nil
		}
		count += n
	}
	return buf
}

// writeBytes writes buf to conn
func writeBytes(conn net.Conn, buf []byte) bool {
	count := 0
	for count < len(buf) {
		n, err := conn.Write(buf[count:])
		if err != nil {
			return false
		}
		count += n
	}
	return true
}

// readMessage reads the next Message from conn
func readMessage(conn net.Conn) *TLVMessage {
	// read header from connection
	headerBytes := readBytes(conn, MessageHeaderLength)
	if headerBytes == nil {
		return nil
	}

	// parse header
	typ := headerBytes[0]
	length := binary.BigEndian.Uint16(headerBytes[1:3])

	// make sure message type is valid
	if typ == MessageTypeNone || typ >= MessageTypeInvalid {
		return nil
	}

	// make sure message length is valid
	if length > MessageMaxLength {
		return nil
	}

	// read remaining message data from connection
	data := readBytes(conn, length-MessageHeaderLength)
	if data == nil {
		return nil
	}

	// return message
	return &TLVMessage{
		typ,
		length,
		data,
	}
}

// writeMessage writes message to conn
func writeMessage(conn net.Conn, message Message) bool {
	b, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}
	tlv := TLVMessage{
		message.GetType(),
		uint16(len(b)) + MessageHeaderLength,
		b,
	}
	buf := tlv.serialize()
	return writeBytes(conn, buf)
}

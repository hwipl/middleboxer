package cmd

import (
	"bytes"
	"encoding/binary"
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

// Message is a TLV message
type Message struct {
	Type   uint8
	Length uint16
	Data   []byte
}

// serialize encodes a message as bytes
func (m *Message) serialize() []byte {
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
func readMessage(conn net.Conn) *Message {
	// read header from connection
	headerBytes := readBytes(conn, MessageHeaderLength)
	if headerBytes == nil {
		return nil
	}

	// parse header
	typ := headerBytes[0]
	length := binary.BigEndian.Uint16(headerBytes[1:3])

	// TODO: check types?

	// check length and read message data from connection
	if length > MessageMaxLength {
		return nil
	}
	data := readBytes(conn, length-MessageHeaderLength)
	if data == nil {
		return nil
	}

	// return message
	return newMessage(typ, data)
}

// writeMessage writes message to conn
func writeMessage(conn net.Conn, message *Message) bool {
	buf := message.serialize()
	return writeBytes(conn, buf)
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

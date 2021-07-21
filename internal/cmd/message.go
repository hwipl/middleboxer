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

package cmd

import (
	"bytes"
	"log"
	"net"
	"testing"
	"time"
)

// TestReadWritebytes tests the readBytes() and writeBytes() functions
func TestReadWriteBytes(t *testing.T) {
	in, out := net.Pipe()
	if err := out.SetDeadline(time.Now().Add(time.Second)); err != nil {
		log.Fatal(err)
	}
	data := []byte{1, 2, 3, 4, 5, 6}
	go func() {
		if !writeBytes(in, data) {
			log.Fatal("error writing to conn")
		}
	}()
	want := data
	got := readBytes(out, uint16(len(data)))
	if bytes.Compare(got, want) != 0 {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestReadWriteMessage tests the readMessage() and writeMessage() functions
func TestReadWriteMessage(t *testing.T) {
	in, out := net.Pipe()
	if err := out.SetDeadline(time.Now().Add(time.Second)); err != nil {
		log.Fatal(err)
	}
	msg := &MessageNop{}
	go func() {
		if !writeMessage(in, msg) {
			log.Fatal("error writing to conn")
		}
	}()
	want := msg
	got := readMessage(out)
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestMessageNop tests nop messages
func TestMessageNop(t *testing.T) {
	in, out := net.Pipe()
	if err := out.SetDeadline(time.Now().Add(time.Second)); err != nil {
		log.Fatal(err)
	}

	msg := &MessageNop{}
	go func() {
		if !writeMessage(in, msg) {
			log.Fatal("error writing to conn")
		}
	}()
	want := msg.GetType()
	got := readMessage(out).GetType()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

// TestMessageRegister tests register messages
func TestMessageRegister(t *testing.T) {
	in, out := net.Pipe()
	if err := out.SetDeadline(time.Now().Add(time.Second)); err != nil {
		log.Fatal(err)
	}

	msg := &MessageRegister{Client: 1}
	go func() {
		if !writeMessage(in, msg) {
			log.Fatal("error writing to conn")
		}
	}()
	want := msg
	got := readMessage(out)
	if got.GetType() != want.GetType() {
		t.Errorf("got %d, want %d", got.GetType(), want.GetType())
	}
	if *got.(*MessageRegister) != *want {
		t.Errorf("got %v, want %v", got, want)
	}
}

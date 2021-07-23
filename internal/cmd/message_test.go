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
	out.SetDeadline(time.Now().Add(time.Second))
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

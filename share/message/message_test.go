package message

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSend(t *testing.T) {
	sampleData := []byte("Hello")
	var buffer bytes.Buffer
	buffer.Grow(HEADER_SIZE + len(sampleData) + 1)

	msgchan := NewMessageChan(&buffer)
	msgchan.Send(Kind(0x33), sampleData)

	expected := []byte("\x33Hello\n")
	assert.Equal(t, expected, buffer.Bytes())
}

func TestReceive(t *testing.T) {
	sampleMessage := []byte("\x33Hello\n")
	var buffer bytes.Buffer
	buffer.Write(sampleMessage)

	msgchan := NewMessageChan(&buffer)
	kind, data, _ := msgchan.Receive()

	assert.Equal(t, Kind(0x33), kind)
	assert.Equal(t, sampleMessage[1:len(sampleMessage)-1], data)
}

func TestConnSend(t *testing.T) {
}

func TestConnReceive(t *testing.T) {}

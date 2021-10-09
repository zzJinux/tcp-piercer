package message

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zzJinux/tcp-piercer/share"
)

func TestConnSendReceive(t *testing.T) {

	port := share.AvailablePort()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Fatalf("prepare: listen failed: %v", err)
	}

	// sampleData := []byte("Hello")
	// expected := []byte("\x33Hello\n")

	conn1, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Fatalf("prepare: listen failed: %v", err)
	}
	defer conn1.Close()

	conn2, err := l.Accept()
	if err != nil {
		t.Fatalf("prepare: accept failed: %v", err)
	}
	defer conn2.Close()

	done1, done2 := make(chan struct{}), make(chan struct{})

	go func() {
		msgchan := NewMessageChan(conn1)
		msgchan.Send(Kind(0x33), []byte("Hello"))
		close(done1)
	}()

	var kind Kind
	var data []byte
	go func() {
		msgchan := NewMessageChan(conn2)
		kind, data, _ = msgchan.Receive()

		close(done2)
	}()

	<-done1
	<-done2

	assert.Equal(t, Kind(0x33), kind)
	assert.Equal(t, []byte("Hello"), data)
}

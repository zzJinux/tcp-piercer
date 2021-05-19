package pnet

import (
	"bufio"
	"bytes"
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	var err error
	testServerAddress, err = net.ResolveTCPAddr("tcp", os.Getenv("TEST_SERVER"))
	if err != nil {
		log.Fatalf("failed to resolve the address: %v", err)
	}
	code := m.Run()
	os.Exit(code)
}

var testServerAddress *net.TCPAddr

func getTestServerAddress() *net.TCPAddr {
	return testServerAddress
}

// See "test/proxy_net_test.sh" for how the environment is set up
func TestDialEchoServer(t *testing.T) {
	assert := assert.New(t)

	//
	// The actual value of ServicePort doesn't matter
	// Use random ports for ServicePort to avoid connection timeout due to the router's conntrack
	// 120s is typical timeout duration.
	//
	dialer := Dialer{ServicePort: 9000 + int(time.Now().Unix()%120), Dialer: net.Dialer{Timeout: time.Duration(5 * time.Second)}}
	address := getTestServerAddress()

	ctx := context.TODO()
	conn, err := dialer.DialContext(ctx, address.Network(), address.String())
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	defer func() {
		// close the connection
		err = conn.Close()
		if err != nil {
			t.Fatalf("Close failed: %v", err)
		}
	}()

	bufConn := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	// receive the initial message from the server
	initialMessage, err := bufConn.ReadBytes('\n')
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	initialMessage = initialMessage[:len(initialMessage)-1]
	tokens := bytes.Split(initialMessage, []byte(" "))
	assert.Equal(4, len(tokens), "the number of tokens of initial message")
	assert.Equal(address.IP.String(), string(tokens[2]))
	assert.Equal(strconv.Itoa(address.Port), string(tokens[3]))

	// send a random message to the server
	_, err = bufConn.Write([]byte("hello\n"))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	err = bufConn.Flush()
	if err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	// receive the corresponding reply
	reply, err := bufConn.ReadBytes('\n')
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	assert.Equal("ECHO: hello\n", string(reply))

}

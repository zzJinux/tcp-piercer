package pnet

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// See "test/proxy_net_test.sh" to find out how the "TestServer" is set up
var echoServerAddr *net.TCPAddr

func TestMain(m *testing.M) {
	var err error
	echoServerAddr, err = net.ResolveTCPAddr("tcp", os.Getenv("ECHO_SERVER"))
	if err != nil {
		log.Fatalf("failed to resolve the address: %v", err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestDialEchoServer(t *testing.T) {
	assert := assert.New(t)

	var err error
	var done = make(chan struct{}) // close on test done
	defer func() { close(done) }()

	//
	// Use random ports for ServicePort to avoid connection timeout due to the router's conntrack
	// 120s is typical timeout duration.
	//
	servicePort := 9000 + int(time.Now().Unix()%120)
	var ready = make(chan struct{})
	go simpleServe(done, ready, servicePort)
	<-ready

	normalDialer := net.Dialer{Timeout: time.Duration(5 * time.Second)}
	normalDialer.LocalAddr, err = net.ResolveTCPAddr("tcp", fmt.Sprint(":", servicePort))
	_, err = normalDialer.Dial(echoServerAddr.Network(), echoServerAddr.String())
	assert.NotNil(err, "Dial call should fail.")

	proxyDialer := Dialer{ServicePort: servicePort, Dialer: net.Dialer{Timeout: time.Duration(5 * time.Second)}}
	proxyConn, err := proxyDialer.Dial(echoServerAddr.Network(), echoServerAddr.String())
	assert.Nilf(err, "proxy dial failed: %v", err)

	defer func() {
		err = proxyConn.Close()
		assert.Nilf(err, "proxy conn close failed: %v", err)
	}()

	echoTest(t, proxyConn)
}

//
// Prepare a listening socket bound to SERVICE_PORT
//
func simpleServe(done <-chan struct{}, ready chan<- struct{}, port int) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	defer l.Close()

	var testDone bool
	go func() {
		<-done
		testDone = true
		l.Close()
	}()

	close(ready)

	for {
		conn, err := l.Accept()
		if err != nil {
			if testDone {
				return nil
			} else {
				return err
			}
		}

		go func(c net.Conn) {
			io.Copy(c, c)
			c.Close()
		}(conn)
	}
}

func echoTest(t *testing.T, conn net.Conn) {
	assert := assert.New(t)

	remoteAddr := conn.RemoteAddr().(*net.TCPAddr)
	bufConn := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	// receive the initial message from the server
	initialMessage, err := bufConn.ReadBytes('\n')
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	initialMessage = initialMessage[:len(initialMessage)-1]
	tokens := bytes.Split(initialMessage, []byte(" "))

	assert.Equal(4, len(tokens), "the number of tokens of initial message")
	assert.Equal(remoteAddr.IP.String(), string(tokens[2]))
	assert.Equal(strconv.Itoa(remoteAddr.Port), string(tokens[3]))

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

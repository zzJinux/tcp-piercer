package test

import (
	"context"
	"encoding/json"
	"flag"
	"net"
	"os"
	"testing"
	"time"

	"github.com/zzJinux/tcp-piercer/client"
	"github.com/zzJinux/tcp-piercer/server"
)

// client flag
var serverAddr = flag.String("serveraddr", "", "address of the server (usage: -serveraddr=1.2.3.4:56")

// server flag
var serverPort = flag.String("port", "", "port which the server listens on")

var resultpath = flag.String("resultpath", "", "path for writing the result to")

var role = flag.String("role", "", "client|server")

func TestInitChannel(t *testing.T) {
	if *resultpath == "" {
		t.Error("-resultpath not specified")
	}

	switch *role {
	case "client":
		runClient(t)
	case "server":
		runServer(t)
	default:
		t.Skipf("unknown role\"%s\"", *role)
	}
}

func runClient(t *testing.T) {
	if *serverAddr == "" {
		t.Error("-serveraddr not specified")
	}

	var servicePort = 8080

	c, _ := client.NewClient(*serverAddr, servicePort)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	session, err := c.InitChannel(ctx)
	if err != nil {
		t.Fatal(err)
	}
	out := make(map[string]string)
	out["private"] = session.PrivateEndpoint.String()
	out["public"] = session.PublicEndpoint.String()

	b, _ := json.Marshal(out)
	os.WriteFile(*resultpath, b, 0666)
}

func runServer(t *testing.T) {
	if *serverPort == "" {
		t.Fatal("-serverPort not specified")
	}
	s, _ := server.NewServer()

	l, _ := net.Listen("tcp", ":"+*serverPort)
	defer l.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	conn, _ := l.Accept()
	defer conn.Close()

	session, err := s.InitChannel(ctx, conn)
	if err != nil {
		t.Fatal(err)
	}
	out := make(map[string]string)
	out["private"] = session.PrivateEndpoint.String()
	out["public"] = session.PublicEndpoint.String()

	b, _ := json.Marshal(out)
	os.WriteFile(*resultpath, b, 0666)
}

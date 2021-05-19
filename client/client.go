package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Client struct {
	conn        net.Conn
	serverAddr  string
	servicePort string
}

func NewClient(serverAddr, servicePort string) (*Client, error) {
	return &Client{nil, serverAddr, servicePort}, nil
}

func (c *Client) Start() (err error) {
	var conn net.Conn
	conn, err = net.Dial("tcp", c.serverAddr)
	if err != nil {
		return
	}

	msg := "My private endpoint is " + conn.LocalAddr().String()
	fmt.Fprintln(conn, msg)
	replyBytes, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		return
	}

	log.Println("Message")
	log.Println("\t" + msg)
	log.Println("Reply")
	log.Println("\t" + string(replyBytes))

	return
}

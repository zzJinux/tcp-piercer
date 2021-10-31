package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/zzJinux/tcp-piercer/share/control"
	"github.com/zzJinux/tcp-piercer/share/message"
	"github.com/zzJinux/tcp-piercer/share/pnet"
	"golang.org/x/sync/errgroup"
)

type Client struct {
	serverAddr  string
	servicePort int

	session *control.Session
}

func NewClient(serverAddr string, servicePort int) (*Client, error) {
	if servicePort == 0 {
		return nil, fmt.Errorf("NewClient: 0 cannot be the service port")
	}
	return &Client{
		serverAddr:  serverAddr,
		servicePort: servicePort,
	}, nil
}

// non-blocking
func (c *Client) StartContext(ctx context.Context) error {
	var err error

	log.Printf("connecting to %s", c.serverAddr)

	// 1. connection
	session, err := c.InitChannel(ctx)
	if err != nil {
		return err
	}
	c.session = session

	log.Printf("connected")
	// 2. handle incomding payloads (encapsulated TCP requests)
	// TODO

	return nil
}

func (c *Client) InitChannel(ctx context.Context) (*control.Session, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("connection: cancelled")
	default:
		// normal
	}

	dialer := pnet.Dialer{ServicePort: c.servicePort, Dialer: net.Dialer{}}
	conn, err := dialer.DialContext(ctx, "tcp", c.serverAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	msgChan := message.NewMessageChan(conn)

	var publicEndpoint, privateEndpoint *net.TCPAddr = nil, conn.LocalAddr().(*net.TCPAddr)

	// The CLIENT informs the SERVER of the address before SNAT
	eg := new(errgroup.Group)

	eg.Go(func() error {
		return msgChan.Send(message.MSG_NATINFO, []byte(privateEndpoint.String()))
	})

	// The SERVER informs the CLIENT of the address after SNAT
	eg.Go(func() error {
		kind, data, err := msgChan.Receive()
		if err != nil {
			return err
		}

		if kind != message.MSG_NATINFO {
			return fmt.Errorf("unexpected message kind: %s", kind.String())
		}

		dataStr := string(data)
		publicEndpoint, err = net.ResolveTCPAddr("tcp", dataStr)
		if err != nil {
			return err
		}

		return nil
	})

	err = eg.Wait()
	if err != nil {
		return nil, fmt.Errorf("initConnection: %w", err)
	}

	return &control.Session{
		Conn:            conn,
		MsgChan:         &msgChan,
		PublicEndpoint:  publicEndpoint,
		PrivateEndpoint: privateEndpoint,
	}, nil
}

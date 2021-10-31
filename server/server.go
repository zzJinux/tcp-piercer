package server

import (
	"context"
	"fmt"
	"net"

	"github.com/zzJinux/tcp-piercer/share/control"
	"github.com/zzJinux/tcp-piercer/share/message"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	sessions control.Sessions
}

func NewServer() (*Server, error) {
	return &Server{}, nil
}

// `StartContext` starts the server (non-blocking)
func (s *Server) StartContext(ctx context.Context, port int) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		// TODO: error handling
		go s.handleConnect(ctx, conn)
	}
}

func (s *Server) handleConnect(ctx context.Context, conn net.Conn) error {
	session, err := s.InitChannel(ctx, conn)
	if err != nil {
		return err
	}

	fmt.Println(session)
	return nil
}

func (s *Server) InitChannel(ctx context.Context, conn net.Conn) (*control.Session, error) {
	var messageChan = message.NewMessageChan(conn)
	var publicEndpoint, privateEndpoint *net.TCPAddr = conn.RemoteAddr().(*net.TCPAddr), nil

	eg := new(errgroup.Group)
	eg.Go(func() error {
		return messageChan.Send(message.MSG_NATINFO, []byte(publicEndpoint.String()))
	})

	eg.Go(func() error {
		kind, data, err := messageChan.Receive()
		if err != nil {
			return err
		}

		if kind != message.MSG_NATINFO {
			return fmt.Errorf("unexpected message kind: %s", kind.String())
		}

		dataStr := string(data)
		privateEndpoint, err = net.ResolveTCPAddr("tcp", dataStr)
		if err != nil {
			return err
		}

		return nil
	})

	err := eg.Wait()
	if err != nil {
		return nil, fmt.Errorf("initConnection: %w", err)
	}

	return &control.Session{
		Conn:            conn,
		MsgChan:         &messageChan,
		PrivateEndpoint: privateEndpoint,
		PublicEndpoint:  publicEndpoint,
	}, nil
}

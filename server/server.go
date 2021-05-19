package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Server struct {
	listener net.Listener
	port     string
}

func NewServer(port string) (*Server, error) {
	return &Server{nil, port}, nil
}

func (s *Server) Start() (err error) {
	var listener net.Listener
	listener, err = net.Listen("tcp", ":"+s.port)
	if err != nil {
		return
	}
	s.listener = listener

	for {
		var conn net.Conn
		conn, err = listener.Accept()
		if err != nil {
			return
		}
		// ctx := context.Background()

		done := make(chan error)
		go func() {
			msg, err := bufio.NewReader(conn).ReadBytes('\n')
			if err != nil {
				done <- err
				return
			}
			// TODO: validate the message

			reply := "Your public endpoint is " + conn.RemoteAddr().String()
			fmt.Fprintln(conn, reply)

			log.Println("Message")
			log.Println("\t" + string(msg))
			log.Println("Reply")
			log.Println("\t" + reply)
		}()
	}
}

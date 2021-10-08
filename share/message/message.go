package message

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

type Kind int

const HEADER_SIZE = 1

const (
	MSG_NIL Kind = iota
	MSG_NATINFO
	MSG_PACKET
)

func (msgt Kind) String() string {
	switch msgt {
	case MSG_NIL:
		return "NIL"
	case MSG_NATINFO:
		return "NATINFO"
	case MSG_PACKET:
		return "PACKET"
	default:
		return fmt.Sprintf("UNKNOWN-%d", int(msgt))
	}
}

func NewMessageChan(rw io.ReadWriter) MessageChan {
	bch := &BaseMessageChan{
		rw:  rw,
		brw: bufio.NewReadWriter(bufio.NewReader(rw), bufio.NewWriter(rw)),
	}

	if conn, ok := rw.(net.Conn); ok {
		return &ConnMessageChan{
			imessageChan: bch,
			conn:         conn,
		}
	} else {
		return bch
	}
}

type MessageChan interface {
	Send(Kind, []byte) error
	Receive() (Kind, []byte, error)
}

type BaseMessageChan struct {
	rw  io.ReadWriter // Underlying ReadWriter
	brw *bufio.ReadWriter
}

func (m *BaseMessageChan) Receive() (Kind, []byte, error) {
	message, err := m.brw.ReadBytes('\n')
	if err != nil {
		return MSG_NIL, nil, err
	}
	kind := Kind(message[0])
	data := message[1 : len(message)-1]

	return kind, data, nil
}

func (m *BaseMessageChan) Send(kind Kind, data []byte) error {
	buf := make([]byte, HEADER_SIZE+len(data)+1)
	buf[0] = byte(kind)
	copy(buf[1:], data)
	buf[1+len(data)] = '\n'

	_, err := m.brw.Write(buf)
	if err != nil {
		return err
	}
	if err = m.brw.Flush(); err != nil {
		return err
	}
	return nil
}

type ConnMessageChan struct {
	imessageChan MessageChan
	conn         net.Conn
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func (m *ConnMessageChan) Receive() (Kind, []byte, error) {
	if err := m.conn.SetReadDeadline(time.Now().Add(m.readTimeout)); err != nil {
		return MSG_NIL, nil, err
	}
	return m.imessageChan.Receive()
}

func (m *ConnMessageChan) Send(kind Kind, data []byte) error {
	if err := m.conn.SetWriteDeadline(time.Now().Add(m.writeTimeout)); err != nil {
		return err
	}
	return m.imessageChan.Send(kind, data)
}

func (m *ConnMessageChan) SetReadTimeout(d time.Duration) {
	m.readTimeout = d
}

func (m *ConnMessageChan) SetWriteTimeout(d time.Duration) {
	m.writeTimeout = d
}

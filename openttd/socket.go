package openttd

import (
	"fmt"
	"net"
)

type SocketConn struct {
	Socket net.Conn
}

func (s *SocketConn) Send(p Packet) error {
	b := p.Bytes()
	_, err := s.Socket.Write(b)

	return err
}

func SocketConnect(h string, p int) (SocketConn, error) {
	s, err := net.Dial("tcp", fmt.Sprintf("%s:%d", h, p))

	c := SocketConn{
		Socket: s,
	}

	return c, err
}

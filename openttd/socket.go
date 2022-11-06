package openttd

import (
	"fmt"
	"net"
)

type SocketConn struct {
	socket net.Conn
}

func (s *SocketConn) Send(p *Packet) error {
	b := p.Bytes()
	_, err := s.socket.Write(b)

	return err
}

func OpenSocket(h string, p int) (*SocketConn, error) {
	s, err := net.Dial("tcp", fmt.Sprintf("%s:%d", h, p))

	c := SocketConn{
		socket: s,
	}

	return &c, err
}

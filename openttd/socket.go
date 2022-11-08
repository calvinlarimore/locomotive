package openttd

import (
	"fmt"
	"net"
)

type socketConn struct {
	socket net.Conn
}

func (s *socketConn) send(p *Packet) error {
	b := p.Bytes()
	_, err := s.socket.Write(b)

	return err
}

func OpenSocket(h string, p int) (*socketConn, error) {
	s, err := net.Dial("tcp", fmt.Sprintf("%s:%d", h, p))

	c := socketConn{
		socket: s,
	}

	return &c, err
}

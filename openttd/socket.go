package openttd

import (
	"encoding/binary"
	"fmt"
	"net"
)

type conn struct {
	socket  net.Conn
	channel chan *Packet
}

func (s *conn) send(p *Packet) error {
	b := p.Bytes()
	_, err := s.socket.Write(b)

	return err
}

func (s *conn) read() (*Packet, error) {
	b := make([]byte, 69)
	_, err := s.socket.Read(b)

	l := binary.LittleEndian.Uint16(b[1:3])

	p := createPacket(b[0])
	p.data = append(p.data, b[3:l]...)

	return &p, err
}

func OpenSocket(h string, p int, ch chan *Packet) (*conn, error) {
	s, err := net.Dial("tcp", fmt.Sprintf("%s:%d", h, p))

	c := conn{
		socket:  s,
		channel: ch,
	}

	go func(c *conn, ch chan *Packet) {
		for {
			p, err := c.read()
			ch <- p

			if err != nil {
				return
				// Properly handle error
			}
		}
	}(&c, ch)

	return &c, err
}

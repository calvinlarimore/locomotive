package openttd

import (
	"encoding/binary"
	"fmt"
	"net"
)

type socket struct {
	conn    net.Conn
	channel chan *Packet
}

func (s *socket) send(p *Packet) error {
	b := p.Bytes()
	_, err := s.conn.Write(b)

	return err
}

func (s *socket) read() (*Packet, error) {
	b := make([]byte, 69)
	_, err := s.conn.Read(b)

	l := binary.LittleEndian.Uint16(b[1:3])

	p := createPacket(b[0])
	p.data = append(p.data, b[3:l]...)

	return &p, err
}

func openSocket(h string, p int, ch chan *Packet) (*socket, error) {
	c, err := net.Dial("tcp", fmt.Sprintf("%s:%d", h, p))

	s := socket{
		conn:    c,
		channel: ch,
	}

	go func(s *socket, c chan *Packet) {
		for {
			p, err := s.read()
			c <- p

			if err != nil {
				return
				// Properly handle error
			}
		}
	}(&s, ch)

	return &s, err
}

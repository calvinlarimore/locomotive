package locomotive

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/calvinlarimore/locomotive/openttd"
)

type socket struct {
	conn net.Conn
}

func (s *socket) send(p *openttd.Packet) error {
	b := p.Bytes()
	_, err := s.conn.Write(b)

	return err
}

func (s *socket) read() (*openttd.Packet, error) {
	b := make([]byte, 32767)
	_, err := s.conn.Read(b)

	l := binary.LittleEndian.Uint16(b[0:2])

	p := openttd.CreatePacket(b[2])
	p.Data = append(p.Data, b[3:l]...)

	return p, err
}

func openSocket(h string, p int) (*socket, error) {
	c, err := net.Dial("tcp", fmt.Sprintf("%s:%d", h, p))

	s := socket{
		conn: c,
	}

	go func(s *socket) {
		for {
			p, err := s.read()
			handlePacket(p)

			if err != nil {
				return
				// Properly handle error
			}
		}
	}(&s)

	return &s, err
}

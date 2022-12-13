package locomotive

import "github.com/calvinlarimore/locomotive/openttd"

type Conn struct {
	socket *socket
}

func (c *Conn) Send(m openttd.ClientMessage) {
	c.socket.send(m.Packet())
}

func Connect(h string, p int) (*Conn, error) {
	s, err := openSocket(h, p)

	c := Conn{
		socket: s,
	}

	return &c, err
}

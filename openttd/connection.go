package openttd

type Conn struct {
	socket *socket
}

func (c *Conn) Send(m ClientMessage) {
	c.socket.send(m.packet())
}

func Connect(h string, p int) (*Conn, error) {
	s, err := openSocket(h, p)

	c := Conn{
		socket: s,
	}

	return &c, err
}

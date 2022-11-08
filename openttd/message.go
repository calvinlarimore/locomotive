package openttd

type ServerMessage interface {
	Type() byte
}

type MessageHandler[T ServerMessage] interface {
	Handle(T)
}

var messageHandlers = make(map[string]MessageHandler[ServerMessage])

type MessageServerError struct {
	Error byte
}

func (m *MessageServerError) Type() byte {
	return 0x03
}

func createMessageServerError(p *Packet) *MessageServerError {
	m := MessageServerError{}

	b, _ := p.Reader().ReadByte()
	m.Error = b

	return &m
}

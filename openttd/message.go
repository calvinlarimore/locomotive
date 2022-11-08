package openttd

type ServerMessage interface {
	Type() byte
}

type ClientMessage interface {
	Type() byte
	packet() *Packet
}

type MessageHandler[T ServerMessage] interface {
	Handle(T)
}

var messageHandlers = make(map[string]MessageHandler[ServerMessage])

// Messages:

// PACKET_SERVER_FULL
type MessageServerFull struct{}

func (m *MessageServerFull) Type() byte { return 0x00 }

func createMessageServerFull(p *Packet) *MessageServerFull { return &MessageServerFull{} }

// PACKET_SERVER_BANNED
type MessageServerBanned struct{}

func (m *MessageServerBanned) Type() byte { return 0x01 }

func createMessageServerBanned(p *Packet) *MessageServerBanned { return &MessageServerBanned{} }

// PACKET_SERVER_ERROR
type MessageServerError struct {
	Error byte
}

func (m *MessageServerError) Type() byte { return 0x03 }

func createMessageServerError(p *Packet) *MessageServerError {
	m := MessageServerError{}

	b, _ := p.Reader().ReadByte()
	m.Error = b

	return &m
}

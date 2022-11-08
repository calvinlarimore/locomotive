package openttd

type ServerMessage interface {
	Type() byte
}

type ClientMessage interface {
	packet() *packet
}

type MessageHandler[T ServerMessage] interface {
	Handle(T)
}

var messageHandlers = make(map[string]MessageHandler[ServerMessage])

func SetMessageHandler(m string, h MessageHandler[ServerMessage]) {
	messageHandlers[m] = h
}

// Server Messages:

// PACKET_SERVER_FULL
type MessageServerFull struct{}

func (m *MessageServerFull) Type() byte { return 0x00 }

func createMessageServerFull(p *packet) *MessageServerFull { return &MessageServerFull{} }

// PACKET_SERVER_BANNED
type MessageServerBanned struct{}

func (m *MessageServerBanned) Type() byte { return 0x01 }

func createMessageServerBanned(p *packet) *MessageServerBanned { return &MessageServerBanned{} }

// PACKET_SERVER_ERROR
type MessageServerError struct {
	Error byte
}

func (m *MessageServerError) Type() byte { return 0x03 }

func createMessageServerError(p *packet) *MessageServerError {
	m := MessageServerError{}

	b, _ := p.Reader().ReadByte()
	m.Error = b

	return &m
}

// Client Messages:

// PACKET_CLIENT_JOIN
type MessageClientJoin struct {
	Name    string
	Company byte
}

func (m *MessageClientJoin) packet() *packet {
	p := createPacket(0x02)

	p.Writer().WriteString(gameVersion)
	p.Writer().WriteUint32(newGRFRevision)
	p.Writer().WriteString(m.Name)
	p.Writer().WriteByte(m.Company)
	p.Writer().WriteByte(0x00) // LEGACY: Used to contain language id

	return &p
}

package openttd

import (
	"fmt"
	"log"
)

type ServerMessage interface {
	Type() byte
}

type ClientMessage interface {
	Packet() *Packet
}

// Server Messages:

// PACKET_SERVER_FULL
type MessageServerFull struct{}

func (MessageServerFull) Type() byte { return 0x00 }

func CreateMessageServerFull(r *packetReader) *MessageServerFull { return &MessageServerFull{} }

// PACKET_SERVER_BANNED
type MessageServerBanned struct{}

func (MessageServerBanned) Type() byte { return 0x01 }

func CreateMessageServerBanned(r *packetReader) *MessageServerBanned { return &MessageServerBanned{} }

// PACKET_SERVER_ERROR
type MessageServerError struct {
	Error byte
}

func (MessageServerError) Type() byte { return 0x03 }

func CreateMessageServerError(r *packetReader) *MessageServerError {
	m := MessageServerError{}

	b, _ := r.ReadByte()
	m.Error = b

	return &m
}

// PACKET_SERVER_GAME_INFO
type MessageServerGameInfo struct {
	ProtocolVersion   byte
	GameScriptVersion uint32
	GameScript        string
	GRFs              []newGRFConfig
	GameDate          date
	StartDate         date
	MaxCompanies      byte
	ComapanyCount     byte
	ServerName        string
	ServerVersion     string
	IsUsingPassword   bool
	MaxClients        byte
	ClientCount       byte
	SpectatorCount    byte
	MapWidth          uint16
	MapHeight         uint16
	Landscape         byte
	IsDedicated       bool
}

func (MessageServerGameInfo) Type() byte { return 0x06 }

func CreateMessageServerGameInfo(r *packetReader) *MessageServerGameInfo {
	m := MessageServerGameInfo{}

	m.ProtocolVersion, _ = r.ReadByte()
	newGRFSerialization := byte(0x00)

	switch m.ProtocolVersion {
	case 0x06:
		newGRFSerialization, _ = r.ReadByte()

		fallthrough
	case 0x05:
		m.GameScriptVersion, _ = r.ReadUint32()
		m.GameScript, _ = r.ReadString(networkNameLength)

		fallthrough
	case 0x04:
		count, _ := r.ReadByte()

		readMD5 := func(r *packetReader) md5 {
			b := make([]byte, 16)
			r.Read(b)
			return *(*[16]byte)(b) // wtf go
		}

		for i := byte(0); i < count; i++ {
			newgrf := newGRFIdentifier{}

			switch newGRFSerialization {
			case 0x00: //NST_GRFID_MD5
				newgrf.ID, _ = r.ReadUint32()
				newgrf.MD5 = readMD5(r)
			case 0x01: //NST_GRFID_MD5_NAME
				newgrf.ID, _ = r.ReadUint32()
				newgrf.MD5 = readMD5(r)
				newgrf.Name, _ = r.ReadString(networkNameLength)
			case 0x02: // NST_LOOKUP_ID
				// TODO: (NEWGRF_IMPL) add lookup when implementing newgrfs
				index, _ := r.ReadUint32()
				log.Fatalln(fmt.Errorf("tried to get newgrf from lookup table. this is not implemented yet! index: %d", index))
			}

			config := newGRFConfig{}
			config.ID = newgrf

			m.GRFs = append(m.GRFs, config)

			log.Printf("\"Loaded\" NewGRF: %#v", newgrf)

			// TODO: (NEWGRF_IMPL) load newgrf
		}

		fallthrough
	case 0x03:
		{
			g, _ := r.ReadUint32()
			m.GameDate = date(g)
			s, _ := r.ReadUint32()
			m.StartDate = date(s)
		}

		fallthrough
	case 0x02:
		m.MaxCompanies, _ = r.ReadByte()
		m.ComapanyCount, _ = r.ReadByte()

		fallthrough
	case 0x01:
		m.ServerName, _ = r.ReadString(networkNameLength)
		m.ServerVersion, _ = r.ReadString(networkVersionLength)

		if m.ProtocolVersion < 0x06 {
			r.ReadByte() // LEGACY: Used to contain the server's language id
		}

		m.IsUsingPassword, _ = r.ReadBool()
		m.MaxClients, _ = r.ReadByte()
		m.ClientCount, _ = r.ReadByte()
		m.SpectatorCount, _ = r.ReadByte()

		if m.ProtocolVersion < 0x03 {
			// LEGACY: Used to contain the server's date info.
			// We still read it because we need it if we try to
			//   connect to an ancient server.

			g, _ := r.ReadUint16()
			m.GameDate = date(int(g) + daysUntilOriginalBaseYear())
			s, _ := r.ReadUint16()
			m.StartDate = date(int(s) + daysUntilOriginalBaseYear())
		}

		if m.ProtocolVersion < 0x06 {
			// LEGACY: Used to contain the server's map name.
			b := byte(0xff)
			for b != 0x00 {
				b, _ = r.ReadByte()
			}
		}

		m.MapWidth, _ = r.ReadUint16()
		m.MapHeight, _ = r.ReadUint16()
		m.Landscape, _ = r.ReadByte()
		m.IsDedicated, _ = r.ReadBool()
	}

	return &m
}

// PACKET_SERVER_NEED_GAME_PASSWORD
type MessageServerNeedGamePassword struct{}

func (MessageServerNeedGamePassword) Type() byte { return 0x0a }

func CreateMessageServerNeedGamePassword(r *packetReader) *MessageServerNeedGamePassword {
	m := MessageServerNeedGamePassword{}
	return &m
}

// PACKET_SERVER_WELCOME
type MessageServerWelcome struct {
	ClientID       uint32
	GenerationSeed uint32
	ServerID       string
}

func (MessageServerWelcome) Type() byte { return 0x0e }

func CreateMessageServerWelcome(r *packetReader) *MessageServerWelcome {
	m := MessageServerWelcome{}

	id, _ := r.ReadUint32()
	m.ClientID = id

	seed, _ := r.ReadUint32()
	m.GenerationSeed = seed

	sid, _ := r.ReadString(networkServerIDLength)
	m.ServerID = sid

	return &m
}

// Client Messages:

// PACKET_CLIENT_JOIN
type MessageClientJoin struct {
	Name    string
	Company byte
}

func (m *MessageClientJoin) Packet() *Packet {
	p := CreatePacket(0x02)

	p.Writer().WriteString(gameVersion)
	p.Writer().WriteUint32(newGRFRevision)
	p.Writer().WriteString(m.Name)
	p.Writer().WriteByte(m.Company)
	p.Writer().WriteByte(0x00) // LEGACY: Used to contain language id

	return p
}

// PACKET_CLIENT_GAME_INFO
type MessageClientGameInfo struct{}

func (m *MessageClientGameInfo) Packet() *Packet { return CreatePacket(0x07) }

// PACKET_CLIENT_GAME_INFO
type MessageClientGamePassword struct {
	Password string
}

func (m *MessageClientGamePassword) Packet() *Packet {
	p := CreatePacket(0x0b)

	p.Writer().WriteString(m.Password)

	return p
}

package openttd

import (
	"encoding/binary"
	"fmt"
	"log"
)

type packet struct {
	packetType byte
	data       []byte
}

type packetReader struct {
	packet *packet
}

func (r *packetReader) Read(b []byte) int {
	for i := range b {
		b[i] = r.packet.data[i]
	}

	r.packet.data = r.packet.data[len(b):]

	return len(b)
}

func (r *packetReader) ReadByte() (byte, error) {
	b := make([]byte, 1)

	l := r.Read(b)
	if l != 1 {
		return 0x00, fmt.Errorf("wrong byte count! expected %d, got: %d", 1, l)
	}

	return b[0], nil
}

func (r *packetReader) ReadUint16() (uint16, error) {
	b := make([]byte, 2)

	l := r.Read(b)
	if l != 2 {
		return uint16(0), fmt.Errorf("wrong byte count! expected %d, got: %d", 2, l)
	}

	return binary.LittleEndian.Uint16(b), nil
}

func (r *packetReader) ReadUint32() (uint32, error) {
	b := make([]byte, 4)

	l := r.Read(b)
	if l != 4 {
		return uint32(0), fmt.Errorf("wrong byte count! expected %d, got: %d", 4, l)
	}

	return binary.LittleEndian.Uint32(b), nil
}

func (r *packetReader) ReadUint64() (uint64, error) {
	b := make([]byte, 8)

	l := r.Read(b)
	if l != 8 {
		return uint64(0), fmt.Errorf("wrong byte count! expected %d, got: %d", 8, l)
	}

	return binary.LittleEndian.Uint64(b), nil
}

func (r *packetReader) ReadBool() (bool, error) {
	b, err := r.ReadByte()

	return b != 0x00, err
}

func (r *packetReader) ReadString(max uint) (string, int) {
	l := max
	b := make([]byte, max)

	for i := uint(0); i < max-1; i++ {
		b[i] = r.packet.data[i]
		if b[i] == 0x00 {
			l = i + 1
			break
		}
	}

	r.packet.data = r.packet.data[l:]

	return string(b[:l]), int(l)
}

type packetWriter struct {
	packet *packet
}

func (w *packetWriter) Write(b []byte) {
	w.packet.data = append(w.packet.data, b...)
}

func (w *packetWriter) WriteByte(b byte) error {
	w.packet.data = append(w.packet.data, b)

	return nil
}

func (w *packetWriter) WriteUint16(i uint16) error {
	w.packet.data = append(w.packet.data, binary.LittleEndian.AppendUint16(make([]byte, 0), i)...)

	return nil
}

func (w *packetWriter) WriteUint32(i uint32) error {
	w.packet.data = append(w.packet.data, binary.LittleEndian.AppendUint32(make([]byte, 0), i)...)

	return nil
}

func (w *packetWriter) WriteUint64(i uint64) error {
	w.packet.data = append(w.packet.data, binary.LittleEndian.AppendUint64(make([]byte, 0), i)...)

	return nil
}

func (w *packetWriter) WriteBool(b bool) error {
	var err error

	if b {
		err = w.WriteByte(0x01)
	} else {
		err = w.WriteByte(0x00)
	}

	return err
}

func (w *packetWriter) WriteString(s string) int {
	w.packet.data = append(w.packet.data, s...)
	w.packet.data = append(w.packet.data, 0x00)

	return len(s) + 1
}

func (p *packet) Reader() *packetReader {
	return &packetReader{
		packet: p,
	}
}

func (p *packet) Writer() *packetWriter {
	return &packetWriter{
		packet: p,
	}
}

func (p *packet) Bytes() []byte {
	b := make([]byte, 0)

	l := uint16(len(p.data) + 3)
	b = append(b, binary.LittleEndian.AppendUint16(make([]byte, 0), l)...)

	b = append(b, p.Type())

	b = append(b, p.data...)

	fmt.Printf("Converting packet to bytes: %d\n", b)

	return b
}

func (p *packet) Type() byte {
	return p.packetType
}

func createPacket(t byte) *packet {
	p := packet{
		packetType: t,
		data:       make([]byte, 0),
	}

	return &p
}

func handlePacket(p *packet) {
	switch p.Type() {
	case 0x00: // PACKET_SERVER_FULL
		m := createMessageServerFull(p.Reader())
		h, ok := messageHandlers["full"].(MessageHandler[MessageServerFull])

		if ok {
			h.Handle(m)
		} else {
			errInvalidHandler(m)
		}

	case 0x01: // PACKET_SERVER_BANNED
		m := createMessageServerBanned(p.Reader())
		h, ok := messageHandlers["banned"].(MessageHandler[MessageServerBanned])

		if ok {
			h.Handle(m)
		} else {
			errInvalidHandler(m)
		}

	case 0x03: // PACKET_SERVER_ERROR
		m := createMessageServerError(p.Reader())
		h, ok := messageHandlers["error"].(MessageHandler[MessageServerError])

		if ok {
			h.Handle(m)
		} else {
			errInvalidHandler(m)
		}

	case 0x06: // PACKET_SERVER_GAME_INFO
		m := createMessageServerGameInfo(p.Reader())
		h, ok := messageHandlers["game_info"].(MessageHandler[MessageServerGameInfo])

		if ok {
			h.Handle(m)
		} else {
			errInvalidHandler(m)
		}

	case 0x0a: // PACKET_SERVER_NEED_GAME_PASSWORD
		m := createMessageServerNeedGamePassword(p.Reader())
		h, ok := messageHandlers["need_game_password"].(MessageHandler[MessageServerNeedGamePassword])

		if ok {
			h.Handle(m)
		} else {
			errInvalidHandler(m)
		}

	case 0x0e: // PACKET_SERVER_WELCOME
		m := createMessageServerWelcome(p.Reader())
		h, ok := messageHandlers["welcome"].(MessageHandler[MessageServerWelcome])

		if ok {
			h.Handle(m)
		} else {
			errInvalidHandler(m)
		}

	default:
		log.Println(fmt.Errorf("unknown packet type 0x%02x", p.Type()))
	}
}

func errInvalidHandler(m ServerMessage) {
	log.Println(fmt.Errorf("message handler for message type \"%T\" has invalid type", m))
}

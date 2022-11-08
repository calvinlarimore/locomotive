package openttd

import (
	"encoding/binary"
	"fmt"
)

type Packet struct {
	packetType byte
	data       []byte
}

type PacketReader struct {
	packet *Packet
}

func (r *PacketReader) Read(b []byte) int {
	for i := range b {
		b[i] = r.packet.data[i]
	}

	r.packet.data = r.packet.data[len(b):]

	return len(b)
}

func (r *PacketReader) ReadByte() (byte, error) {
	b := make([]byte, 1)

	l := r.Read(b)
	if l != 1 {
		return 0x00, fmt.Errorf("wrong byte count! expected %d, got: %d", 1, l)
	}

	return b[0], nil
}

func (r *PacketReader) ReadUint16() (uint16, error) {
	b := make([]byte, 2)

	l := r.Read(b)
	if l != 2 {
		return uint16(0), fmt.Errorf("wrong byte count! expected %d, got: %d", 2, l)
	}

	return binary.LittleEndian.Uint16(b), nil
}

func (r *PacketReader) ReadUint32() (uint32, error) {
	b := make([]byte, 4)

	l := r.Read(b)
	if l != 4 {
		return uint32(0), fmt.Errorf("wrong byte count! expected %d, got: %d", 4, l)
	}

	return binary.LittleEndian.Uint32(b), nil
}

func (r *PacketReader) ReadUint64() (uint64, error) {
	b := make([]byte, 8)

	l := r.Read(b)
	if l != 8 {
		return uint64(0), fmt.Errorf("wrong byte count! expected %d, got: %d", 8, l)
	}

	return binary.LittleEndian.Uint64(b), nil
}

func (r *PacketReader) ReadBool() (bool, error) {
	b, err := r.ReadByte()

	return b != 0x00, err
}

func (r *PacketReader) ReadString(max uint) (string, int) {
	l := max
	b := make([]byte, max)

	for i := uint(0); i < max-1; i++ {
		b[i] = r.packet.data[i]
		if b[i] == 0x00 {
			l = i + 1
		}
	}

	r.packet.data = r.packet.data[l:]

	return string(b[:l]), int(l)
}

type PacketWriter struct {
	packet *Packet
}

func (w *PacketWriter) Write(b []byte) {
	w.packet.data = append(w.packet.data, b...)
}

func (w *PacketWriter) WriteByte(b byte) error {
	w.packet.data = append(w.packet.data, b)

	return nil
}

func (w *PacketWriter) WriteUint16(i uint16) error {
	w.packet.data = append(w.packet.data, binary.LittleEndian.AppendUint16(make([]byte, 0), i)...)

	return nil
}

func (w *PacketWriter) WriteUint32(i uint32) error {
	w.packet.data = append(w.packet.data, binary.LittleEndian.AppendUint32(make([]byte, 0), i)...)

	return nil
}

func (w *PacketWriter) WriteUint64(i uint64) error {
	w.packet.data = append(w.packet.data, binary.LittleEndian.AppendUint64(make([]byte, 0), i)...)

	return nil
}

func (w *PacketWriter) WriteBool(b bool) error {
	var err error

	if b {
		err = w.WriteByte(0x01)
	} else {
		err = w.WriteByte(0x00)
	}

	return err
}

func (w *PacketWriter) WriteString(s string) int {
	w.packet.data = append(w.packet.data, s...)
	w.packet.data = append(w.packet.data, 0x00)

	return len(s) + 1
}

func (p *Packet) Reader() *PacketReader {
	return &PacketReader{
		packet: p,
	}
}

func (p *Packet) Writer() *PacketWriter {
	return &PacketWriter{
		packet: p,
	}
}

func (p *Packet) Bytes() []byte {
	b := make([]byte, 3)
	b[0] = p.Type()

	// TODO: Length

	b = append(b, p.data...)

	return b
}

func (p *Packet) Type() byte {
	return p.packetType
}

func createPacket(t byte) Packet {
	p := Packet{
		packetType: t,
		data:       make([]byte, 0),
	}

	return p
}

func handlePacket(p *Packet) {
	switch p.Type() {
	case 0x00: // PACKET_SERVER_FULL
		m := createMessageServerError(p)
		messageHandlers["full"].Handle(m)
	case 0x01: // PACKET_SERVER_BANNED
		m := createMessageServerError(p)
		messageHandlers["banned"].Handle(m)
	case 0x03: // PACKET_SERVER_ERROR
		m := createMessageServerError(p)
		messageHandlers["error"].Handle(m)
	}
}

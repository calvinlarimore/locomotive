package openttd

import (
	"encoding/binary"
	"fmt"
)

type Packet struct {
	packetType byte
	Data       []byte
}

type packetReader struct {
	packet *Packet
}

func (r *packetReader) Read(b []byte) int {
	for i := range b {
		b[i] = r.packet.Data[i]
	}

	r.packet.Data = r.packet.Data[len(b):]

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
		b[i] = r.packet.Data[i]
		if b[i] == 0x00 {
			l = i + 1
			break
		}
	}

	r.packet.Data = r.packet.Data[l:]

	return string(b[:l]), int(l)
}

type packetWriter struct {
	packet *Packet
}

func (w *packetWriter) Write(b []byte) {
	w.packet.Data = append(w.packet.Data, b...)
}

func (w *packetWriter) WriteByte(b byte) error {
	w.packet.Data = append(w.packet.Data, b)

	return nil
}

func (w *packetWriter) WriteUint16(i uint16) error {
	w.packet.Data = append(w.packet.Data, binary.LittleEndian.AppendUint16(make([]byte, 0), i)...)

	return nil
}

func (w *packetWriter) WriteUint32(i uint32) error {
	w.packet.Data = append(w.packet.Data, binary.LittleEndian.AppendUint32(make([]byte, 0), i)...)

	return nil
}

func (w *packetWriter) WriteUint64(i uint64) error {
	w.packet.Data = append(w.packet.Data, binary.LittleEndian.AppendUint64(make([]byte, 0), i)...)

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
	w.packet.Data = append(w.packet.Data, s...)
	w.packet.Data = append(w.packet.Data, 0x00)

	return len(s) + 1
}

func (p *Packet) Reader() *packetReader {
	return &packetReader{
		packet: p,
	}
}

func (p *Packet) Writer() *packetWriter {
	return &packetWriter{
		packet: p,
	}
}

func (p *Packet) Bytes() []byte {
	b := make([]byte, 0)

	l := uint16(len(p.Data) + 3)
	b = append(b, binary.LittleEndian.AppendUint16(make([]byte, 0), l)...)

	b = append(b, p.Type())

	b = append(b, p.Data...)

	fmt.Printf("Converting packet to bytes: %d\n", b)

	return b
}

func (p *Packet) Type() byte {
	return p.packetType
}

func CreatePacket(t byte) *Packet {
	p := Packet{
		packetType: t,
		Data:       make([]byte, 0),
	}

	return &p
}

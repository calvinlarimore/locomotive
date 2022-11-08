package openttd

type Conn interface {
	send(*Packet) error
}

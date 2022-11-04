package openttd

type Conn interface {
	Send(Packet) error
}

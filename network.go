package locomotive

import (
	"fmt"
	"log"

	"github.com/calvinlarimore/locomotive/openttd"
)

type MessageHandler[T openttd.ServerMessage] interface {
	Handle(*T)
}

var messageHandlers = make(map[string]any)

func SetMessageHandler(m string, h any) {
	messageHandlers[m] = h
}

func handlePacket(p *openttd.Packet) {
	switch p.Type() {
	case 0x00: // PACKET_SERVER_FULL
		m := openttd.CreateMessageServerFull(p.Reader())
		h, ok := messageHandlers["full"].(MessageHandler[openttd.MessageServerFull])

		if ok {
			h.Handle(m)
		} else {
			errInvalidHandler(m)
		}

	case 0x01: // PACKET_SERVER_BANNED
		m := openttd.CreateMessageServerBanned(p.Reader())
		h, ok := messageHandlers["banned"].(MessageHandler[openttd.MessageServerBanned])

		if ok {
			h.Handle(m)
		} else {
			errInvalidHandler(m)
		}

	case 0x03: // PACKET_SERVER_ERROR
		m := openttd.CreateMessageServerError(p.Reader())
		h, ok := messageHandlers["error"].(MessageHandler[openttd.MessageServerError])

		if ok {
			h.Handle(m)
		} else {
			errInvalidHandler(m)
		}

	case 0x06: // PACKET_SERVER_GAME_INFO
		m := openttd.CreateMessageServerGameInfo(p.Reader())
		h, ok := messageHandlers["game_info"].(MessageHandler[openttd.MessageServerGameInfo])

		if ok {
			h.Handle(m)
		} else {
			errInvalidHandler(m)
		}

	case 0x0a: // PACKET_SERVER_NEED_GAME_PASSWORD
		m := openttd.CreateMessageServerNeedGamePassword(p.Reader())
		h, ok := messageHandlers["need_game_password"].(MessageHandler[openttd.MessageServerNeedGamePassword])

		if ok {
			h.Handle(m)
		} else {
			errInvalidHandler(m)
		}

	case 0x0e: // PACKET_SERVER_WELCOME
		m := openttd.CreateMessageServerWelcome(p.Reader())
		h, ok := messageHandlers["welcome"].(MessageHandler[openttd.MessageServerWelcome])

		if ok {
			h.Handle(m)
		} else {
			errInvalidHandler(m)
		}

	default:
		log.Println(fmt.Errorf("unknown packet type 0x%02x", p.Type()))
	}
}

func errInvalidHandler(m openttd.ServerMessage) {
	log.Println(fmt.Errorf("message handler for message type \"%T\" has invalid type", m))
}

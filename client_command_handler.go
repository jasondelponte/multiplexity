package multiplexity

import (
	"log"
)

type ClientCommandHandler struct {
	config         *Config
	serverMsgWrite MessageWriteFn
	clientMsgWrite MessageWriteFn
}

func NewClientCommandHandler(config *Config, serverMsgWrite, clientMsgWrite MessageWriteFn) *ClientCommandHandler {
	return &ClientCommandHandler{
		config:         config,
		serverMsgWrite: serverMsgWrite,
		clientMsgWrite: clientMsgWrite,
	}
}

func (c *ClientCommandHandler) Handle(command Command) {
	log.Println("Server command received", command.ToString())

	switch command.Type() {
	case CommandTypeClientConnect:

	case CommandTypeClientQuit:

	default:
		log.Fatalln("Unknown client command:", command.ToString())
	}
}

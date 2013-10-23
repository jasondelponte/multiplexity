package multiplexity

import (
	"log"
)

type ServerCommandHandler struct {
	config         *Config
	serverMsgWrite MessageHandlerFn
	clientMsgWrite MessageHandlerFn
}

func NewServerCommandHandler(config *Config, serverMsgWrite, clientMsgWrite MessageHandlerFn) *ServerCommandHandler {
	return &ServerCommandHandler{
		config:         config,
		serverMsgWrite: serverMsgWrite,
		clientMsgWrite: clientMsgWrite,
	}
}

func (s *ServerCommandHandler) Handle(command Command) {
	log.Println("Server command received", command.ToString())

	switch command.Type() {
	case CommandTypeServerConnect:
		s.sendUserCreds()

	case CommandTypeServerQuit:

	case CommandTypeServerMessage:
		s.handleServerMessage(command.(*ServerMessageCommand).Message)

	default:
		log.Fatalln("Unknown server command:", command.ToString())
	}
}

func (s *ServerCommandHandler) sendUserCreds() {
	userMsg := &Message{
		Command:  "USER",
		Params:   []string{s.config.Server.User, ".", "."},
		Trailing: s.config.Server.RealName,
	}
	nickMsg := &Message{
		Command: "NICK",
		Params:  []string{s.config.Server.Nick},
	}

	s.serverMsgWrite(userMsg)
	s.serverMsgWrite(nickMsg)
}

func (s *ServerCommandHandler) handleServerMessage(message *Message) {
	switch message.Command {
	case "PING":
		pongMsg := &Message{
			Command:  "PONG",
			Trailing: message.Trailing,
		}
		s.serverMsgWrite(pongMsg)

	// Capture these
	case RPL_WELCOME:
	case RPL_YOURHOST:
	case RPL_CREATED:
	case RPL_MYINFO:
	// case RPL_ISUPPORT:

	default:
		s.clientMsgWrite(message)
	}
}

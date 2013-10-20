package multiplexity

import (
	"log"
)

type ServerCommandHandler struct {
	config         *Config
	serverMsgWrite MessageWriteFn
	clientMsgWrite MessageWriteFn
}

func NewServerCommandHandler(config *Config, serverMsgWrite, clientMsgWrite MessageWriteFn) *ServerCommandHandler {
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

	case CommandTypeMessage:
		s.handleServerMessage(command.(*MessageCommand).Message)

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
	default:
		s.clientMsgWrite(message)
	}
}

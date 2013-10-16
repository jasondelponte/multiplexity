package multiplexity

import (
	"log"
)

type ServerCommandHandler struct {
	config       *Config
	serverCmd    CommandChan
	serverMsgOut MessageChan
	clientMsgOut MessageChan
}

func NewServerCommandHandler(config *Config, serverCmd CommandChan, serverMsgOut, clientMsgOut MessageChan) *ServerCommandHandler {
	return &ServerCommandHandler{
		config:       config,
		serverCmd:    serverCmd,
		serverMsgOut: serverMsgOut,
		clientMsgOut: clientMsgOut,
	}
}

func (s *ServerCommandHandler) Run() {
	for {
		command, _ := <-s.serverCmd
		s.handleServerCommand(command)
	}
}

func (s *ServerCommandHandler) handleServerCommand(command Command) {
	log.Println("Server command received", command.ToString())

	if command.Type() == CommandTypeServerConnect {
		s.sendUserCreds()
	} else if command.Type() == CommandTypeMessage {
		s.handleServerMessage(command.(*MessageCommand).Message)
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

	s.serverMsgOut <- userMsg
	s.serverMsgOut <- nickMsg
}

func (s *ServerCommandHandler) handleServerMessage(message *Message) {
	if message.Command == "PING" {
		pongMsg := &Message{
			Command:  "PONG",
			Trailing: message.Trailing,
		}
		s.serverMsgOut <- pongMsg
		return
	}

	s.clientMsgOut <- message
}

package multiplexity

import (
	"fmt"
	"log"
)

type ClientCommandHandler struct {
	config         *Config
	serverMsgWrite MessageHandlerFn
	clientMsgWrite MessageHandlerFn

	clients ClientList
}

func NewClientCommandHandler(config *Config, serverMsgWrite, clientMsgWrite MessageHandlerFn) *ClientCommandHandler {
	return &ClientCommandHandler{
		config:         config,
		serverMsgWrite: serverMsgWrite,
		clientMsgWrite: clientMsgWrite,

		clients: make(ClientList, 0, 5),
	}
}

func (c *ClientCommandHandler) Handle(command Command) {
	log.Println("Client command received", command.ToString())

	switch command.Type() {
	case CommandTypeClientConnect:

	case CommandTypeClientQuit:
		c.removeClient(command.(*ClientQuitCommand).Client)

	case CommandTypeClientMessage:
		c.handleClientMessage(command.(*ClientMessageCommand))

	default:
		log.Fatalln("Unknown client command:", command.ToString())
	}
}

func (c *ClientCommandHandler) handleClientMessage(command *ClientMessageCommand) {
	switch command.Message.Command {
	case "PING":
		pongMsg := &Message{
			Command:  "PONG",
			Trailing: command.Message.Trailing,
		}
		c.serverMsgWrite(pongMsg)

	case "NICK":
		if c.hasClient(command.Client) {
		} else {
			welcomeMsg := &Message{
				Command:  RPL_WELCOME,
				Trailing: fmt.Sprintf("Howdy %s", command.Message.Params[0]),
			}
			command.Client.Write(welcomeMsg)
			c.clients = append(c.clients, command.Client)
		}

	case "USER":
		// Ignore user command, only care about nick

	case "QUIT":
		c.removeClient(command.Client)
		quitMsg := &Message{
			Command:  "QUIT",
			Trailing: fmt.Sprintf("Bye %s!", "nick"),
		}
		command.Client.Write(quitMsg)

	default:
		c.serverMsgWrite(command.Message)
	}
}

func (c *ClientCommandHandler) hasClient(client *Client) bool {
	for _, v := range c.clients {
		if client.Id == v.Id {
			return true
		}
	}
	return false
}

func (c *ClientCommandHandler) removeClient(client *Client) {
	at := -1
	for i, v := range c.clients {
		if client.Id == v.Id {
			at = i
		}
	}

	if at >= 0 {
		c.clients = append(c.clients[:at], c.clients[at+1:]...)
	}
}

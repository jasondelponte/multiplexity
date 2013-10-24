package multiplexity

import (
	"fmt"
)

type CommandType int

const (
	CommandTypeClientConnect = iota
	CommandTypeClientQuit
	CommandTypeClientMessage

	CommandTypeServerConnect
	CommandTypeServerQuit
	CommandTypeServerMessage
)

type Command interface {
	Type() CommandType
	String() string
}

type CommandChan chan Command

// Client Connect Command
type ClientConnectCommand struct {
	Client *Client
}

func (c ClientConnectCommand) Type() CommandType {
	return CommandType(CommandTypeClientConnect)
}

func (c ClientConnectCommand) String() string {
	return fmt.Sprintf("ClientConnectCommand from: %d", c.Client.Id)
}

// Client Quit Command
type ClientQuitCommand struct {
	Client *Client
}

func (c ClientQuitCommand) Type() CommandType {
	return CommandType(CommandTypeClientQuit)
}

func (c ClientQuitCommand) String() string {
	return fmt.Sprintf("ClientQuitCommand from: %d", c.Client.Id)
}

// Mesage Command
type ClientMessageCommand struct {
	Message *Message
	Client  *Client
}

func (c ClientMessageCommand) Type() CommandType {
	return CommandType(CommandTypeClientMessage)
}

func (c ClientMessageCommand) String() string {
	return fmt.Sprintf("ClientMessageCommand from: %d Message: %s", c.Client.Id, c.Message)
}

// Server Connect Command
type ServerConnectCommand struct {
	Server *ServerClient
}

func (c ServerConnectCommand) Type() CommandType {
	return CommandType(CommandTypeServerConnect)
}

func (c ServerConnectCommand) String() string {
	return fmt.Sprintf("ServerConnectCommand")
}

// Server Quit Command
type ServerQuitCommand struct {
	Server *ServerClient
}

func (c ServerQuitCommand) Type() CommandType {
	return CommandType(CommandTypeServerQuit)
}

func (c ServerQuitCommand) String() string {
	return fmt.Sprintf("ServerQuitCommand")
}

// Mesage Command
type ServerMessageCommand struct {
	Message *Message
}

func (c ServerMessageCommand) Type() CommandType {
	return CommandType(CommandTypeServerMessage)
}

func (c ServerMessageCommand) String() string {
	return fmt.Sprintf("ServerMessageCommand Message: %s", c.Message)
}

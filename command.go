package multiplexity

import (
	"fmt"
)

type CommandType int

const (
	CommandTypeClientConnect = iota
	CommandTypeClientQuit

	CommandTypeServerConnect
	CommandTypeServerQuit

	CommandTypeMessage
)

type Command interface {
	Type() CommandType
	ToString() string
}

type CommandChan chan Command

// Client Connect Command
type ClientConnectCommand struct {
	FromId int
	Client *Client
}

func (c ClientConnectCommand) Type() CommandType {
	return CommandType(CommandTypeClientConnect)
}

func (c ClientConnectCommand) ToString() string {
	return fmt.Sprintf("ClientConnectCommand FromId: %d", c.FromId)
}

// Client Quit Command
type ClientQuitCommand struct {
	FromId int
	Client *Client
}

func (c ClientQuitCommand) Type() CommandType {
	return CommandType(CommandTypeClientQuit)
}

func (c ClientQuitCommand) ToString() string {
	return fmt.Sprintf("ClientQuitCommand FromId: %d", c.FromId)
}

// Server Connect Command
type ServerConnectCommand struct {
	FromId int
	Server *ServerClient
}

func (c ServerConnectCommand) Type() CommandType {
	return CommandType(CommandTypeServerConnect)
}

func (c ServerConnectCommand) ToString() string {
	return fmt.Sprintf("ServerConnectCommand FromId: %d", c.FromId)
}

// Server Quit Command
type ServerQuitCommand struct {
	FromId int
	Server *ServerClient
}

func (c ServerQuitCommand) Type() CommandType {
	return CommandType(CommandTypeServerQuit)
}

func (c ServerQuitCommand) ToString() string {
	return fmt.Sprintf("ServerQuitCommand FromId: %d", c.FromId)
}

// Mesage Command
type MessageCommand struct {
	FromId  int
	Message *Message
}

func (c MessageCommand) Type() CommandType {
	return CommandType(CommandTypeMessage)
}

func (c MessageCommand) ToString() string {
	return fmt.Sprintf("MessageCommand FromId: %d Message: %s", c.FromId, c.Message.ToString())
}

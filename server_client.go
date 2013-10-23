package multiplexity

import (
	"log"
	"net"
)

type ServerClient struct {
	conn        *Connection
	commandChan CommandChan
	IsConnected bool
}

const ServerClientId = -1

func NewServerClient(commandChan CommandChan) *ServerClient {
	server := &ServerClient{
		commandChan: commandChan,
		IsConnected: false,
	}
	return server
}

func (s *ServerClient) Connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalln("Failed to dial", address, err)
		return err
	}

	s.conn = NewConnection(conn, s.onReadMessage)

	return nil
}

func (s *ServerClient) ReadWritePumps() {
	s.IsConnected = true
	s.commandChan <- &ServerConnectCommand{
		Server: s,
	}

	s.conn.ReadWritePumps()

	s.commandChan <- &ServerQuitCommand{
		Server: s,
	}
}

func (s *ServerClient) Disconnect() {
	log.Println("Disconnecting ServerClient")
	s.conn.Disconnect()
}

func (s *ServerClient) Write(message *Message) {
	if s.conn != nil && s.conn.WriteChan != nil {
		s.conn.WriteChan <- message
	} else {
		log.Fatalln("Attempting to write message to server, but is not alive", message.ToString())
	}
}

func (s *ServerClient) onReadMessage(message *Message) {
	s.commandChan <- &ServerMessageCommand{
		Message: message,
	}
}

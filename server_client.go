package multiplexity

import (
	"log"
	"net"
)

type ServerClient struct {
	conn        *Connection
	commandChan CommandChan
	readChan    MessageChan
	IsConnected bool
}

const ServerClientId = -1

func NewServerClient(commandChan CommandChan) *ServerClient {
	return &ServerClient{
		commandChan: commandChan,
		readChan:    make(MessageChan),
		IsConnected: false,
	}
}

func (s *ServerClient) Connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalln("Failed to dial", address, err)
		return err
	}

	s.conn = NewConnection(conn, s.readChan)

	return nil
}

func (s *ServerClient) ReadWritePumps() {
	go s.readIntercept()

	s.IsConnected = true
	s.commandChan <- &ServerConnectCommand{
		FromId: ServerClientId,
		Server: s,
	}

	s.conn.ReadWritePumps()

	s.commandChan <- &ServerQuitCommand{
		FromId: ServerClientId,
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

func (s *ServerClient) readIntercept() {
	for {
		message, ok := <-s.readChan
		if !ok {
			log.Println("ServerClient read intercept closed")
			break
		}
		s.commandChan <- &MessageCommand{
			FromId:  ServerClientId,
			Message: message,
		}
	}
}

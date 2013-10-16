package multiplexity

import (
	"log"
	"net"
)

type ClientServer struct {
	commandChan CommandChan
}

func NewClientServer(commandChan CommandChan) *ClientServer {
	server := &ClientServer{
		commandChan: commandChan,
	}
	return server
}

func (s *ClientServer) Listen(address string) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalln("Failed to listen on", address)
		return
	}

	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalln("Recevied error while listening for connections", err)
			continue
		}

		go s.handleConnection(0, conn)
	}
}

func (s *ClientServer) Write(message *Message) {
	log.Println("ClientServer asked to post message", message.ToString())
}

func (s *ClientServer) handleConnection(newId int, conn net.Conn) {
	client := NewClient(newId, conn, s.commandChan)

	s.commandChan <- &ClientConnectCommand{
		FromId: newId,
		Client: client,
	}

	client.ReadWritePumps()

	s.commandChan <- &ClientQuitCommand{
		FromId: newId,
		Client: client,
	}
}

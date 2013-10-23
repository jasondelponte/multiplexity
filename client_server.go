package multiplexity

import (
	"log"
	"net"
)

type ClientServer struct {
	commandChan CommandChan
	clientIds   int
}

func NewClientServer(commandChan CommandChan) *ClientServer {
	server := &ClientServer{
		commandChan: commandChan,
		clientIds:   0,
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

		go s.handleConnection(s.clientIds, conn)
		s.clientIds += 1
	}
}

func (s *ClientServer) handleConnection(newId int, conn net.Conn) {
	client := NewClient(newId, conn, s.commandChan)

	s.commandChan <- &ClientConnectCommand{
		Client: client,
	}

	client.ReadWritePumps()

	s.commandChan <- &ClientQuitCommand{
		Client: client,
	}
}

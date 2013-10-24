package multiplexity

import (
	"fmt"
	"time"
)

type Proxy struct {
	config *Config

	serverClient *ServerClient
	clientServer *ClientServer

	clientCmdChan CommandChan
	serverCmdChan CommandChan

	clients ClientList
}

func NewProxy(config *Config) *Proxy {
	return &Proxy{
		config:        config,
		clientCmdChan: make(CommandChan),
		serverCmdChan: make(CommandChan),
		clients:       make(ClientList, 0, 0),
	}
}

func (p *Proxy) Start() {
	p.createClientServer()
	p.createServerClient()

	serverCmdHnd := NewServerCommandHandler(p.config, p.serverWriteMsg, p.clientWriteMsg)
	clientCmdHnd := NewClientCommandHandler(p.config, p.serverWriteMsg, p.clientWriteMsg)

	pingTick := time.Tick(time.Minute)
	for {
		select {
		case <-pingTick:
			pingMsg := &Message{
				Command:  "PING",
				Trailing: "something really cool just happened",
			}
			p.clientWriteMsg(pingMsg)

		case cmd, _ := <-p.serverCmdChan:
			serverCmdHnd.Handle(cmd)

		case cmd, _ := <-p.clientCmdChan:
			clientCmdHnd.Handle(cmd)

			// hack client list needs to be outside of proxy and cmd handelr
			p.clients = clientCmdHnd.clients
		}
	}
}

func (p *Proxy) serverWriteMsg(message *Message) {
	if p.serverClient != nil && p.serverClient.IsConnected {
		p.serverClient.Write(message)
	}
}
func (p *Proxy) clientWriteMsg(message *Message) {
	for _, client := range p.clients {
		outMsg := message.Copy()

		switch message.Command {
		case "PING": // do nothing
		case "JOIN": // sub client info
			outMsg.Prefix = fmt.Sprintf("%s!~%s@%s", client.Nick, client.Nick, client.RemoteHost)

		default:
			outMsg.Prefix = p.config.Client.Hostname
		}
		client.Write(outMsg)
	}
}

func (p *Proxy) createServerClient() error {
	p.serverClient = NewServerClient(p.serverCmdChan)
	p.serverClient.Connect(p.config.Server.Address)
	go p.serverClient.ReadWritePumps()

	return nil
}

func (p *Proxy) createClientServer() error {
	p.clientServer = NewClientServer(p.clientCmdChan)
	go p.clientServer.Listen(p.config.Client.ListenAddress)

	return nil
}

package multiplexity

import ()

type Proxy struct {
	config *Config

	serverClient *ServerClient
	clientServer *ClientServer

	clientCmdChan CommandChan
	serverCmdChan CommandChan
}

func NewProxy(config *Config) *Proxy {
	return &Proxy{
		config:        config,
		clientCmdChan: make(CommandChan),
		serverCmdChan: make(CommandChan),
	}
}

func (p *Proxy) Start() {
	p.createClientServer()
	p.createServerClient()

	serverCmdHnd := NewServerCommandHandler(p.config, p.serverWriteMsg, p.clientWriteMsg)
	clientCmdHnd := NewClientCommandHandler(p.config, p.serverWriteMsg, p.clientWriteMsg)

	for {
		select {
		case cmd, _ := <-p.serverCmdChan:
			serverCmdHnd.Handle(cmd)

		case cmd, _ := <-p.clientCmdChan:
			clientCmdHnd.Handle(cmd)
		}
	}
}

func (p *Proxy) serverWriteMsg(message *Message) {
	if p.serverClient != nil && p.serverClient.IsConnected {
		p.serverClient.Write(message)
	}
}
func (p *Proxy) clientWriteMsg(message *Message) {
	if p.clientServer != nil {
		p.clientServer.Write(message)
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

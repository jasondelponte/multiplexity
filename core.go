package multiplexity

import ()

type Core struct {
	config *Config

	serverClient *ServerClient
	clientServer *ClientServer

	clientCmdChan CommandChan
	clientMsgChan MessageChan

	serverCmdChan CommandChan
	serverMsgChan MessageChan
}

func NewCore(config *Config) *Core {
	return &Core{
		config:        config,
		clientCmdChan: make(CommandChan),
		clientMsgChan: make(MessageChan),
		serverCmdChan: make(CommandChan),
		serverMsgChan: make(MessageChan),
	}
}

func (c *Core) StartProxy() {
	c.createClientServer()
	c.createServerClient()

	serverCmdHnd := NewServerCommandHandler(c.config, c.serverCmdChan, c.serverMsgChan, c.clientMsgChan)
	go serverCmdHnd.Run()

	for {
		select {
		case msg, _ := <-c.clientMsgChan:
			if c.clientServer != nil {
				c.clientServer.Write(msg)
			}
		case msg, _ := <-c.serverMsgChan:
			if c.serverClient != nil && c.serverClient.IsConnected {
				c.serverClient.Write(msg)
			}
		}
	}
}

func (c *Core) createServerClient() error {
	c.serverClient = NewServerClient(c.serverCmdChan)
	c.serverClient.Connect(c.config.Server.Address)
	go c.serverClient.ReadWritePumps()

	return nil
}

func (c *Core) createClientServer() error {
	c.clientServer = NewClientServer(c.clientCmdChan)
	go c.clientServer.Listen(c.config.Client.ListenAddress)

	return nil
}

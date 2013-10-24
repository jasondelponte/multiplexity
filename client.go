package multiplexity

import (
	"log"
	"net"
	"strings"
)

type Client struct {
	Id          int
	Nick        string
	RemoteAddr  string
	RemoteHost  string
	conn        *Connection
	commandChan CommandChan
}

type ClientList []*Client

func NewClient(id int, conn net.Conn, commandChan CommandChan) *Client {
	client := &Client{
		Id:          id,
		Nick:        "unknown",
		commandChan: commandChan,
	}
	client.reverseLookupClient(conn)
	client.conn = NewConnection(conn, client.onReadMessage)
	return client
}

func (c *Client) ReadWritePumps() {
	c.conn.ReadWritePumps()
}

func (c *Client) Disconnect() {
	log.Println("Disconnecting client")
	c.conn.Disconnect()
}

func (c *Client) Write(message *Message) {
	if c.conn != nil && c.conn.WriteChan != nil {
		c.conn.WriteChan <- message
	} else {
		log.Fatalln("Attempting to write message to client, but is not alive", message)
	}
}

func (c *Client) onReadMessage(message *Message) {
	c.commandChan <- &ClientMessageCommand{
		Message: message,
		Client:  c,
	}
}

func (c *Client) reverseLookupClient(conn net.Conn) {
	c.RemoteAddr = strings.Split(conn.RemoteAddr().String(), ":")[0]

	host, err := net.LookupAddr(c.RemoteAddr)
	if err != nil {
		log.Println("Failed to reverse lookup hostname for", c.RemoteAddr, c.Id)
	}

	if len(host) > 0 {
		c.RemoteHost = host[0]
	} else {
		c.RemoteHost = c.RemoteAddr
	}
}

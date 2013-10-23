package multiplexity

import (
	"log"
	"net"
)

type Client struct {
	Id          int
	conn        *Connection
	commandChan CommandChan
}

type ClientList []*Client

func NewClient(id int, conn net.Conn, commandChan CommandChan) *Client {
	client := &Client{
		Id:          id,
		commandChan: commandChan,
	}
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
		log.Fatalln("Attempting to write message to client, but is not alive", message.ToString())
	}
}

func (c *Client) onReadMessage(message *Message) {
	c.commandChan <- &ClientMessageCommand{
		Message: message,
		Client:  c,
	}
}

package multiplexity

import (
	"log"
	"net"
)

type Client struct {
	Id          int
	conn        *Connection
	readChan    MessageChan
	commandChan CommandChan
}

func NewClient(id int, conn net.Conn, commandChan CommandChan) *Client {
	readChan := make(MessageChan)
	return &Client{
		Id:          id,
		conn:        NewConnection(conn, readChan),
		readChan:    readChan,
		commandChan: commandChan,
	}
}

func (c *Client) ReadWritePumps() {
	go c.readIntercept()
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

func (c *Client) readIntercept() {
	for {
		message, ok := <-c.readChan
		if !ok {
			log.Println("Client read intercept closed")
			break
		}
		c.commandChan <- &MessageCommand{
			FromId:  c.Id,
			Message: message,
		}
	}
}

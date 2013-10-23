package multiplexity

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Connection struct {
	conn       net.Conn
	closeWrite chan bool
	onReadFn   MessageHandlerFn
	WriteChan  MessageChan
}

func NewConnection(conn net.Conn, onReadFn MessageHandlerFn) *Connection {
	return &Connection{
		conn:       conn,
		closeWrite: make(chan bool),
		onReadFn:   onReadFn,
		WriteChan:  make(MessageChan),
	}
}

func (c *Connection) ReadWritePumps() {
	defer func() {
		log.Println("Closing Connection write channel and connection")
		close(c.WriteChan)
		c.conn.Close()
	}()

	go c.reader()
	c.writer()
}

func (c *Connection) Disconnect() {
	c.closeWrite <- true
}

func (c *Connection) reader() {
	reader := bufio.NewReader(c.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Connection read error", err, "closing reader")
			c.closeWrite <- true
			break
		}

		c.onReadFn(ParseMessage(message))
	}
	log.Println("Connection reader closing")
}

func (c *Connection) writer() {
	for {
		select {
		case <-c.closeWrite:
			log.Println("Connection closing writer")
			return

		case message, ok := <-c.WriteChan:
			if !ok {
				log.Println("Connection write channel closed")
				return
			}

			wrote, err := fmt.Fprint(c.conn, message.Serialize())
			if err != nil {
				log.Println("Connection write failed,", err)
				return
			}

			log.Println("Wrote to connection", message.ToString(), wrote)
		}
	}
	log.Println("Connection writer closing")
}

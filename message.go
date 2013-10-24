package multiplexity

import (
	"regexp"
	"strings"
)

type Message struct {
	Prefix   string
	Command  string
	Params   []string
	Trailing string
}

var (
	messageLineEndReg = regexp.MustCompile("((\r)?\n)+")
)

type MessageChan chan *Message
type MessageHandlerFn func(*Message)

func (m Message) Copy() *Message {
	msg := &Message{
		Prefix:   m.Prefix,
		Command:  m.Command,
		Params:   make([]string, 0, len(m.Params)),
		Trailing: m.Trailing,
	}
	msg.Params = append(msg.Params, m.Params...)

	return msg
}

func ParseMessage(message string) *Message {
	message = messageLineEndReg.ReplaceAllString(message, "")
	prefixEnd := 0
	trailingStart := len(message)

	msg := &Message{}

	// Prefix
	colon := strings.Index(message, ":")
	if colon == 0 {
		prefixEnd = strings.Index(message, " ")
		msg.Prefix = message[1:prefixEnd]
		prefixEnd = prefixEnd + 1
	}

	// Trailing, Defined as first spac e+colon to occur in the message
	colon = strings.Index(message[prefixEnd:], " :")
	if colon != -1 {
		trailingStart = prefixEnd + colon
		msg.Trailing = message[trailingStart+2:]
	}

	// Command and params of th emessage
	body := strings.Split(message[prefixEnd:trailingStart], " ")

	if len(body) == 0 {
		//TODO need error reporting
		return nil
	}

	// Body
	msg.Command = body[0]

	// Params
	for _, param := range body[1:] {
		if len(param) == 0 {
			continue
		}
		msg.Params = append(msg.Params, param)
	}

	return msg
}

func (m Message) Serialize() string {
	return m.ToString() + "\r\n"
}

func (m Message) ToString() string {
	msg := ""

	if len(m.Prefix) > 0 {
		msg += ":" + m.Prefix + " "
	}

	msg += m.Command

	if len(m.Params) > 0 {
		msg += " " + strings.Join(m.Params, " ")
	}

	if len(m.Trailing) > 0 {
		msg += " :" + m.Trailing
	}

	return msg
}

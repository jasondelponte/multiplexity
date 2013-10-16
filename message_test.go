package multiplexity

import (
	"testing"
)

func TestParseMessage(t *testing.T) {
	tstMsg := ":CalebDelnay!calebd@localhost PRIVMSG #mychannel :Hello everyone!"

	p := ParseMessage(tstMsg)

	if p == nil {
		t.Fatal("Failed to parse message")
	}

	if p.Prefix != "CalebDelnay!calebd@localhost" {
		t.Error("Failed to parse message prefix,", p.Prefix)
	}
	if p.Command != "PRIVMSG" {
		t.Error("Command is not PRIVMSG,", p.Command)
	}
	if len(p.Params) != 1 || p.Params[0] != "#mychannel" {
		t.Error("Params is not [#mychannel],", p.Params, len(p.Params))
	}
	if p.Trailing != "Hello everyone!" {
		t.Error("Trailing not equal Hello everyone!,", p.Trailing)
	}
}

func TestParseMessageNoPrefix(t *testing.T) {
	tstMsg := "PRIVMSG #mychannel :Hello everyone!"

	p := ParseMessage(tstMsg)

	if p == nil {
		t.Fatal("Failed to parse message")
	}

	if len(p.Prefix) != 0 {
		t.Error("Preix is not empty", p.Prefix)
	}
	if p.Command != "PRIVMSG" {
		t.Error("Command is not PRIVMSG,", p.Command)
	}
	if len(p.Params) != 1 || p.Params[0] != "#mychannel" {
		t.Error("Params is not [#mychannel],", p.Params, len(p.Params))
	}
	if p.Trailing != "Hello everyone!" {
		t.Error("Trailing not equal Hello everyone!,", p.Trailing)
	}
}

func TestParseMessageNoParams(t *testing.T) {
	tstMsg := ":CalebDelnay!calebd@localhost PRIVMSG :Hello everyone!"

	p := ParseMessage(tstMsg)

	if p == nil {
		t.Fatal("Failed to parse message")
	}

	if p.Prefix != "CalebDelnay!calebd@localhost" {
		t.Error("Preix is not CalebDelnay!calebd@localhost", p.Prefix)
	}
	if p.Command != "PRIVMSG" {
		t.Error("Command is not PRIVMSG,", p.Command)
	}
	if len(p.Params) != 0 {
		t.Error("Params is not empty,", p.Params, len(p.Params))
	}
	if p.Trailing != "Hello everyone!" {
		t.Error("Trailing not equal Hello everyone!,", p.Trailing)
	}
}

func TestParseMessageNoTrail(t *testing.T) {
	tstMsg := "PRIVMSG #mychannel"

	p := ParseMessage(tstMsg)

	if p == nil {
		t.Fatal("Failed to parse message")
	}

	if len(p.Prefix) != 0 {
		t.Error("Preix is not empty", p.Prefix)
	}
	if p.Command != "PRIVMSG" {
		t.Error("Command is not PRIVMSG,", p.Command)
	}
	if len(p.Params) != 1 || p.Params[0] != "#mychannel" {
		t.Error("Params is not [#mychannel],", p.Params, len(p.Params))
	}
	if len(p.Trailing) != 0 {
		t.Error("Trailing is not empty,", p.Trailing)
	}
}

func TestParseMessageWithPrefixNoTrail(t *testing.T) {
	tstMsg := ":CalebDelnay!calebd@localhost PRIVMSG #mychannel"

	p := ParseMessage(tstMsg)

	if p == nil {
		t.Fatal("Failed to parse message")
	}

	if p.Prefix != "CalebDelnay!calebd@localhost" {
		t.Error("Preix is not CalebDelnay!calebd@localhost", p.Prefix)
	}
	if p.Command != "PRIVMSG" {
		t.Error("Command is not PRIVMSG,", p.Command)
	}
	if len(p.Params) != 1 || p.Params[0] != "#mychannel" {
		t.Error("Params is not [#mychannel],", p.Params, len(p.Params))
	}
	if len(p.Trailing) != 0 {
		t.Error("Trailing is not empty,", p.Trailing)
	}
}

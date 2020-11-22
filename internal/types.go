package internal

import (
	"fmt"
	"net"
)

// Conn represents a net.Connection with some extra meta data like an ID
type Conn interface {
	net.Conn
	ID() interface{}
}

// Result is an interface that holds the value or error; preferably returned
// from a called CommandHandler
type Result interface {
	error
	fmt.Stringer
	Value() interface{}
	Ok() bool
}

// Command is an interface that represents a parsed actionable request
type Command interface {
	Input() string
	Action() string
	Body() string
}

// CommandHandler is an interface to represent a handler to act as a callback
type CommandHandler interface {
	Call(Conn, Command) Result
}

// CommandRouter is an interface that connects an end point to a CommandHandler
type CommandRouter interface {
	Handle(string, CommandHandler)
	HandleFunc(string, func(Conn, Command) Result)
	Parse(string) Result
	ParseAndCall(Conn, string) Result
}

package handler

import (
	"github.com/the-forges/example-net/internal"
)

// QuitHandler closes the connection, effectively quitting from the connections
// point of view
func QuitHandler(conn internal.Conn, cmd internal.Command) internal.Result {
	if err := conn.Close(); err != nil {
		return result{cmd, err}
	}
	return result{cmd, nil}
}

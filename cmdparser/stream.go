package cmdparser

import "net"

type Stream struct {
	net.Conn
	id interface{}
}

func (s Stream) ID() interface{} {
	return s.id
}

func NewStream(conn net.Conn, id interface{}) Stream {
	return Stream{conn, id}
}

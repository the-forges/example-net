package model

import "net"

// User struct to hold user information
type User struct {
	ID       int64
	Conn     net.Conn
	UserName string
}

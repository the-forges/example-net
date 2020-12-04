package util

import (
	"fmt"
	"net"
)

func WriteMessage(c net.Conn, msg string) {
	fmt.Fprintf(c, "%s\n", msg)
}

package handler

import (
	"context"
	"fmt"
	"net"
	"strings"
	"the-forges/example-net/model"
)

func writeMessage(c net.Conn, msg string) {
	fmt.Fprintf(c, "%s\n", msg)
}

func BroadcastMessageHandler(ctx context.Context, conn net.Conn, args ...string) (bool, error){
	if len(args) == 0 {
		return false, fmt.Errorf("empty message")
	}

	users, err := usersFromContext(ctx)
	if err != nil {
		return false, err
	}

	for _, user := range  users{
		writeMessage(user.Conn, strings.Join(args, " "))
	}

	return true, nil
}

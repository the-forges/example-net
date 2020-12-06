package handler

import (
	"context"
	"fmt"
	"net"
	"strings"
	"the-forges/example-net/util"
)

func BroadcastMessageHandler(ctx context.Context, conn net.Conn, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("empty message")
	}

	users, err := usersFromContext(ctx)
	if err != nil {
		return err
	}

	for _, user := range users.Iter() {
		if user == nil {
			continue
		}
		util.WriteMessage(user.Conn, strings.Join(args, " "))
	}

	return nil
}

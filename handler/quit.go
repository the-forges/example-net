package handler

import (
	"context"
	"fmt"
	"net"
	"the-forges/example-net/util"
)

func QuitHandler(ctx context.Context, conn net.Conn, args ...string) error {
	users, err := usersFromContext(ctx)
	if err != nil {
		return err
	}
	id, ok := ctx.Value(util.CtxID).(int64)
	if !ok || id == 0 {
		return fmt.Errorf("missing connection ID in context")
	}
	if _, err = users.FindUser(id); err != nil {
		return err
	}
	connected, ok := ctx.Value(util.CtxConnected).(*bool)
	if !ok {
		return fmt.Errorf("connection error")
	}
	if err := conn.Close(); err != nil {
		return err
	}
	// If there is an error, I don't want to write over the current UsersMap
	// with the returned UsersMap.
	newusers, err := users.KickUser(id)
	if err != nil {
		return err
	}
	// Now it's safe to swap the UsersMaps
	users = newusers
	*connected = false
	// broadcastMessage(fmt.Sprintf("Connection: %v, ended.", id))
	return nil
}

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
	id, ok := ctx.Value(util.CtxID).(int)
	if !ok || id == 0 {
		return fmt.Errorf("missing connection ID in context")
	}
	user, ok := (*users)[id]
	if !ok || user == nil {
		return fmt.Errorf("cannot find user")
	}
	connected, ok := ctx.Value(util.CtxConnected).(*bool)
	if !ok {
		return fmt.Errorf("connection error")
	}
	if err := conn.Close(); err != nil {
		return err
	}
	delete(*users, id)
	*connected = false
	// broadcastMessage(fmt.Sprintf("Connection: %v, ended.", id))

	return nil
}

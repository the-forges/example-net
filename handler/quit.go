package handler

import (
	"context"
	"fmt"
	"net"
)

func QuitHandler(ctx context.Context, conn net.Conn, args ...string) (bool, error) {
	users, err := usersFromContext(ctx)
	if err != nil {
		return false, err
	}
	id, ok := ctx.Value("ID").(int)
	if !ok || id == 0 {
		return false, fmt.Errorf("missing connection ID in context")
	}
	user, ok := (*users)[id]
	if !ok {
		return false, fmt.Errorf("cannot find user")
	}

	if err := conn.Close(); err != nil {
		return false, err
	}

	
	delete(users, id)
	// broadcastMessage(fmt.Sprintf("Connection: %v, ended.", id))

	return true, nil
}


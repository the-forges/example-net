package Handler

import (
	"context"
	"fmt"
	"net"
	"strings"
	"the-forges/example-net/model"
	"the-forges/example-net/util"
)

func UpdateUserNameHandler(ctx context.Context, conn net.Conn, args ...string) error {
	if len(args) <= 0 {
		return fmt.Errorf("missing name argument")
	}
	users, err := usersFromContext(ctx)
	if err != nil {
		return err
	}
	id, ok := ctx.Value(util.CtxID).(int)
	if !ok {
		return fmt.Errorf("cannot find ID")
	}

	user, ok := (*users)[id]
	if !ok {
		return fmt.Errorf("cannot find user")
	}

	user.UserName = strings.Join(args, " ")
	return nil
}

func usersFromContext(ctx context.Context) (*map[int]*model.User, error) {
	users, ok := ctx.Value(util.CtxUsers).(*map[int]*model.User)
	if !ok || users == nil {
		return nil, fmt.Errorf("cannot find users")
	}

	return users, nil
}

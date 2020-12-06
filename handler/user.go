package handler

import (
	"context"
	"fmt"
	"log"
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
	id, ok := ctx.Value(util.CtxID).(int64)
	if !ok {
		return fmt.Errorf("cannot find ID")
	}

	user, err := users.FindUser(id)
	if err != nil {
		return err
	}

	user.UserName = strings.Join(args, " ")
	return nil
}

func usersFromContext(ctx context.Context) (model.UsersMap, error) {
	users, ok := ctx.Value(util.CtxUsers).(model.UsersMap)
	if !ok || users == nil {
		log.Println(users)
		return nil, fmt.Errorf("cannot find users")
	}

	return users, nil
}

package Handler

import (
	"context"
	"fmt"
	"net"
	"strings"
	"the-forges/example-net/model"
)

func UpdateUserNameHandler(ctx context.Context, conn net.Conn, args ...string) (bool, error){
	if len(args) <= 0 {
		return false, fmt.Errorf("missing name argument")
	}
	users, err := usersFromContext(ctx)
	if err != nil {
		return false, err
	}
	id, ok := ctx.Value("ID").(int)
	if !ok {
		return false, fmt.Errorf("cannot find ID")
	}

	user, ok := (*users)[id]
	if !ok {
		return false, fmt.Errorf("cannot find user")
	}

	user.UserName = strings.Join(args, " ")
	return true, nil
}

func usersFromContext(ctx context.Context) (*map[int]*model.User, error){
	users, ok := ctx.Value("Users").(*map[int]*model.User)
	if !ok || users == nil{
		return nil, fmt.Errorf("cannot find users")
	}

	return users, nil
}
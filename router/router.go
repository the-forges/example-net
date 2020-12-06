package router

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"the-forges/example-net/model"
	"the-forges/example-net/util"
	"time"
)

var (
	routes       = make(map[string]Handler, 0)
	defaultRoute Handler
	users        = model.NewUsersMap()
)

type Handler func(context.Context, net.Conn, ...string) error

func HandleFunc(command string, h Handler) {
	command = strings.TrimSpace(command)
	routes[command] = h
}

func DefaultFunc(h Handler) {
	defaultRoute = h
}

// Parse - takes in command checks map of cmds to see if there is one present
func Parse(command string) (Handler, []string, error) {
	command = strings.TrimSpace(command)
	args := make([]string, 0)
	if strings.HasPrefix(command, "/") {
		parts := strings.SplitN(command, " ", 2)
		args = append(args, parts[1:]...)
	} else {
		args = append(args, command)
	}

	h, ok := routes[command]
	if !ok {
		if defaultRoute != nil {
			return defaultRoute, args, nil
		}
		return nil, args, fmt.Errorf("cannot find command")
	}
	return h, args, nil
}

func Listen(server net.Listener) error {
	log.Println("Server running")
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println(err)
			return err
		}

		go connectionHandler(conn)
	}
}

// connection handler
func connectionHandler(conn net.Conn) {
	id := time.Now().Unix()
	user := &model.User{ID: id, Conn: conn}
	users = users.RecieveUser(user)
	connected := true
	log.Printf("id: %d, user: %#v, users: %#v, connected: %v", id, user, users, connected)
	ctx := context.Background()
	ctx = context.WithValue(ctx, util.CtxID, id)
	ctx = context.WithValue(ctx, util.CtxUsers, users)
	ctx = context.WithValue(ctx, util.CtxConnected, &connected)
	log.Printf("ctx: %#v", ctx)
	for {
		util.WriteMessage(conn, "Enter command")
		buf := bufio.NewReader(conn)
		req, err := buf.ReadBytes(byte('\n'))
		if err != nil {
			if err.Error() != "EOF" {
				log.Printf("error: %s", err)
			}
			continue
		}
		body := string(req)
		h, args, err := Parse(body)
		if err != nil {
			util.WriteMessage(conn, err.Error())
			continue
		}
		if err := h(ctx, conn, args...); err != nil {
			util.WriteMessage(conn, err.Error())
		}
		if !connected {
			break
		}
	}
	log.Printf("%v disconnected. You have %d connections.", id, users.Len())
}

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
	routes = make(map[string]Handler, 0)
	users  = model.NewUsersMap()
)

type Handler func(context.Context, net.Conn, ...string) error

func HandleFunc(command string, handler Handler) {
	command = strings.TrimSpace(command)
	routes[command] = handler
}

// Parse - takes in command checks map of cmds to see if there is one present
func Parse(command string) (Handler, error) {
	command = strings.TrimSpace(command)
	handler, ok := routes[command]
	if !ok {
		return nil, fmt.Errorf("cannot find command")
	}
	return handler, nil
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
			return
		}
		body := string(req)
		h, err := Parse(body)
		if err != nil {
			util.WriteMessage(conn, err.Error())
			return
		}
		if err := h(ctx, conn); err != nil {
			util.WriteMessage(conn, err.Error())
		}
		if !connected {
			break
		}
	}
	log.Printf("%v disconnected. You have %d connections.", id, users.Len())
}

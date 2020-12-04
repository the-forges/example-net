package router

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"the-forges/example-net/model"
	"time"
)

//TODO add command
//parse command
//execute enpoint handler if command is found
//return errors where necessary

var routes = make(map[string]Handler, 0)
var users = make(map[int64]model.User)


type Handler func(context.Context, net.Conn, ...string) (bool, error)

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
		
		// conns[connID] = conn
		// log.Printf("%v connected. You have %d connections.", connID, len(conns))


	}
}

// connection handler
func connectionHandler(conn net.Conn) {
	// TODO: Check from in may for pre exiting id (loop through map)
	connID := time.Now().Unix()
	user := model.User{ID: connID, Conn: conn}
	users[connID] = user
	for {
		writeMessage(conn, "Enter command")
		buf := bufio.NewReader(conn)
		req, err := buf.ReadBytes(byte('\n'))
		if err != nil {
			if err.Error() != "EOF" {
				log.Printf("error: %s", err)
			}
			return
		}
		body := string(req)
		if !commandParser(id, body, conn) {
			break
		}

	}
	log.Printf("%v disconnected. You have %d connections.", id, len(conns))
}
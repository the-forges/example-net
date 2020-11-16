package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

var conns = make(map[int64]net.Conn)
var users = make(map[int64]User)

// User struct to hold user information
type User struct {
	ID       int64
	Conn     net.Conn
	UserName string
}

func main() {
	server, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server running")
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		connID := time.Now().Unix()
		conns[connID] = conn
		log.Printf("%v connected. You have %d connections.", connID, len(conns))

		go userHandler(connID, conn)
		go connectionHandler(connID, conn)
	}
}

func writeMessage(c net.Conn, msg string) {
	fmt.Fprintf(c, "%s\n", msg)
}

func broadcastMessage(msg string) {
	if msg == "" {
		return
	}

	for _, conn := range conns {
		writeMessage(conn, msg)
	}
}

func userHandler(id int64, c net.Conn) {
	userbuf := bufio.NewReader(c)
	writeMessage(c, "Enter username")
	reqbuf, err := userbuf.ReadBytes('\n')
	if err != nil {
		if err.Error() != "EOF" {
			log.Printf("error: %s", err)
		}
		return
	}
	userName := string(reqbuf)
	userName = strings.TrimSpace(userName)
	user := User{ID: id, Conn: c, UserName: userName}
	users[id] = user
	writeMessage(c, fmt.Sprintf("Username: %s, ID: %d was created.", user.UserName, user.ID))
	
}

func connectionHandler(id int64, c net.Conn) {
	buf := bufio.NewReader(c)

	for {
		req, err := buf.ReadBytes(byte('\n'))
		if err != nil {
			if err.Error() != "EOF" {
				log.Printf("error: %s", err)
			}
			break
		}
		body := string(req)
		msg := fmt.Sprintf("%d wrote %s", id, body)
		broadcastMessage(msg)
		body = strings.TrimSpace(body)
		if body == "quit" {
			break
		}
	}
	c.Close()
	delete(conns, id)
	broadcastMessage(fmt.Sprintf("Connection: %v, ended.", id))
	log.Printf("%v disconnected. You have %d connections.", id, len(conns))
}

// commandParser - Parse commands and call a function don't broadcast commands but text
// Build user profile

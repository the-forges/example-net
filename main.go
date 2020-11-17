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

func createUser(id int64, c net.Conn) {
	userbuf := bufio.NewReader(c)
	writeMessage(c, "Enter username")
	reqbuf, err := userbuf.ReadBytes(byte('\n'))
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
	
	
}

func printUserInfo(id int64, c net.Conn) {
	writeMessage(c, fmt.Sprintf("Username: %s, ID: %d was created.", users[id].UserName, users[id].ID))
}

func printAllUserInfo(id int64, c net.Conn) {
	for _, user := range users {
		writeMessage(c, fmt.Sprintf("Username: %s, ID: %d was created.", user.UserName, user.ID))
	}
}

func userCommands(id int64, c net.Conn) {
	writeMessage(c, "Enter command")
	buf := bufio.NewReader(c)
	req, err := buf.ReadBytes(byte('\n'))
		if err != nil {
			if err.Error() != "EOF" {
				log.Printf("error: %s", err)
			}
			return
		}
		body := string(req)
		commandParser(id, body, c)
		// msg := fmt.Sprintf("%s wrote %s", users[id].UserName, body)
		// broadcastMessage(msg)
}

func connectionHandler(id int64, c net.Conn) {
	body := ""

	for {
		createUser(id, c)
		userCommands(id,c)
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
func commandParser(id int64, cmd string, c net.Conn) {
	switch cmd {
	case "":
		return
	case "/printUserInfo":
		printUserInfo(id,c)
	case "/printAllUserInfo":
		printAllUserInfo(id, c)
	case "/broadcastMessage":
		broadcastMessage(fmt.Sprintf("%s says Hello!", users[id].UserName))
	case "/personalMessage":
		
	case "/quit":
		return
	}
}
// Build user profile

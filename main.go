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
		go func(id int64, c net.Conn) {
			buf := bufio.NewReader(c)
			for {
				fmt.Fprintf(c, "> ")
				req, err := buf.ReadBytes(byte('\n'))
				if err != nil {
					if err.Error() != "EOF" {
						log.Printf("error: %s", err)
					}
					break
				}
				body := string(req)
				writeMessage(c, body)
				body = strings.TrimSpace(body)
				if body == "quit"{
					break
				}
			}
			c.Close()
			delete(conns, id)
			log.Printf("%v disconnected. You have %d connections.", id, len(conns))
		}(connID, conn)
	}
}

func writeMessage(c net.Conn, msg string) {
	fmt.Fprintf(c, "You wrote: %s\n", msg)
}
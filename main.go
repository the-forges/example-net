package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/the-forges/example-net/cmdparser"
	"github.com/the-forges/example-net/handler"
	"github.com/the-forges/example-net/internal"
)

var streams = make(map[int64]cmdparser.Stream)

func main() {
	router := cmdparser.NewRouter()
	router.HandleFunc("/quit", handler.QuitHandler)

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
		stream := cmdparser.NewStream(conn, time.Now().Unix())
		streams[stream.ID().(int64)] = stream
		log.Printf("%v connected. You have %d connections.", stream.ID(), len(streams))
		go func(s internal.Conn) {
			buf := bufio.NewReader(s)
			for {
				fmt.Fprintf(s, "> ")
				req, err := buf.ReadBytes(byte('\n'))
				if err != nil {
					if err.Error() != "EOF" {
						log.Printf("error: %s", err)
					}
					break
				}
				body := string(req)
				writeMessage(s, body)
				if r := router.ParseAndCall(internal.Conn(s), body); !r.Ok() {
					log.Printf("[%s/%d]\n", r.Error(), s.ID())
					writeMessage(
						s,
						fmt.Sprintf("There was an problem with your command: %s", r.Error()),
					)
				}
			}
			s.Close()
			delete(streams, s.ID().(int64))
			log.Printf("%v disconnected. You have %d connections.", s.ID(), len(streams))
		}(stream)
	}
}

func writeMessage(c net.Conn, msg string) {
	fmt.Fprintf(c, "%s\n", msg)
}

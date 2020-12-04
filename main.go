package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"the-forges/example-net/handler"
	"the-forges/example-net/router"
	"time"
)

func main() {
	router.HandleFunc("/quit", handler.QuitHandler)
	router.HandleFunc("", handler.BroadcastMessageHandler)
	server, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(router.Listen(server))
}

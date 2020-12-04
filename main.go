package main

import (
	"log"
	"net"
	"the-forges/example-net/handler"
	"the-forges/example-net/router"
)

func main() {
	router.HandleFunc("/name :name", handler.ChangeUserNameHandler)
	router.HandleFunc("/quit", handler.QuitHandler)
	router.HandleFunc("", handler.BroadcastMessageHandler)
	server, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(router.Listen(server))
}

package main

import (
	order "cbt_timelapses_backend/m/v2/internal/instances"
	"cbt_timelapses_backend/m/v2/internal/scripts"
	"cbt_timelapses_backend/m/v2/internal/ws"
	"context"
	"log"
	"time"
)

var ctx = context.Background()

func main() {

	server := ws.CreateServer(messageHandler)

	for {
		server.WriteMessageAll([]byte("Hello"))
		time.Sleep(1 * time.Second)
	}
}

func messageHandler(message []byte, server *ws.Server) {
	log.Println(string(message))
	newOrder, _ := order.FromJSON(message)
	go scripts.CreateFakeTimelapse(newOrder, server)
}

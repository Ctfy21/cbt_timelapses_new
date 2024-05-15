package main

import (
	order "cbt_timelapses_backend/m/v2/internal/instances"
	"cbt_timelapses_backend/m/v2/internal/scripts"
	"cbt_timelapses_backend/m/v2/internal/ws"
	"log"
	"time"
)

func main() {

	server := ws.CreateServer(messageHandler)

	for {
		server.WriteMessageAll([]byte("Hello"))
		time.Sleep(1 * time.Second)
	}
}

func messageHandler(message []byte, server *ws.Server) {
	log.Println(string(message))
	newOrder := order.CreateOrder("sb1", "centertable", "2024-05-10_00-00-00", "2024-05-11_00-00-00", order.Status_waiting)
	server.WriteMessageAll(newOrder.ToJSON())
	go scripts.CreateFakeTimelapse(newOrder, server)
}

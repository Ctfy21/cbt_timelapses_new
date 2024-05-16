package main

import (
	"cbt_timelapses_backend/m/v2/configs"
	"cbt_timelapses_backend/m/v2/internal/database"
	order "cbt_timelapses_backend/m/v2/internal/instances"
	"cbt_timelapses_backend/m/v2/internal/scripts"
	"cbt_timelapses_backend/m/v2/internal/ws"
	"log"
	"strconv"
)

var exit = make(chan bool)

func main() {

	ws.CreateServer(messageHandler)
	var _ = <-exit

}

func messageHandler(message []byte, server *ws.Server) {
	log.Println(string(message))
	newOrder, err := order.FromJSON(message)
	if err != nil || newOrder.OrderJSON.Status != configs.STATUS_WAITING {
		log.Println(err)
		return
	}
	newID := postOrderToDB(message, server)
	go scripts.CreateFakeTimelapse(newOrder, server, newID)
}

func postOrderToDB(message []byte, server *ws.Server) uint64 {
	newID := database.GetIncrId(server.RedisDB, "OrderID")
	database.SetJSONData(server.RedisDB, "Order:"+strconv.FormatUint(newID, 10), message)
	return newID
}

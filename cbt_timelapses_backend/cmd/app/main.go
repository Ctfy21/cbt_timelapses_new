package main

import (
	"cbt_timelapses_backend/m/v2/configs"
	"cbt_timelapses_backend/m/v2/internal/database"
	order "cbt_timelapses_backend/m/v2/internal/instances"
	"cbt_timelapses_backend/m/v2/internal/timelapse"
	"cbt_timelapses_backend/m/v2/internal/ws"
	"log"
	"strconv"
)

var exit = make(chan bool)
var ids = make(chan int, configs.MAX_TIMELAPSE_ORDER)
var orders = make(chan *order.OrderJSONType, configs.MAX_TIMELAPSE_ORDER)

func main() {
	server := ws.CreateServer(messageHandler)
	go timelapse.CreateQueue(orders, server, ids)
	var _ = <-exit

}

func messageHandler(message []byte, server *ws.Server) {
	log.Println(string(message))
	newOrder, err := order.FromJSON(message)
	if err != nil || newOrder.OrderJSON.Status != configs.STATUS_WAITING {
		log.Println(err)
		return
	}
	newID, err := postOrderToDB(newOrder, server)
	if err != nil {
		log.Println(err)
		return
	}
	ids <- newID
	orders <- newOrder
}

func postOrderToDB(order *order.OrderJSONType, server *ws.Server) (int, error) {
	newID := database.GetIncrId(server.DB, "OrderID")
	order.OrderJSON.Id = newID
	val, err := order.ToJSON()
	if err != nil {
		return 0, err
	}
	database.SetJSONData(server.DB, "Order:"+strconv.Itoa(newID), val)
	json := database.GetJSONArrayValuesFromKeyPattern(server.DB, "Order:*", false)
	server.WriteMessageAll(json)
	return newID, nil
}

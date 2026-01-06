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
	log.Println("=== Received message:", string(message))
	newOrder, err := order.FromJSON(message)
	if err != nil {
		log.Println("ERROR: Failed to parse JSON:", err)
		return
	}
	log.Printf("=== Parsed order: Room=%s, Camera=%s, Status=%d\n", 
		newOrder.OrderJSON.Room, newOrder.OrderJSON.Camera, newOrder.OrderJSON.Status)
	
	if newOrder.OrderJSON.Status != configs.STATUS_WAITING {
		log.Printf("ERROR: Invalid status %d, expected %d\n", newOrder.OrderJSON.Status, configs.STATUS_WAITING)
		return
	}
	
	log.Println("=== Saving order to DB...")
	newID, err := postOrderToDB(newOrder, server)
	if err != nil {
		log.Println("ERROR: Failed to save to DB:", err)
		return
	}
	log.Printf("=== Order saved with ID: %d\n", newID)
	
	ids <- newID
	orders <- newOrder
	log.Println("=== Order added to queue")
}

func postOrderToDB(order *order.OrderJSONType, server *ws.Server) (int, error) {
	log.Println("=== Getting new ID from DB...")
	newID := database.GetIncrId(server.DB, "OrderID")
	log.Printf("=== Got new ID: %d\n", newID)
	
	order.OrderJSON.Id = newID
	val, err := order.ToJSON()
	if err != nil {
		log.Println("ERROR: Failed to marshal order to JSON:", err)
		return 0, err
	}
	log.Printf("=== Order JSON: %s\n", string(val))
	
	log.Printf("=== Saving to DB with key: Order:%d\n", newID)
	database.SetJSONData(server.DB, "Order:"+strconv.Itoa(newID), val)
	
	log.Println("=== Getting all orders from DB...")
	json := database.GetJSONArrayValuesFromKeyPattern(server.DB, "Order:*", false)
	log.Printf("=== All orders JSON: %s\n", string(json))
	
	log.Println("=== Broadcasting to all clients...")
	server.WriteMessageAll(json)
	
	return newID, nil
}

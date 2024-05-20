package timelapse

import (
	"cbt_timelapses_backend/m/v2/configs"
	"cbt_timelapses_backend/m/v2/internal/database"
	order "cbt_timelapses_backend/m/v2/internal/instances"
	"cbt_timelapses_backend/m/v2/internal/ws"
	"log"
	"os/exec"
	"strconv"
)

func CreateQueue(newOrder chan *order.OrderJSONType, server *ws.Server, newID chan int) {
	for {
		id := <-newID
		order := <-newOrder
		CreateTimelapse(order, server, id)
	}
}

func CreateTimelapse(order *order.OrderJSONType, server *ws.Server, id int) {
	err := exec.Command(configs.DIRECTORY_FOLDER_SCRIPT,
		"--dir",
		configs.SCREENSHOTS_FOLDER+"/"+order.OrderJSON.Room+"/"+order.OrderJSON.Camera,
		"--start",
		order.OrderJSON.StartDate,
		"--end",
		order.OrderJSON.EndDate).Run()

	order.OrderJSON.Status = configs.STATUS_OK
	if err != nil {
		order.OrderJSON.Status = configs.STATUS_ERROR
		log.Println(err)
	}

	val, err := order.ToJSON()
	if err != nil {
		log.Println("Error during JSON Order marshalling: ", err)
		return
	}

	database.SetJSONData(server.RedisDB, "Order:"+strconv.Itoa(id), val)
	server.WriteMessageAll(val)
}

//func CreateFakeTimelapse(order *order.OrderJSONType, server *ws.Server, newID int) {
//
//	err := exec.Command("ping", "google.com").Run()
//	order.OrderJSON.Status = configs.STATUS_OK
//	if err != nil {
//		order.OrderJSON.Status = configs.STATUS_ERROR
//		log.Println(err)
//	}
//
//	val, err := order.ToJSON()
//	if err != nil {
//		log.Println("Error during JSON Order marshalling: ", err)
//		return
//	}
//
//	database.SetJSONData(server.RedisDB, "Order:"+strconv.Itoa(newID), val)
//	server.WriteMessageAll(val)
//
//}

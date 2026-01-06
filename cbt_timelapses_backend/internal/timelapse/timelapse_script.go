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
	scriptPath := configs.DIRECTORY_FOLDER_SCRIPT
	dirPath := configs.SCREENSHOTS_FOLDER + "/" + order.OrderJSON.Room + "/" + order.OrderJSON.Camera
	
	log.Printf("=== Creating timelapse ID:%d, Script:%s, Dir:%s, Start:%s, End:%s\n",
		id, scriptPath, dirPath, order.OrderJSON.StartDate, order.OrderJSON.EndDate)
	
	cmd := exec.Command(scriptPath,
		"--dir", dirPath,
		"--start", order.OrderJSON.StartDate,
		"--end", order.OrderJSON.EndDate)
	
	// Захватываем stdout и stderr
	output, err := cmd.CombinedOutput()
	
	log.Printf("=== Script output for ID:%d:\n%s\n", id, string(output))

	order.OrderJSON.Status = configs.STATUS_OK
	if err != nil {
		order.OrderJSON.Status = configs.STATUS_ERROR
		log.Printf("=== ERROR creating timelapse ID:%d: %v\n", id, err)
	} else {
		log.Printf("=== SUCCESS creating timelapse ID:%d\n", id)
	}

	val, err := order.ToJSON()
	if err != nil {
		log.Println("Error during JSON Order marshalling: ", err)
		return
	}

	database.SetJSONData(server.DB, "Order:"+strconv.Itoa(id), val)
	
	// Отправляем полный список заказов всем клиентам
	json := database.GetJSONArrayValuesFromKeyPattern(server.DB, "Order:*", false)
	server.WriteMessageAll(json)
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

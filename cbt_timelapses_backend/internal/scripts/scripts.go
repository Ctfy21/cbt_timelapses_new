package scripts

import (
	"cbt_timelapses_backend/m/v2/configs"
	"cbt_timelapses_backend/m/v2/internal/database"
	order "cbt_timelapses_backend/m/v2/internal/instances"
	"cbt_timelapses_backend/m/v2/internal/ws"
	"log"
	"os/exec"
	"strconv"
)

func CreateTimelapse(order *order.Order, server *ws.Server) {
	err := exec.Command(configs.DIRECTORY_FOLDER_SCRIPT,
		"--dir",
		configs.SCREENSHOTS_FOLDER+"/"+order.Room+"/"+order.Camera,
		"--start",
		order.StartDate+"_00-00-00",
		"--end",
		order.EndDate+"_00-00-00").Run()

	if err != nil {
		server.WriteMessageAll([]byte("Error"))
		log.Println(err)
	}

}

func CreateFakeTimelapse(order *order.OrderJSONType, server *ws.Server, newID int) {

	err := exec.Command("ping", "google.com").Run()
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

	database.SetJSONData(server.RedisDB, "Order:"+strconv.Itoa(newID), val)
	server.WriteMessageAll(val)

}

package scripts

import (
	"cbt_timelapses_backend/m/v2/configs"
	order "cbt_timelapses_backend/m/v2/internal/instances"
	"cbt_timelapses_backend/m/v2/internal/ws"
	"log"
	"os/exec"
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

	val, _ := order.ToJSON()

	server.WriteMessageAll(val)

}

func CreateFakeTimelapse(order *order.Order, server *ws.Server) {
	err := exec.Command("ping", "google.com").Run()
	order.Status = configs.STATUS_OK
	val, _ := order.ToJSON()
	if err != nil {
		log.Println(err)
		server.WriteMessageAll([]byte("Error"))
		return
	}

	server.WriteMessageAll(val)

}

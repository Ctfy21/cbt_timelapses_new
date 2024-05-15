package scripts

import (
	"cbt_timelapses_backend/m/v2/configs"
	order "cbt_timelapses_backend/m/v2/internal/instances"
	"cbt_timelapses_backend/m/v2/internal/ws"
	"fmt"
	"log"
	"os/exec"
	"time"
)

func CreateTimelapse(order *order.Order) {
	err := exec.Command(configs.DIRECTORY_FOLDER_SCRIPT,
		"--dir",
		configs.SCREENSHOTS_FOLDER+"/"+order.Room+"/"+order.Camera,
		"--start",
		order.StartDate+"_00-00-00",
		"--end",
		order.EndDate+"_00-00-00").Run()

	if err != nil {
		log.Fatal(err)
	}

}

func CreateFakeTimelapse(order *order.Order, server *ws.Server) {
	err := exec.Command("ls", "-l").Run()
	time.Sleep(5 * time.Second)
	if err != nil {
		fmt.Println(err)
		server.WriteMessageAll([]byte("Error"))
		return
	}
	server.WriteMessageAll(order.ToJSON())

}

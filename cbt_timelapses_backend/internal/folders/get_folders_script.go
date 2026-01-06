package folders

import (
	"cbt_timelapses_backend/m/v2/configs"
	"log"
	"os"
)

func GetScreenshotsFolders() map[string][]string {
	rooms := make(map[string][]string)
	entriesRoom, err := os.ReadDir(configs.SCREENSHOTS_FOLDER)
	if err != nil {
		log.Println("Error reading SCREENSHOTS_FOLDER:", err)
		return rooms
	}
	
	for _, eR := range entriesRoom {
		if eR.IsDir() {
			roomPath := configs.SCREENSHOTS_FOLDER + "/" + eR.Name()
			entriesCamera, err := os.ReadDir(roomPath)
			if err != nil {
				log.Println("Error reading room folder:", roomPath, err)
				continue
			}
			
			cameras := []string{}
			for _, eC := range entriesCamera {
				if eC.IsDir() {
					cameras = append(cameras, eC.Name())
				}
			}
			
			if len(cameras) > 0 {
				rooms[eR.Name()] = cameras
			}
		}
	}
	
	log.Printf("Found folders: %+v\n", rooms)
	return rooms
}

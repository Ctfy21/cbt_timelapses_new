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
		log.Fatal(err)
		return make(map[string][]string)
	}
	for _, eR := range entriesRoom {
		if eR.IsDir() {
			entriesCamera, err := os.ReadDir(configs.SCREENSHOTS_FOLDER + "/" + eR.Name())
			if err != nil {
				log.Fatal(err)
				return make(map[string][]string)
			}
			for _, eC := range entriesCamera {
				if eC.IsDir() {
					rooms[eR.Name()] = append(rooms[eR.Name()], eC.Name())
				}
			}
		}
	}
	return rooms
}

package database

import (
	"cbt_timelapses_backend/m/v2/configs"
	"cbt_timelapses_backend/m/v2/internal/folders"
	order "cbt_timelapses_backend/m/v2/internal/instances"
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

func StartClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}

func SetJSONData(rc *redis.Client, key string, data []byte) {
	err := rc.Set(context.Background(), key, data, configs.ORDER_TTL).Err()
	if err != nil {
		log.Println("Redis set JSON error:", err)
	}
}

func GetJSONData(rc *redis.Client, key string) []byte {
	val, err := rc.Get(context.Background(), key).Bytes()
	if err != nil {
		log.Println("Redis get JSON error:", err)
	}
	return val
}

func GetIncrId(rc *redis.Client, key string) int {
	rc.Incr(context.Background(), key)
	val, err := rc.Get(context.Background(), key).Int()
	if err != nil {
		log.Println("Redis error GetIncrId:", err)
	}
	log.Println("Redis get incremented ID val:", val)
	return val
}

func GetJSONArrayValuesFromKeyPattern(rc *redis.Client, key string, needFolders bool) []byte {

	var newOrdersJson order.OrdersJSONType
	if needFolders {
		jsonFolders := folders.GetScreenshotsFolders()
		newOrdersJson = order.OrdersJSONType{OrdersJSON: []string{}, Folders: jsonFolders}
	} else {
		newOrdersJson = order.OrdersJSONType{OrdersJSON: []string{}}
	}

	result, err := rc.Keys(context.Background(), key).Result()
	if err != nil {
		log.Println("Error during get orders operation: ", err)
	}

	for _, curKey := range result {
		newOrdersJson.OrdersJSON = append(newOrdersJson.OrdersJSON, string(GetJSONData(rc, curKey)))
	}

	j, err := json.Marshal(newOrdersJson)

	if err != nil {
		log.Println("Error during json operation in GetJSONArrayValuesFromKeyPattern: ", err)
		return []byte("")
	}

	return j
}

// func FlushDB(rdb *redis.Client) {
// 	err := rdb.FlushDB(context.Background()).Err()
// 	if err != nil {
// 		panic(err)
// 	}
// }

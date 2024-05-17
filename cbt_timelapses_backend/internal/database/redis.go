package database

import (
	"cbt_timelapses_backend/m/v2/configs"
	order "cbt_timelapses_backend/m/v2/internal/instances"
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log"
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

func GetIncrId(rc *redis.Client, key string) uint64 {
	val, _ := rc.Get(context.Background(), key).Uint64()
	rc.Incr(context.Background(), key)
	log.Println("Redis get incremented ID val:", val)
	return val
}

func GetJSONArrayValuesFromKeyPattern(rc *redis.Client, key string) []byte {

	newOrdersJson := order.OrdersJSONType{OrdersJSON: []string{}}

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

func FlushDB(rdb *redis.Client) {
	err := rdb.FlushDB(context.Background()).Err()
	if err != nil {
		panic(err)
	}
}

package order

import (
	"encoding/json"
	"fmt"
	"log"
)

const (
	Status_ok      int = 200
	Status_waiting int = 300
	Status_error   int = 400
)

type Order struct {
	Room      string `json:"Room"`
	Camera    string `json:"Camera"`
	StartDate string `json:"StartDate"`
	EndDate   string `json:"EndDate"`
	Status    int    `json:"Status"`
}

func CreateOrder(room, camera, startDate, endDate string, status int) *Order {
	order := Order{
		Room:      room,
		Camera:    camera,
		StartDate: startDate,
		EndDate:   endDate,
		Status:    status,
	}

	return &order
}

func (order *Order) ToJSON() []byte {
	j, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("Error occured during marshaling. Error: %s", err.Error())
	}
	fmt.Printf("JSON: %s\n", string(j))
	return j
}

func FromJSON(jsonValue []byte) Order {
	var newOrder Order
	err := json.Unmarshal(jsonValue, &newOrder)
	if err != nil {
		log.Fatalf("Error occured during unmarshaling. Error: %s", err.Error())
	}
	fmt.Printf("Order struct: %s\n", newOrder)
	return newOrder
}

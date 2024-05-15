package order

import (
	"encoding/json"
	"fmt"
	"log"
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

func (newOrder *Order) ToJSON() ([]byte, error) {
	j, err := json.Marshal(newOrder)
	if err != nil {
		log.Fatalf("Error occured during marshaling. Error: %s", err.Error())
	}
	fmt.Printf("JSON: %s\n", string(j))
	return j, err
}

func FromJSON(jsonValue []byte) (*Order, error) {
	var newOrder Order
	err := json.Unmarshal(jsonValue, &newOrder)
	if err != nil {
		log.Fatalf("Error occured during unmarshaling. Error: %s", err.Error())
	}
	fmt.Printf("Order struct: %s\n", newOrder)
	return &newOrder, err
}

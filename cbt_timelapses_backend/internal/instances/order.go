package order

import (
	"encoding/json"
	"fmt"
	"log"
)

type Order struct {
	Id        int    `json:"Id,omitempty"`
	Room      string `json:"Room"`
	Camera    string `json:"Camera"`
	StartDate string `json:"StartDate"`
	EndDate   string `json:"EndDate"`
	Status    int    `json:"Status"`
}

type OrdersJSONType struct {
	OrdersJSON []string            `json:"Orders"`
	Folders    map[string][]string `json:"Folders"`
}

type OrderJSONType struct {
	OrderJSON Order `json:"Order"`
}

func (newOrder *OrderJSONType) ToJSON() ([]byte, error) {
	j, err := json.Marshal(newOrder)
	if err != nil {
		log.Fatalf("Error occured during marshaling. Error: %s", err.Error())
	}
	fmt.Printf("JSON: %s\n", string(j))
	return j, err
}

func FromJSON(jsonValue []byte) (*OrderJSONType, error) {
	var newOrder OrderJSONType
	err := json.Unmarshal(jsonValue, &newOrder)
	if err != nil {
		log.Fatalf("Error occured during unmarshaling. Error: %s", err.Error())
	}
	fmt.Printf("Order struct: %+v\n", newOrder)
	return &newOrder, err
}

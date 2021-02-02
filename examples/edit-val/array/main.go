package main

import (
	"fmt"

	"github.com/tamboto2000/jsonextract/v3"
)

type Item struct {
	Qty  int64  `json:"qty"`
	Name string `json:"name"`
	ID   int32
}

func main() {
	raw := []byte(`["Hello world!", 1, -2, 0.3]`)
	jsons, err := jsonextract.FromBytes(raw)
	if err != nil {
		panic(err.Error())
	}

	// there is only one json object
	json := jsons[0]
	fmt.Println("original raw json:", string(json.Runes()))

	// add string
	json.AddItem("added string")

	// add int32
	json.AddItem(int32(32))

	// add object from map
	json.AddItem(map[string]interface{}{
		"name":  "Franklin Collin Tamboto",
		"email": "tamboto2000@gmail.com",
		"id":    1,
	})

	// add object from struct
	json.AddItem(&Item{
		Qty:  23,
		Name: "ROG Phone 3",
		ID:   50,
	})

	fmt.Println("edited raw json:", string(json.Runes()))
}

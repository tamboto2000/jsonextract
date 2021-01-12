package main

import (
	"fmt"

	"github.com/tamboto2000/jsonextract/v2"
)

func main() {
	jsons, err := jsonextract.FromFile("test.txt")
	if err != nil {
		panic(err.Error())
	}

	// there is only on JSON object
	json := jsons[0]
	// change id value
	json.KeyVal["id"].Val = 2
	// reparse JSON
	if err := json.ReParse(); err != nil {
		panic(err.Error())
	}

	// print result
	fmt.Println(string(json.Raw.Bytes()))
}

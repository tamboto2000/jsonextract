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
	for _, val := range json.Vals {
		if val.Kind == jsonextract.Int {
			val.Val = 69
		}

		if val.Kind == jsonextract.String {
			val.Val = "This string is edited"
		}

		if val.Kind == jsonextract.Object {
			// value with key "key" is string, see test.txt
			val.KeyVal["key"].Val = "This string is edited"
		}
	}

	if err := json.ReParse(); err != nil {
		panic(err.Error())
	}

	// print result
	fmt.Println(string(json.Raw.Bytes()))
}

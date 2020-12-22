package main

import (
	"fmt"

	"github.com/tamboto2000/jsonextract/v2"
)

func main() {
	jsons, err := jsonextract.FromFile("array.txt")
	if err != nil {
		panic(err.Error())
	}

	for _, json := range jsons {
		if json.Kind == jsonextract.Array {
			for _, val := range json.Vals {
				fmt.Println("arr val:", val.Val)
			}
		}
	}

	if err := jsonextract.Save(jsons); err != nil {
		panic(err.Error())
	}
}

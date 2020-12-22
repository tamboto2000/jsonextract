package main

import (
	"fmt"

	"github.com/tamboto2000/jsonextract/v2"
)

func main() {
	jsons, err := jsonextract.FromFile("numeric.txt")
	if err != nil {
		panic(err.Error())
	}

	for _, json := range jsons {
		if json.Kind == jsonextract.Int {
			fmt.Println("num val:", json.Val.(int))
		} else {
			fmt.Println("num val:", json.Val.(float64))
		}
	}

	if err := jsonextract.Save(jsons); err != nil {
		panic(err.Error())
	}
}

package main

import (
	"fmt"

	"github.com/tamboto2000/jsonextract/v2"
)

func main() {
	jsons, err := jsonextract.FromFile("null.txt")
	if err != nil {
		panic(err.Error())
	}

	for _, json := range jsons {
		fmt.Println("null val:", json.Val)
	}

	if err := jsonextract.Save(jsons); err != nil {
		panic(err.Error())
	}
}

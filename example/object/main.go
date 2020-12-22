package main

import (
	"fmt"

	"github.com/tamboto2000/jsonextract/v2"
)

func main() {
	jsons, err := jsonextract.FromFile("object.txt")
	if err != nil {
		panic(err.Error())
	}

	for _, json := range jsons {
		if json.Kind == jsonextract.Object {
			for key, val := range json.KeyVal {
				fmt.Println(key+":", val)
			}
		}
	}

	if err := jsonextract.Save(jsons); err != nil {
		panic(err.Error())
	}
}

package main

import (
	"fmt"

	"github.com/tamboto2000/jsonextract/v2"
)

func main() {
	raw := []byte("null null {null}")

	jsons, err := jsonextract.FromBytes(raw)
	if err != nil {
		panic(err.Error())
	}

	for _, json := range jsons {
		if json.Kind == jsonextract.Null {
			fmt.Println("raw:", string(json.RawRunes()))
		}
	}
}

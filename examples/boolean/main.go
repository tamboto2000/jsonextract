package main

import (
	"fmt"

	"github.com/tamboto2000/jsonextract/v2"
)

func main() {
	raw := []byte(`
	true
	false
	fals
	`)

	jsons, err := jsonextract.FromBytes(raw)
	if err != nil {
		panic(err.Error())
	}

	for _, json := range jsons {
		if json.Kind == jsonextract.Boolean {
			fmt.Println("raw:", string(json.RawRunes()))
			fmt.Println("val:", json.Boolean())
		}
	}
}

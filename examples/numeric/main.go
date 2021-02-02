package main

import (
	"fmt"

	"github.com/tamboto2000/jsonextract/v3"
)

func main() {
	raw := []byte(`
	123
	12.3
	12.3e+4
	12.3e-4
	12.3e4
	123e+4
	123e-4
	123e4
	0.0
	-123
	-12.3e-4
	-12.3e+4
	`)

	jsons, err := jsonextract.FromBytes(raw)
	if err != nil {
		panic(err.Error())
	}

	for _, json := range jsons {
		if json.Kind() == jsonextract.Integer {
			fmt.Println("int value")
			fmt.Println("\traw:", string(json.Runes()))
			fmt.Println("\tval:", json.Integer())
		}

		if json.Kind() == jsonextract.Float {
			fmt.Println("float value")
			fmt.Println("\traw:", string(json.Runes()))
			fmt.Println("\tval:", json.Float())
		}
	}
}

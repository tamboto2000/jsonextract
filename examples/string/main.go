package main

import (
	"fmt"

	"github.com/tamboto2000/jsonextract/v2"
)

func main() {
	raw := []byte(`
	"Hello World! \u9090
	blablabla"
	`)
	jsons, err := jsonextract.FromBytes(raw)
	if err != nil {
		panic(err.Error())
	}

	for _, json := range jsons {
		if json.Kind() == jsonextract.String {
			fmt.Println("raw:", string(json.Runes()))
			fmt.Println("val:", json.String())
		}
	}
}

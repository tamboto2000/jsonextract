package main

import (
	"fmt"

	"github.com/tamboto2000/jsonextract/v3"
)

func main() {
	raw := []byte(`[1, "2", -3, 0.4, -0.5, "\"6\"", null, ["inner array!", 123], true, false]`)
	jsons, err := jsonextract.FromBytes(raw)
	if err != nil {
		panic(err.Error())
	}

	for _, json := range jsons {
		if json.Kind() == jsonextract.Array {
			fmt.Println("raw:", string(json.Runes()))
			fmt.Println("array vals:")
			vals := json.Array()
			for _, val := range vals {
				fmt.Println("\t", string(val.Runes()))
			}
		}
	}
}

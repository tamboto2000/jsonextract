package main

import (
	"fmt"

	"github.com/tamboto2000/jsonextract/v2"
)

func main() {
	// Extract int JSONs, ignores zero int
	// jsons, err := jsonextract.FromFileWithOpt("test.txt", jsonextract.Option{ParseInt: true, IgnoreZeroInt: true})
	jsons, err := jsonextract.FromFileWithOpt("test.txt", jsonextract.Option{
		ParseInt:      true,
		IgnoreZeroInt: true,
	})
	if err != nil {
		panic(err.Error())
	}

	for _, json := range jsons {
		fmt.Println(json.Val)
	}
}

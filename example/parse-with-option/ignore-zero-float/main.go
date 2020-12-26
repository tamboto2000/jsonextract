package main

import (
	"fmt"

	"github.com/tamboto2000/jsonextract/v2"
)

func main() {
	// Extract float JSONs, ignores zero float
	jsons, err := jsonextract.FromFileWithOpt("test.txt", jsonextract.Option{
		ParseFloat:      true,
		IgnoreZeroFloat: true,
	})
	if err != nil {
		panic(err.Error())
	}

	for _, json := range jsons {
		fmt.Println(json.Val.(float64))
	}
}

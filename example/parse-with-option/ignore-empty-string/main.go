package main

import (
	"fmt"

	"github.com/tamboto2000/jsonextract/v2"
)

func main() {
	jsons, err := jsonextract.FromFileWithOpt("test.txt", jsonextract.Option{
		ParseStr:       true,
		IgnoreEmptyStr: true,
	})

	if err != nil {
		panic(err.Error())
	}

	for _, json := range jsons {
		fmt.Println(json.Val)
	}
}

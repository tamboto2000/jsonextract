package main

import (
	"fmt"

	"github.com/tamboto2000/jsonextract/v2"
)

func main() {
	raw := []byte(`
	{
		"1": 1,
		"2": -2,
		"3": 0.3,
		"4": -0.4,
		"5": true,
		"6": false,
		"7": null,
		"8": [1, -2, 0.3, -0.4, true, false, null, {"a":"b"}],
		"\u9090": "chinese"
	}
	
	{
		"innerObj": {"msg":"hello world"}
	}
	`)

	jsons, err := jsonextract.FromBytes(raw)
	if err != nil {
		panic(err.Error())
	}

	for _, json := range jsons {
		// if json.Kind == jsonextract.Object {
		// 	fmt.Println("raw:", string(json.RawRunes()))
		// }

		fmt.Println("raw:", string(json.RawRunes()))
	}
}

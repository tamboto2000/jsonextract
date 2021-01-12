package main

import (
	"bytes"

	"github.com/tamboto2000/jsonextract/v2"
)

func main() {
	byts := []byte(`{
		"1": 1,
		"2": -2,
		"3": 0.3,
		"4": -0.4,		
		"5": true,
		"6": false,
		"7": null,
		"8": [1, -2, 0.3, -0.4, true, false, null, {"a":"b"}, 0.0e-1, 1e+2]
	}
	
	{}`)

	r := bytes.NewReader(byts)
	jsons, err := jsonextract.FromReader(r)
	if err != nil {
		panic(err.Error())
	}

	// save result to path
	if err := jsonextract.SaveToPath(jsons, "from_reader.json"); err != nil {
		panic(err.Error())
	}
}

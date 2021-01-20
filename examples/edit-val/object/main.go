package main

import (
	"fmt"

	"github.com/tamboto2000/jsonextract/v2"
)

func main() {
	raw := []byte(`{
		"key":"val",
		"intVal": 2
		}`)
	jsons, err := jsonextract.FromBytes(raw)
	if err != nil {
		panic(err.Error())
	}

	// only a json object is detected
	json := jsons[0]
	fmt.Println("original raw:", string(json.RawRunes()))
	keyVal := json.Object()

	// when you edit a value, it automatically generate new raw json
	// edit field "key" value.
	if err := keyVal["key"].EditStr("Hello World"); err != nil {
		panic(err.Error())
	}

	// edit field intVal
	if err := keyVal["intVal"].EditInt(69); err != nil {
		panic(err.Error())
	}

	// print raw JSON
	fmt.Println("edited raw:", string(json.RawRunes()))

	// print edited value
	str, _ := keyVal["key"].String()
	fmt.Println(`field "key" val:`, str)
	fmt.Println(`field "intVal" val:`, keyVal["intVal"].Integer())
}

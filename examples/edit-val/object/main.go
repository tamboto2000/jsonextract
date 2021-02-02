package main

import (
	"fmt"

	"github.com/tamboto2000/jsonextract/v3"
)

type User struct {
	Name  string `json:"name"`
	Email string
}

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
	fmt.Println("original raw:", string(json.Runes()))
	keyVal := json.Object()

	// adding a new value automatically generate new raw json

	// add new field
	json.AddField("newStr", "Added")

	// add new field from pointer
	strPtr := "ptr str"
	json.AddField("*string", &strPtr)

	// add nil value
	json.AddField("null", nil)

	// add int32
	json.AddField("int32", int32(32))

	// add float32
	json.AddField("float32", float32(0.69))

	// add bool
	json.AddField("bool", true)

	// add array
	json.AddField("array", []interface{}{"Hello", "World!", 1, -2, 0.3, nil})

	// add object with map
	json.AddField("mapObj", map[interface{}]interface{}{
		"string": "Hello world",
		"int32":  int32(32),
		"null":   nil,
		123:      456,
	})

	// add value with integer key
	json.AddField(123, "integer key")

	// add object with struct
	json.AddField("structObj", User{
		Name:  "Franklin Collin Tamboto",
		Email: "tamboto2000@gmail.com",
	})

	// when you edit a value, it automatically generate new raw json

	// edit field "key" value.
	keyVal["key"].SetStr("Hello World")

	// edit field intVal
	keyVal["intVal"].SetInt(69)

	// delete a field
	if json.DeleteField(123) {
		fmt.Println("deleted field", 123)
	}

	// print raw JSON
	fmt.Println("edited raw:", string(json.Runes()))
	// print items len
	fmt.Println("len:", json.Len())
}

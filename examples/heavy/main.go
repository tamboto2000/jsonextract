package main

import "github.com/tamboto2000/jsonextract/v2"

func main() {
	jsons, err := jsonextract.FromFile("test.txt")

	if err != nil {
		panic(err.Error())
	}

	newList := make([]*jsonextract.JSON, 0)
	for _, json := range jsons {
		if json.Kind() == jsonextract.Array {
			newList = append(newList, json)
		}

		if json.Kind() == jsonextract.Object {
			newList = append(newList, json)
		}
	}

	if err := jsonextract.Save(newList); err != nil {
		panic(err.Error())
	}
}

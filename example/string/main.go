package main

import "github.com/tamboto2000/jsonextract"

func main() {
	jsons, err := jsonextract.FromFile("string.txt")
	if err != nil {
		panic(err.Error())
	}

	if err := jsonextract.Save(jsons); err != nil {
		panic(err.Error())
	}
}

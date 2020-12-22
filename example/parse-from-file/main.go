package main

import "github.com/tamboto2000/jsonextract/v2"

func main() {
	jsons, err := jsonextract.FromFile("test.txt")
	if err != nil {
		panic(err.Error())
	}

	// save result to path
	if err := jsonextract.SaveToPath(jsons, "from_file.json"); err != nil {
		panic(err.Error())
	}
}

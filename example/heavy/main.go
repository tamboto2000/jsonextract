package main

import "github.com/tamboto2000/jsonextract/v2"

func main() {
	jsons, err := jsonextract.FromFileWithOpt("test.html", jsonextract.Option{ParseObj: true})
	// jsons, err := jsonextract.FromFile("test.html")

	if err != nil {
		panic(err.Error())
	}

	if err := jsonextract.Save(jsons); err != nil {
		panic(err.Error())
	}
}

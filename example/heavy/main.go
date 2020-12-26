package main

import "github.com/tamboto2000/jsonextract/v2"

func main() {
	jsons, err := jsonextract.FromFileWithOpt("test.txt", jsonextract.Option{ParseArray: true, IgnoreEmptyArray: true})

	if err != nil {
		panic(err.Error())
	}

	if err := jsonextract.Save(jsons); err != nil {
		panic(err.Error())
	}
}

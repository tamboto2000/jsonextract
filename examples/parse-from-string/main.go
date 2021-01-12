package main

import "github.com/tamboto2000/jsonextract/v2"

func main() {
	str := `{
		"1": 1,
		"2": -2,
		"3": 0.3,
		"4": -0.4,		
		"5": true,
		"6": false,
		"7": null,
		"8": [1, -2, 0.3, -0.4, true, false, null, {"a":"b"}, 0.0e-1, 1e+2]
	}
	
	{}`

	jsons, err := jsonextract.FromString(str)
	if err != nil {
		panic(err.Error())
	}

	// save result to path
	if err := jsonextract.SaveToPath(jsons, "from_str.json"); err != nil {
		panic(err.Error())
	}
}

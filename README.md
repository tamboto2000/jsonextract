# JSONExtract

[![Go Reference](https://pkg.go.dev/badge/github.com/tamboto2000/jsonextract.svg)](https://pkg.go.dev/github.com/tamboto2000/jsonextract)

Package jsonextract is a small library for extracting JSON from a string, it extract a possible valid JSONs from a string or text. Right now only latin characters (a-z, A-Z, 0-9) are supported. This package did not guarantee 100% success rate of parsing, so it is highly recommended to check if the JSONs that you get is valid.

### Installation
JSONExtract require Go v1.14 or higher

```sh
$ GO111MODULE go get github.com/tamboto2000/jsonextract
```

### Example
```go
package main

import "github.com/tamboto2000/jsonextract"

func main() {
    // extract JSON from a file
	jsons, err := jsonextract.JSONFromFile("index.html")
	if err != nil {
		panic(err.Error())
	}

    // save result
	if err := jsonextract.SaveToPath(jsons, extracted); err != nil {
		panic(err.Error())
	}
	
	// validate
	valid, invalid := jsonextract.Validate(jsons)
	// save or use...
}
```

License
----

MIT
# JSONExtract

[![Go Reference](https://pkg.go.dev/badge/github.com/tamboto2000/jsonextract.svg)](https://pkg.go.dev/github.com/tamboto2000/jsonextract)

Package jsonextract is a small library for extracting JSON from a string, it extract a possible valid JSONs from a string or text. This package did not guarantee 100% success rate of parsing, so it is highly recommended to check if the JSONs that you get is valid.

# Update!

v2 is in progress! the next version will 100% guarantee every single JSON that parsed is valid! Features that will be included are:

- Preprocessed JSON, meaning that you can access the parsed JSON value directly without having to unmarshal the bytes!
- Target desired JSON kind/type to be parsed, so no more bulky result with unwanted JSON

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
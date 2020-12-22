# JSONExtract

[![Go Reference](https://pkg.go.dev/badge/github.com/tamboto2000/jsonextract.svg)](https://pkg.go.dev/github.com/tamboto2000/jsonextract/v2)

Package jsonextract is a library for extracting JSON from a given source, such as string, bytes, file, and io.Reader

# Update!
v2 is here! Features that included in this release:

- Preprocessed JSON, meaning that you can access the parsed JSON value directly without having to unmarshal the bytes!
- Target desired JSON kind/type to be parsed, so no more bulky result with unwanted JSON

### Installation
JSONExtract require Go v1.14 or higher

```sh
$ GO111MODULE go get github.com/tamboto2000/jsonextract/v2
```

# Examples

### Parse from string
```go
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
```

### Parse from bytes
```go
package main

import "github.com/tamboto2000/jsonextract/v2"

func main() {
	byts := []byte(`{
		"1": 1,
		"2": -2,
		"3": 0.3,
		"4": -0.4,		
		"5": true,
		"6": false,
		"7": null,
		"8": [1, -2, 0.3, -0.4, true, false, null, {"a":"b"}, 0.0e-1, 1e+2]
	}
	
	{}`)

	jsons, err := jsonextract.FromBytes(byts)
	if err != nil {
		panic(err.Error())
	}

	// save result to path
	if err := jsonextract.SaveToPath(jsons, "from_byts.json"); err != nil {
		panic(err.Error())
	}
}
```

### Parse from reader
```go
package main

import (
	"bytes"

	"github.com/tamboto2000/jsonextract/v2"
)

func main() {
	byts := []byte(`{
		"1": 1,
		"2": -2,
		"3": 0.3,
		"4": -0.4,		
		"5": true,
		"6": false,
		"7": null,
		"8": [1, -2, 0.3, -0.4, true, false, null, {"a":"b"}, 0.0e-1, 1e+2]
	}
	
	{}`)

	r := bytes.NewReader(byts)
	jsons, err := jsonextract.FromReader(r)
	if err != nil {
		panic(err.Error())
	}

	// save result to path
	if err := jsonextract.SaveToPath(jsons, "from_reader.json"); err != nil {
		panic(err.Error())
	}
}
```

### Parse from file
```go
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
```

See [Documentation](https://pkg.go.dev/github.com/tamboto2000/jsonextract) for more information

License
----

MIT
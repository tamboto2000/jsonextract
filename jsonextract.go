// Package jsonextract is a small library for extracting JSON from a string, it extract a possible valid JSONs from a string or text.
// Right now only latin characters (a-z, A-Z, 0-9) are supported.
// This package did not guarantee 100% success rate of parsing, so it is highly recommended to check if the JSONs that you get is valid
package jsonextract

import (
	"bytes"
	"io"
	"os"
)

// JSONFromStr extract every possible JSONs from a string
func JSONFromStr(data string) ([][]byte, error) {
	byts := []byte(data)
	return JSONFromBytes(byts)
}

// JSONFromBytes extract every possible JSONs in an array of bytes
func JSONFromBytes(str []byte) ([][]byte, error) {
	reader := bytes.NewReader(str)
	return parse(reader)
}

// JSONFromReader extract every possible JSONs in a io.Reader
func JSONFromReader(reader io.Reader) ([][]byte, error) {
	return parse(reader)
}

// JSONFromFile extract every possible JSONs in a file in path
func JSONFromFile(path string) ([][]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return JSONFromReader(f)
}

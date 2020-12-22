// Package jsonextract is a small library for extracting JSON from a string, it extract a possible valid JSONs from a string or text.
// Right now only latin characters (a-z, A-Z, 0-9) are supported.
// This package did not guarantee 100% success rate of parsing, so it is highly recommended to check if the JSONs that you get is valid
package jsonextract

import "os"

// Option determine what kind of objects should be parsed
type Option struct {
	ParseInt   bool
	ParseFloat bool
	ParseBool  bool
	ParseObj   bool
	ParseArray bool
	ParseNull  bool
}

// // FromBytes extract JSONs from bytes
// func FromBytes(byts []byte) ([]*JSON, error) {
// 	r := readFromBytes(byts)
// 	return runParser(r)
// }

// // FromReader extract JSONs from reader io.Reader
// func FromReader(reader io.Reader) ([]*JSON, error) {
// 	r, err := readFromReader(reader)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return runParser(r)
// }

// FromFile extract JSONs from file in path
func FromFile(path string) ([]*JSON, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	r, err := readFromReader(f)
	if err != nil {
		return nil, err
	}

	return parseAll(r)
}

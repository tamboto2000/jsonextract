// Package jsonextract is a library for extracting any valid JSONs from given source
package jsonextract

import (
	"io"
	"os"
)

// FromString extract JSONs from string
func FromString(str string) ([]*JSON, error) {
	r := readFromString(str)
	return parseAll(r)
}

// FromBytes extract JSONs from bytes
func FromBytes(byts []byte) ([]*JSON, error) {
	r := readFromBytes(byts)
	return parseAll(r)
}

// FromReader extract JSONs from reader io.Reader
func FromReader(reader io.Reader) ([]*JSON, error) {
	r, err := readFromReader(reader)
	if err != nil {
		return nil, err
	}

	return parseAll(r)
}

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

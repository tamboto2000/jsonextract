// Package jsonextract is a library for extracting any valid JSONs from given source
package jsonextract

import (
	"io"
	"os"
)

// // Option determine what kind of objects and with what criteria should be parsed
// type Option struct {
// 	ParseStr         bool
// 	ParseInt         bool
// 	ParseFloat       bool
// 	ParseBool        bool
// 	ParseObj         bool
// 	ParseArray       bool
// 	ParseNull        bool
// 	IgnoreEmptyStr   bool
// 	IgnoreZeroInt    bool
// 	IgnoreZeroFloat  bool
// 	IgnoreFalseBool  bool
// 	IgnoreTrueBool   bool
// 	IgnoreEmptyObj   bool
// 	IgnoreEmptyArray bool
// 	IgnoreNull       bool
// }

// // DefaultOption default option for parsing
// var DefaultOption = Option{
// 	ParseStr:   true,
// 	ParseInt:   true,
// 	ParseFloat: true,
// 	ParseBool:  true,
// 	ParseObj:   true,
// 	ParseArray: true,
// 	ParseNull:  true,
// }

// // FromString extract JSONs from string
// func FromString(str string) ([]*JSON, error) {
// 	r := readFromString(str)
// 	return runParser(r, DefaultOption)
// }

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

// // FromStringWithOpt extract JSONs from string with options
// func FromStringWithOpt(str string, opt Option) ([]*JSON, error) {
// 	r := readFromString(str)
// 	return runParser(r, opt)
// }

// // FromBytesWithOpt extract JSONs from bytes with options
// func FromBytesWithOpt(byts []byte, opt Option) ([]*JSON, error) {
// 	r := readFromBytes(byts)
// 	return runParser(r, opt)
// }

// // FromReaderWithOpt extract JSONs from reader io.Reader with options
// func FromReaderWithOpt(reader io.Reader, opt Option) ([]*JSON, error) {
// 	r, err := readFromReader(reader)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return runParser(r, opt)
// }

// // FromFileWithOpt extract JSONs from file in path with options
// func FromFileWithOpt(path string, opt Option) ([]*JSON, error) {
// 	f, err := os.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer f.Close()

// 	r, err := readFromReader(f)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return runParser(r, opt)
// }

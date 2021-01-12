package jsonextract

import (
	"errors"
	"io"
	"unicode"
)

// JSON kinds
const (
	Object = iota
	Array
	String
	Integer
	Float
	Boolean
	Null
)

// JSON represent JSON val
type JSON struct {
	Kind int
	val  interface{}
	raw  []rune
}

func (json *JSON) push(char rune) {
	json.raw = append(json.raw, char)
}

func (json *JSON) pushRns(chars []rune) {
	json.raw = append(json.raw, chars...)
}

// RawRunes return parsed raw runes of JSON
func (json *JSON) RawRunes() []rune {
	return json.raw
}

// RawBytes return parsed raw bytes of JSON
func (json *JSON) RawBytes() []byte {
	return runesToUTF8(json.raw)
}

// Array return array of JSON. Will return error if Kind != Array
func (json *JSON) Array() ([]*JSON, error) {
	if json.Kind != Array {
		return nil, errors.New("value is not array")
	}

	return json.val.([]*JSON), nil
}

// String return string value. Will return error if Kind != String
func (json *JSON) String() (string, error) {
	if json.Kind != String {
		return "", errors.New("value is not string")
	}

	return json.val.(string), nil
}

// Integer return int value. Will return error if Kind != Integer
func (json *JSON) Integer() (int64, error) {
	if json.Kind != Integer {
		return 0, errors.New("value is not int")
	}

	return json.val.(int64), nil
}

// Float return float value. Will return error if Kind != Float
func (json *JSON) Float() (float64, error) {
	if json.Kind != Float {
		return 0, errors.New("value is not float")
	}

	return json.val.(float64), nil
}

// Boolean return bool value. Will return error if Kind != Boolean
func (json *JSON) Boolean() (bool, error) {
	if json.Kind != Boolean {
		return false, errors.New("value is not bool")
	}

	return json.val.(bool), nil
}

// Null will not return any value except error. Will return error if Kind != Boolean
func (json *JSON) Null() error {
	if json.Kind != Null {
		return errors.New("value is not null")
	}

	return nil
}

func parseAll(r reader) ([]*JSON, error) {
	jsons := make([]*JSON, 0)
	for {
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		json, err := parse(r, char)
		if err != nil {
			if err == errInvalid {
				r.UnreadRune()
				continue
			}

			return nil, err
		}

		if json != nil {
			jsons = append(jsons, json)
		}
	}

	return jsons, nil
}

func parse(r reader, char rune) (*JSON, error) {
	// parse string
	if char == '"' {
		json, err := parseString(r)
		if err != nil {
			return nil, err
		}

		return json, nil
	}

	// parse numeric
	if unicode.IsNumber(char) || char == '-' {
		json, err := parseNumeric(r, char)
		if err != nil {
			return nil, err
		}

		return json, nil
	}

	// parse boolean
	if char == 't' || char == 'f' {
		json, err := parseBool(r, char)
		if err != nil {
			return nil, err
		}

		return json, nil
	}

	// parse null
	if char == 'n' {
		json, err := parseNull(r)
		if err != nil {
			return nil, err
		}

		return json, nil
	}

	// parse array
	if char == '[' {
		json, err := parseArray(r)
		if err != nil {
			return nil, err
		}

		return json, nil
	}

	// parse object
	if char == '{' {
		json, err := parseObj(r)
		if err != nil {
			return nil, err
		}

		return json, nil
	}

	return nil, nil
}

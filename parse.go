package jsonextract

import (
	jsonenc "encoding/json"
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
	// this will be assigned if value inside array of object
	parent *JSON
}

// Reparsing JSON to produce new raw JSON bytes and runes.
// This method is called if value inside JSON is edited
func (json *JSON) reParse() error {
	if json.Kind != Object && json.Kind != Array {
		byts, err := jsonenc.Marshal(json.val)
		if err != nil {
			return err
		}

		// currently I don't know how to directly convert []byte to []rune...
		json.raw = []rune(string(byts))
	}

	if json.Kind == Object {
		keyVals := json.val.(map[string]*JSON)
		newRaw := make([]rune, 0)
		newRaw = append(newRaw, '{')
		objLen := len(keyVals)
		for key, val := range keyVals {
			objLen--
			if err := val.reParse(); err != nil {
				return err
			}

			// Key sequence
			newRaw = append(newRaw, '"')
			newRaw = append(newRaw, []rune(key)...)
			newRaw = append(newRaw, '"')
			newRaw = append(newRaw, ':')

			// value
			newRaw = append(newRaw, val.RawRunes()...)

			// delimiter
			if objLen == 0 {
				newRaw = append(newRaw, '}')
			} else {
				newRaw = append(newRaw, ',')
			}
		}

		json.raw = newRaw
	}

	if json.Kind == Array {
		vals := json.Array()
		newRaw := make([]rune, 0)
		newRaw = append(newRaw, '[')
		arrLen := len(vals)
		for i, val := range vals {
			if err := val.reParse(); err != nil {
				return err
			}

			// value
			newRaw = append(newRaw, val.RawRunes()...)

			// delimiter
			if i == arrLen-1 {
				newRaw = append(newRaw, ']')
			} else {
				newRaw = append(newRaw, ',')
			}
		}

		json.raw = newRaw
	}

	return nil
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

// Object return map of JSON. Will panic if Kind != Object
func (json *JSON) Object() map[string]*JSON {
	if json.Kind != Object {
		panic("value is not object")
	}

	return json.val.(map[string]*JSON)
}

// Array return array of JSON. Will panic if Kind != Array
func (json *JSON) Array() []*JSON {
	if json.Kind != Array {
		panic("value is not array")
	}

	return json.val.([]*JSON)
}

// String return string value. Will panic if Kind != String
func (json *JSON) String() (string, error) {
	if json.Kind != String {
		return "", errors.New("value is not string")
	}

	return json.val.(string), nil
}

// Integer return int value. Will panic if Kind != Integer
func (json *JSON) Integer() int64 {
	if json.Kind != Integer {
		panic("value is not int")
	}

	return json.val.(int64)
}

// Float return float value. Will panic if Kind != Float
func (json *JSON) Float() float64 {
	if json.Kind != Float {
		panic("value is not float")
	}

	return json.val.(float64)
}

// Boolean return bool value. Will panic if Kind != Boolean
func (json *JSON) Boolean() bool {
	if json.Kind != Boolean {
		panic("value is not bool")
	}

	return json.val.(bool)
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

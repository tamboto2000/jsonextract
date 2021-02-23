package jsonextract

import (
	"bytes"
	jsonenc "encoding/json"
	"fmt"
	"io"
	"reflect"
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

// JSON represent json val
type JSON struct {
	kind   int
	val    interface{}
	raw    []rune
	parent *JSON
}

// reparsing JSON to produce new raw json bytes and runes.
// This method is called if value inside JSON is edited
func (json *JSON) reParse() {
	if json.kind != Object && json.kind != Array {
		byts, _ := jsonenc.Marshal(json.val)

		// convert bytes to runes
		json.raw = readAllRunes(bytes.NewReader(byts))
	}

	if json.kind == Object {
		keyVals := json.val.(map[interface{}]*JSON)
		newRaw := make([]rune, 0)
		newRaw = append(newRaw, '{')
		objLen := len(keyVals)
		for key, val := range keyVals {
			objLen--
			val.reParse()

			switch key.(type) {
			case string:
				// Key sequence for string field name/key
				newRaw = append(newRaw, '"')
				newRaw = append(newRaw, []rune(key.(string))...)
				newRaw = append(newRaw, '"')
				newRaw = append(newRaw, ':')

			default:
				// Key sequence for integer field name/key
				keyStr := fmt.Sprintf("%d", key)
				newRaw = append(newRaw, '"')
				newRaw = append(newRaw, []rune(keyStr)...)
				newRaw = append(newRaw, '"')
				newRaw = append(newRaw, ':')
			}

			// value
			newRaw = append(newRaw, val.raw...)

			// delimiter
			if objLen == 0 {
				newRaw = append(newRaw, '}')
			} else {
				newRaw = append(newRaw, ',')
			}
		}

		json.raw = newRaw
	}

	if json.kind == Array {
		vals := json.Array()
		newRaw := make([]rune, 0)
		newRaw = append(newRaw, '[')
		arrLen := len(vals)
		for i, val := range vals {
			val.reParse()

			// value
			newRaw = append(newRaw, val.raw...)

			// delimiter
			if i == arrLen-1 {
				newRaw = append(newRaw, ']')
			} else {
				newRaw = append(newRaw, ',')
			}
		}

		json.raw = newRaw
	}
}

func (json *JSON) push(char rune) {
	json.raw = append(json.raw, char)
}

func (json *JSON) pushRns(chars []rune) {
	json.raw = append(json.raw, chars...)
}

// Kind return json kind
func (json *JSON) Kind() int {
	return json.kind
}

// Runes return raw runes of extracted json
func (json *JSON) Runes() []rune {
	return json.raw
}

// Bytes return raw bytes of extracted json
func (json *JSON) Bytes() []byte {
	return runesToUTF8(json.raw)
}

// Object return map of JSON. Will panic if kind is not Object
func (json *JSON) Object() map[interface{}]*JSON {
	if json.kind != Object {
		panic("value is not object")
	}

	return json.val.(map[interface{}]*JSON)
}

// Array return array of JSON. Will panic if kind is not Array
func (json *JSON) Array() []*JSON {
	if json.kind != Array {
		panic("value is not array")
	}

	return json.val.([]*JSON)
}

// String return string value. Will panic if kind is not String
func (json *JSON) String() string {
	if json.kind != String {
		panic("value is not string")
	}

	return json.val.(string)
}

// Integer return int value. Will panic if kind is not Integer
func (json *JSON) Integer() int64 {
	if json.kind != Integer {
		panic("value is not int")
	}

	return convertIntTo64(json.val)
}

// Float return float value. Will panic if kind is not Float
func (json *JSON) Float() float64 {
	if json.kind != Float {
		panic("value is not float")
	}

	return convertFloatTo64(json.val)
}

// Boolean return bool value. Will panic if kind is not Boolean
func (json *JSON) Boolean() bool {
	if json.kind != Boolean {
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

// convert any type of integer to int64
func convertIntTo64(i interface{}) int64 {
	refVal := reflect.ValueOf(i)

	// signed integer
	if refVal.Kind() == reflect.Int {
		return int64(i.(int))
	}

	if refVal.Kind() == reflect.Int8 {
		return int64(i.(int8))
	}

	if refVal.Kind() == reflect.Int16 {
		return int64(i.(int16))
	}

	if refVal.Kind() == reflect.Int32 {
		return int64(i.(int32))
	}

	// unsigned integer
	if refVal.Kind() == reflect.Uint {
		return int64(i.(uint))
	}

	if refVal.Kind() == reflect.Uint8 {
		return int64(i.(uint8))
	}

	if refVal.Kind() == reflect.Uint16 {
		return int64(i.(uint16))
	}

	if refVal.Kind() == reflect.Uint32 {
		return int64(i.(uint32))
	}

	if refVal.Kind() == reflect.Uint64 {
		return int64(i.(uint64))
	}

	return i.(int64)
}

// convert any type of float to float64
func convertFloatTo64(i interface{}) float64 {
	refVal := reflect.ValueOf(i)

	if refVal.Kind() == reflect.Float32 {
		return float64(i.(float32))
	}

	return i.(float64)
}

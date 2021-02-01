package jsonextract

import (
	"bytes"
	jsonenc "encoding/json"
	"errors"
	"reflect"
)

// EditStr edit string value. Will panic if JSON.Kind != String
func (json *JSON) EditStr(str string) {
	if json.kind != String {
		panic("value is not string")
	}

	json.val = str

	getParent(json).reParse()
}

// EditInt edit int value. Will panic if JSON.Kind != Integer
func (json *JSON) EditInt(i int64) {
	if json.kind != Integer {
		panic("value is not int")
	}

	json.val = i

	getParent(json).reParse()
}

// EditFloat edit float value. Will panic if JSON.Kind != Float
func (json *JSON) EditFloat(i float64) {
	if json.kind != Float {
		panic("value is not float")
	}

	json.val = i

	getParent(json).reParse()
}

// EditBool edit bool value. Will panic if JSON.Kind != Boolean
func (json *JSON) EditBool(b bool) {
	if json.kind != Boolean {
		panic("value is not bool")
	}

	json.val = b

	getParent(json).reParse()
}

// DeleteField delete object field by key. Will panic if JSON.Kind != Object
func (json *JSON) DeleteField(key string) {
	if json.kind != Object {
		panic("value is not object")
	}

	keyVal := json.val.(map[string]*JSON)
	delete(keyVal, key)
	json.val = keyVal

	getParent(json).reParse()
}

// DeleteElm delete element on index i. Will panic if JSON.Kind != Array
func (json *JSON) DeleteElm(i int) {
	if json.kind != Array {
		panic("value is not array")
	}

	vals := json.val.([]*JSON)
	if len(vals) == 0 {
		return
	}

	vals = append(vals[:i], vals[i+1:]...)
	json.val = vals

	getParent(json).reParse()
}

// AddField add new field to object. Will panic if JSON.Kind != Object
// Valid type for val are string, int, float, map[string]interface{}, struct, and nil
func (json *JSON) AddField(key string, val interface{}) error {
	if json.kind != Object {
		panic("value is not object")
	}

	keyval := json.val.(map[string]*JSON)

	// nil value
	if val == nil {
		newJSON := &JSON{
			kind:   Null,
			val:    nil,
			raw:    []rune("null"),
			parent: json,
		}

		keyval[key] = newJSON

		getParent(json).reParse()
		return nil
	}

	refVal := reflect.ValueOf(val)
	for refVal.Kind() == reflect.Interface || refVal.Kind() == reflect.Ptr {
		refVal = refVal.Elem()
	}

	// non array or map
	if refVal.Kind() != reflect.Slice && refVal.Kind() != reflect.Array &&
		refVal.Kind() != reflect.Map {

		// generate new json bytes
		jsbyts, _ := jsonenc.Marshal(val)
		newJSON := &JSON{
			val:    val,
			raw:    readAllRunes(bytes.NewReader(jsbyts)),
			parent: json,
		}

		// integer
		if isValIsInteger(refVal) {
			newJSON.kind = Integer
		}

		// float
		if isValFloat(refVal) {
			newJSON.kind = Float
		}

		// string
		if refVal.Kind() == reflect.String {
			newJSON.kind = String
		}

		getParent(json).reParse()
	}

	return errors.New("type not supported")
}

// check if reflection value kind is integer
func isValIsInteger(val reflect.Value) bool {
	// integer
	if val.Kind() == reflect.Int || val.Kind() == reflect.Int8 ||
		val.Kind() == reflect.Int16 || val.Kind() == reflect.Int32 ||
		val.Kind() == reflect.Int64 {
		return true
	}

	// unsigned integer
	if val.Kind() == reflect.Uint || val.Kind() == reflect.Uint8 ||
		val.Kind() == reflect.Uint16 || val.Kind() == reflect.Uint32 ||
		val.Kind() == reflect.Uint64 {
		return true
	}

	return false
}

func isValFloat(val reflect.Value) bool {
	if val.Kind() == reflect.Float32 || val.Kind() == reflect.Float64 {
		return true
	}

	return false
}

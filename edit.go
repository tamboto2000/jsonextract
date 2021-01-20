package jsonextract

import (
	"errors"
)

// EditStr edit string value. Wil return error if JSON.Kind != String
func (json *JSON) EditStr(str string) error {
	if json.Kind != String {
		return errors.New("value is not string")
	}

	json.val = str

	return getParent(json).reParse()
}

// EditInt edit int value. Wil return error if JSON.Kind != Integer
func (json *JSON) EditInt(i int64) error {
	if json.Kind != Integer {
		return errors.New("value is not int")
	}

	json.val = i

	return getParent(json).reParse()
}

// EditFloat edit float value. Wil return error if JSON.Kind != Float
func (json *JSON) EditFloat(i float64) error {
	if json.Kind != Float {
		return errors.New("value is not float")
	}

	json.val = i

	return getParent(json).reParse()
}

// EditBool edit bool value. Wil return error if JSON.Kind != Boolean
func (json *JSON) EditBool(b bool) error {
	if json.Kind != Boolean {
		return errors.New("value is not bool")
	}

	json.val = b

	return getParent(json).reParse()
}

// DeleteField delete object field by key. Will return error if JSON.Kind != Object
func (json *JSON) DeleteField(key string) error {
	if json.Kind != Object {
		return errors.New("value is not object")
	}

	keyVal := json.val.(map[string]*JSON)
	delete(keyVal, key)
	json.val = keyVal

	return getParent(json).reParse()
}

// DeleteElm delete element on index i. Will return error if JSON.Kind != Array
func (json *JSON) DeleteElm(i int) error {
	if json.Kind != Array {
		return errors.New("value is not array")
	}

	vals := json.val.([]*JSON)
	if len(vals) == 0 {
		return nil
	}

	vals = append(vals[:i], vals[i+1:]...)
	json.val = vals

	return getParent(json).reParse()
}

// AddField add new field to object. Will return error if JSON.Kind != Object
// func (json *JSON) AddField(key string, val interface{}) error {
// 	if json.Kind != Object {
// 		return errors.New("value is not object")
// 	}

// 	// refVal := reflect.ValueOf(val)

// 	return nil
// }

package jsonextract

import (
	"errors"
	"fmt"
	"reflect"
)

// SetStr set string value. Will panic if kind is not String
func (json *JSON) SetStr(str string) {
	if json.kind != String {
		panic("value is not string")
	}

	json.val = str

	getParent(json).reParse()
}

// SetInt set int value. Will panic if kind is not Integer
func (json *JSON) SetInt(i int64) {
	if json.kind != Integer {
		panic("value is not int")
	}

	json.val = i

	getParent(json).reParse()
}

// SetFloat set float value. Will panic if kind is not Float
func (json *JSON) SetFloat(i float64) {
	if json.kind != Float {
		panic("value is not float")
	}

	json.val = i

	getParent(json).reParse()
}

// SetBool set bool value. Will panic if kind is not Boolean
func (json *JSON) SetBool(b bool) {
	if json.kind != Boolean {
		panic("value is not bool")
	}

	json.val = b

	getParent(json).reParse()
}

// DeleteField delete object field by key. Will panic if kind is not Object.
// key must be string or an integer type, otherwise panic will occur
func (json *JSON) DeleteField(key interface{}) bool {
	if json.kind != Object {
		panic("value is not object")
	}

	valKey := reflect.ValueOf(key)
	for valKey.Kind() == reflect.Ptr || valKey.Kind() == reflect.Interface {
		valKey = valKey.Elem()
	}

	if valKey.Kind() != reflect.String && !isValInteger(valKey) {
		panic("key must be string or an integer type")
	}

	var keyStr string
	if valKey.Kind() == reflect.String {
		keyStr = valKey.String()
	} else if isValInteger(valKey) {
		keyStr = fmt.Sprintf("%d", valKey.Interface())
	}

	keyVal := json.val.(map[interface{}]*JSON)
	if _, ok := keyVal[keyStr]; !ok {
		return false
	}

	delete(keyVal, keyStr)
	json.val = keyVal

	getParent(json).reParse()

	return true
}

// DeleteItem delete element on index i. Will panic if kind is not Array.
// Will panic if i bigger than items count minus 1 (len(vals)-1, the last index).
// If item is existed, true is returned, otherwise false
func (json *JSON) DeleteItem(i int) bool {
	if json.kind != Array {
		panic("value is not array")
	}

	vals := json.val.([]*JSON)
	valLen := len(vals)
	if valLen == 0 {
		return false
	}

	if i > valLen-1 {
		panic("out of range")
	}

	vals = append(vals[:i], vals[i+1:]...)
	json.val = vals

	getParent(json).reParse()

	return true
}

// AddField add new field to object. Will panic if kind is not Object.
// Will panic if val is invalid json value.
// Valid key is string or integer type.
// If val is map, and key is not string nor an integer type, panic will occur
func (json *JSON) AddField(key interface{}, val interface{}) {
	if json.kind != Object {
		panic("value is not object")
	}

	keyval := json.val.(map[interface{}]*JSON)

	newJSON, err := generateJSON(val)
	if err != nil {
		panic(err.Error())
	}

	newJSON.parent = json

	keyVal := reflect.ValueOf(key)
	for keyVal.Kind() == reflect.Interface {
		keyVal = keyVal.Elem()
	}

	if keyVal.Kind() != reflect.String && !isValInteger(keyVal) {
		panic("key must be string or an integer type")
	}

	// var fk string
	// if keyVal.Kind() == reflect.String {
	// 	fk = keyVal.String()
	// } else if isValInteger(keyVal) {
	// 	fk = fmt.Sprintf("%d", keyVal.Interface())
	// }

	keyval[key] = newJSON
	getParent(json).reParse()
}

// AddItem add new item in json array. Will panic if kind is not Array.
// Will panic if val is invalid json value
func (json *JSON) AddItem(val interface{}) {
	if json.kind != Array {
		panic("value is not array")
	}

	newJSON, err := generateJSON(val)
	if err != nil {
		panic(err.Error())
	}

	vals := json.val.([]*JSON)
	vals = append(vals, newJSON)

	json.val = vals
	getParent(json).reParse()
}

// Len return items count inside JSON.
// If JSON is not Array or Object, Len will return 0, so please check the JSON kind first
func (json *JSON) Len() int {
	if json.kind == Object {
		vals := json.val.(map[interface{}]*JSON)
		return len(vals)
	}

	if json.kind == Array {
		vals := json.val.([]*JSON)
		return len(vals)
	}

	return 0
}

// check if reflection value kind is integer
func isValInteger(val reflect.Value) bool {
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

// generate JSON from val
func generateJSON(val interface{}) (*JSON, error) {
	// nil value
	if val == nil {
		newJSON := &JSON{
			kind: Null,
			val:  nil,
		}

		return newJSON, nil
	}

	// if val is pointer or interface, iterate Value.Elem to get the real value
	refVal := reflect.ValueOf(val)
	for refVal.Kind() == reflect.Interface || refVal.Kind() == reflect.Ptr {
		refVal = refVal.Elem()
	}

	// string
	if refVal.Kind() == reflect.String {
		var str string
		newRefVal := reflect.ValueOf(&str)
		newRefVal.Elem().Set(refVal)
		newJSON := &JSON{
			kind: String,
			val:  str,
		}

		return newJSON, nil
	}

	// integer
	if isValInteger(refVal) {
		var i interface{}
		newRefVal := reflect.ValueOf(&i)
		newRefVal.Elem().Set(refVal)
		newJSON := &JSON{
			kind: Integer,
			val:  i,
		}

		return newJSON, nil
	}

	// float
	if isValFloat(refVal) {
		var i interface{}
		newRefVal := reflect.ValueOf(&i)
		newRefVal.Elem().Set(refVal)
		newJSON := &JSON{
			kind: Float,
			val:  i,
		}

		return newJSON, nil
	}

	// bool
	if refVal.Kind() == reflect.Bool {
		var b bool
		newRefVal := reflect.ValueOf(&b)
		newRefVal.Elem().Set(refVal)
		newJSON := &JSON{
			kind: Float,
			val:  b,
		}

		return newJSON, nil
	}

	// array
	if refVal.Kind() == reflect.Array || refVal.Kind() == reflect.Slice {
		newJSON := &JSON{kind: Array}
		arr := make([]*JSON, refVal.Len())
		for i := 0; i < refVal.Len(); i++ {
			item := refVal.Index(i)

			if item.IsZero() {
				njs, _ := generateJSON(nil)
				arr[i] = njs

				continue
			}

			for item.Kind() == reflect.Ptr || item.Kind() == reflect.Interface {
				item = item.Elem()
			}

			njs, err := generateJSON(item.Interface())
			if err != nil {
				return nil, err
			}

			arr[i] = njs
		}

		newJSON.val = arr
		return newJSON, nil
	}

	// object from map
	if refVal.Kind() == reflect.Map {
		newJSON := &JSON{kind: Object}
		newMap := make(map[interface{}]*JSON)
		for _, key := range refVal.MapKeys() {
			for key.Kind() == reflect.Interface {
				key = key.Elem()
			}

			if key.Kind() != reflect.String && !isValInteger(key) {
				return nil, errors.New("map key must be string or an integer type")
			}

			item := refVal.MapIndex(key)

			if item.IsZero() {
				njs, _ := generateJSON(nil)
				newMap[key.Interface()] = njs

				continue
			}

			for item.Kind() == reflect.Ptr || item.Kind() == reflect.Interface {
				item = item.Elem()
			}

			njs, err := generateJSON(item.Interface())
			if err != nil {
				return nil, err
			}

			newMap[key.Interface()] = njs
		}

		newJSON.val = newMap
		return newJSON, nil
	}

	// object from struct
	if refVal.Kind() == reflect.Struct {
		valType := refVal.Type()
		newJSON := &JSON{kind: Object}
		newMap := make(map[interface{}]*JSON)
		for i := 0; i < refVal.NumField(); i++ {
			// field name
			var fk string

			// field from reflect.Value
			fv := refVal.Field(i)
			// field from reflect.Type
			ft := valType.Field(i)

			// if there is no json tag, use field name as json field name
			if fk = ft.Tag.Get("json"); fk == "" {
				fk = ft.Name
			} 

			nj, err := generateJSON(fv.Interface())
			if err != nil {
				return nil, err
			}

			newMap[fk] = nj
		}

		newJSON.val = newMap
		return newJSON, nil
	}

	return nil, errors.New("type unsupported")
}

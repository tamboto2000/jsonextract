package jsonextract

import (
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

// DeleteItem delete element on index i. Will panic if JSON.Kind != Array
func (json *JSON) DeleteItem(i int) {
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

// AddField add new field to object. Will panic if JSON.Kind != Object.
// Will panic if val is invalid json value.
// If val is map, and key is not string, panic will occur
func (json *JSON) AddField(key string, val interface{}) {
	if json.kind != Object {
		panic("value is not object")
	}

	keyval := json.val.(map[string]*JSON)

	newJSON, err := generateJSON(val)
	if err != nil {
		panic(err.Error())
	}

	newJSON.parent = json

	keyval[key] = newJSON
	getParent(json).reParse()
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

func isValIsFloat(val reflect.Value) bool {
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
	if isValIsInteger(refVal) {
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
	if isValIsFloat(refVal) {
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
		newMap := make(map[string]*JSON)
		for _, key := range refVal.MapKeys() {
			if key.Kind() != reflect.String {
				return nil, errors.New("map key must be string")
			}

			item := refVal.MapIndex(key)

			if item.IsZero() {
				njs, _ := generateJSON(nil)
				newMap[key.String()] = njs

				continue
			}

			for item.Kind() == reflect.Ptr || item.Kind() == reflect.Interface {
				item = item.Elem()
			}

			njs, err := generateJSON(item.Interface())
			if err != nil {
				return nil, err
			}

			newMap[key.String()] = njs
		}

		newJSON.val = newMap
		return newJSON, nil
	}

	// object from struct
	if refVal.Kind() == reflect.Struct {
		valType := refVal.Type()
		newJSON := &JSON{kind: Object}
		newMap := make(map[string]*JSON)
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

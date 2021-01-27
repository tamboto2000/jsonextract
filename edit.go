package jsonextract

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
// func (json *JSON) AddField(key string, val interface{}) error {
// 	if json.kind != Object {
// 		panic("value is not object")
// 	}

// 	refVal := reflect.ValueOf(val)
// 	if refVal.Kind() == reflect.Ptr || refVal.Kind() == reflect.Interface {
// 		refVal = refVal.Elem()
// 		if refVal.IsZero() {

// 		}
// 	}

// 	return nil
// }

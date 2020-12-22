package jsonextract

import (
	"io"
)

var (
	openCurlBrack  = byte(123)
	closeCurlBrack = byte(125)
	openBrack      = byte(91)
	closeBrack     = byte(93)
	quot           = byte(34)
	colon          = byte(58)
	coma           = byte(44)
	backSlash      = byte(92)
	slash          = byte(47)
	minus          = byte(45)
	plus           = byte(43)
	dot            = byte(46)
)

var (
	letters   = []byte{97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90}
	numbers   = []byte{49, 50, 51, 52, 53, 54, 55, 56, 57, 48}
	escapable = []byte{97, 98, 116, 110, 118, 102, 114, 117, slash, backSlash, quot}
	// "\a\b\t\n\v\f\r "
	syntax       = []byte{7, 8, 9, 10, 11, 12, 13, 32}
	nullStr      = []byte{110, 117, 108, 108}
	trueStr      = []byte{116, 114, 117, 101}
	falseStr     = []byte{102, 97, 108, 115, 101}
	exponentChar = []byte{101, 69}
	hexChars     = []byte{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102}
)

// Kind represent JSON kind or type
type Kind int

// Every JSON kind
const (
	String = Kind(iota)
	Int
	// float64
	Float
	Boolean
	Object
	Array
	Null
)

// Raw store raw JSON bytes
type Raw struct {
	byts []byte
}

// Bytes return stored raw json
func (r *Raw) Bytes() []byte {
	return r.byts
}

func (r *Raw) push(byt byte) {
	r.byts = append(r.byts, byt)
}

func (r *Raw) pushBytes(byts []byte) {
	r.byts = append(r.byts, byts...)
}

// JSON contains JSON data and its kind
type JSON struct {
	// JSON kind/type
	Kind Kind
	// Only apply to Int and Float
	WithExponent bool
	// Key and value pair, like "key" : "val", only assigned when Kind == Object
	KeyVal map[string]*JSON
	// Values for Array, only assigned if Kind == Array
	Vals []*JSON
	// If data is just a single object, like Int, Float, Boolean, not Object nor Array
	Val interface{}

	// Stores raw JSON bytes
	Raw *Raw
}

// parse all detected JSON
func parseAll(r reader) ([]*JSON, error) {
	jsons := make([]*JSON, 0)
	for {
		json, err := parse(r)
		if err == nil {
			jsons = append(jsons, json)
		}

		if err != nil && err == io.EOF {
			break
		}
	}

	return jsons, nil
}

// parse with selected JSON kind/type
func parseWithOpt(r reader, opt Option) ([]*JSON, error) {
	jsons := make([]*JSON, 0)
	for {
		json, err := parse(r)
		if err == nil {
			if json.Kind == String && opt.ParseStr {
				jsons = append(jsons, json)
			}

			if json.Kind == Int && opt.ParseInt {
				jsons = append(jsons, json)
			}

			if json.Kind == Float && opt.ParseFloat {
				jsons = append(jsons, json)
			}

			if json.Kind == Boolean && opt.ParseBool {
				jsons = append(jsons, json)
			}

			if json.Kind == Null && opt.ParseNull {
				jsons = append(jsons, json)
			}

			if json.Kind == Array && opt.ParseArray {
				jsons = append(jsons, json)
			}

			if json.Kind == Object && opt.ParseObj {
				jsons = append(jsons, json)
			}
		}

		if err != nil && err == io.EOF {
			break
		}
	}

	return jsons, nil
}

// parse any detected JSON
func parse(r reader) (*JSON, error) {
	for {
		char, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		if isCharSyntax(char) {
			continue
		}

		if !isCharValidBeginObj(char) {
			continue
		}

		// if char is a beginning for an object, start finding JSON
		r.UnreadByte()
		break
	}

	// try to parse string
	str, err := parseStr(r)
	if err != nil {
		r.UnreadByte()

		if err != errUnmatch {
			return nil, err
		}
	} else {
		return str, nil
	}

	// try to parse numeric
	num, err := parseNum(r)
	if err != nil {
		r.UnreadByte()

		if err != errUnmatch {
			return nil, err
		}
	} else {
		return num, nil
	}

	// try to parse boolean
	bl, err := parseBool(r)
	if err != nil {
		r.UnreadByte()

		if err != errUnmatch {
			return nil, err
		}
	} else {
		return bl, nil
	}

	// try to parse null
	nl, err := parseNull(r)
	if err != nil {
		r.UnreadByte()

		if err != errUnmatch {
			return nil, err
		}
	} else {
		return nl, nil
	}

	// try to parse array
	arr, err := parseArr(r)
	if err != nil {
		r.UnreadByte()

		if err != errUnmatch {
			return nil, err
		}
	} else {
		return arr, nil
	}

	// try to parse object
	obj, err := parseObj(r)
	if err != nil {
		r.UnreadByte()

		if err != errUnmatch {
			return nil, err
		}
	} else {
		return obj, nil
	}

	return nil, nil
}

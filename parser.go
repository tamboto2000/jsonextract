package jsonextract

import (
	"io"
	"strconv"
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
	dot            = byte(46)
)

var (
	letters   = []byte{97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90}
	numbers   = []byte{49, 50, 51, 52, 53, 54, 55, 56, 57, 48}
	escapable = []byte{97, 98, 116, 110, 118, 102, 114, slash, backSlash, quot}
	// "\a\b\t\n\v\f\r "
	syntax   = []byte{7, 8, 9, 10, 11, 12, 13, 32}
	nullStr  = []byte{110, 117, 108, 108}
	trueStr  = []byte{116, 114, 117, 101}
	falseStr = []byte{102, 97, 108, 115, 101}
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
	// DELETE
	// fmt.Println(string([]byte{byt}))
	r.byts = append(r.byts, byt)
}

func (r *Raw) pushBytes(byts []byte) {
	r.byts = append(r.byts, byts...)
}

// JSON contains JSON data and its kind
type JSON struct {
	Kind Kind
	// Key and value pair, like "key" : "val", only assigned when Kind == Object
	KeyVal map[string]*JSON
	// Values for Array, only assigned if Kind == Array
	Vals []JSON
	// If data is just a single object, like Int, Float, Boolen, not Object nor Array
	Val interface{}

	// Stores raw JSON bytes
	Raw *Raw
}

func parse(r reader) ([]*JSON, error) {
	jsons := make([]*JSON, 0)
	for {
		// check if char is a probably an object
		char, err := r.ReadByte()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}

			break
		}

		if !isCharValidBeginObj(char) {
			continue
		}

		r.UnreadByte()

		if json, err := parseStr(r); err != nil {
			if err == errInvalid || err == errUnmatch {
				r.UnreadByte()
			} else {
				if err == io.EOF {
					return jsons, nil
				}

				return nil, err
			}
		} else {
			jsons = append(jsons, json)
		}

		if json, err := parseNum(r); err != nil {
			if err == errInvalid || err == errUnmatch {
				r.UnreadByte()
			} else {
				if err == io.EOF {
					return jsons, nil
				}

				return nil, err
			}
		} else {
			jsons = append(jsons, json)
		}

		if json, err := parseNull(r); err != nil {
			if err == errInvalid || err == errUnmatch {
				r.UnreadByte()
			} else {
				if err == io.EOF {
					return jsons, nil
				}

				return nil, err
			}
		} else {
			jsons = append(jsons, json)
		}

		if json, err := parseBool(r); err != nil {
			if err == errInvalid || err == errUnmatch {
				r.UnreadByte()
			} else {
				if err == io.EOF {
					return jsons, nil
				}

				return nil, err
			}
		} else {
			jsons = append(jsons, json)
		}

		if json, err := parseObj(r); err != nil {
			if err == errInvalid || err == errUnmatch {
				r.UnreadByte()
			} else {
				if err == io.EOF {
					return jsons, nil
				}

				return nil, err
			}
		} else {
			jsons = append(jsons, json)
		}
	}

	return jsons, nil
}

func parseStr(r reader) (*JSON, error) {
	raw := new(Raw)
	char, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	onEscape := false
	if char == quot {
		raw.push(char)
		json := &JSON{Kind: String, Raw: raw}
		temp := make([]byte, 0)

		for {
			char, err := r.ReadByte()
			if err != nil {
				return nil, err
			}

			if onEscape {
				if isCharEscapable(char) {
					onEscape = false
					raw.push(char)
					temp = append(temp, char)
					continue
				}

				return nil, errInvalid
			}

			if char == backSlash {
				onEscape = true

				raw.push(char)
				temp = append(temp, char)
				continue
			}

			if char == quot {
				raw.push(char)
				json.Val = string(temp)
				return json, nil
			}

			if escChar := escapeChar(char); escChar != nil {
				// DELETE
				// fmt.Println(string(escChar))
				for _, c := range escChar {
					raw.push(c)
				}

				temp = append(temp, char)
				continue
			}

			raw.push(char)
			temp = append(temp, char)
		}
	}

	// DELETE
	// fmt.Println("unmatch string", string([]byte{char}))
	return nil, errUnmatch
}

func parseNum(r reader) (*JSON, error) {
	char, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	// DELETE
	// fmt.Println(string([]byte{char}))

	isMinus := false
	raw := new(Raw)
	if char == minus {
		// check if the next char is numeric
		c, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		if !isCharNumber(c) {
			// DELETE
			// fmt.Println("invalid num", string([]byte{c}))
			return nil, errInvalid
		}

		raw.push(char)
		char = c
		r.UnreadByte()
		isMinus = true
	}

	if isCharNumber(char) {
		json := &JSON{Kind: Int, Raw: raw}
		if !isMinus {
			raw.push(char)
		}

		isFloat := false
		for {
			char, err := r.ReadByte()
			if err != nil {
				if err == io.EOF {
					if json.Kind == Int {
						i, err := strconv.Atoi(string(raw.Bytes()))
						if err != nil {
							return nil, err
						}

						r.UnreadByte()
						json.Val = i
						return json, nil
					}

					if json.Kind == Float {
						i, err := strconv.ParseFloat(string(raw.Bytes()), 64)
						if err != nil {
							return nil, err
						}

						r.UnreadByte()
						json.Val = i
						return json, nil
					}
				}

				return nil, err
			}

			if isCharNumber(char) {
				raw.push(char)
				continue
			}

			if char == dot {
				if isFloat {
					// DELETE
					// fmt.Println("invalid num", string([]byte{char}))
					return nil, errInvalid
				}

				c, err := r.ReadByte()
				if err != nil {
					return nil, err
				}

				if !isCharNumber(c) {
					// DELETE
					// fmt.Println("invalid num", string([]byte{c}))
					return nil, errInvalid
				}

				isFloat = true
				r.UnreadByte()
				raw.push(char)
				json.Kind = Float
				continue
			}

			if json.Kind == Int {
				i, err := strconv.Atoi(string(raw.Bytes()))
				if err != nil {
					return nil, err
				}

				r.UnreadByte()
				json.Val = i
				return json, nil
			}

			if json.Kind == Float {
				i, err := strconv.ParseFloat(string(raw.Bytes()), 64)
				if err != nil {
					return nil, err
				}

				r.UnreadByte()
				json.Val = i
				return json, nil
			}
		}
	}

	// DELETE
	// fmt.Println("unmatch num", string([]byte{char}))
	return nil, errUnmatch
}

func parseNull(r reader) (*JSON, error) {
	char, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	// n
	if char == 110 {
		raw := new(Raw)
		raw.push(char)
		for i, c := range nullStr {
			if i == 0 {
				continue
			}

			char, err := r.ReadByte()
			if err != nil {
				return nil, err
			}

			if char != c {
				return nil, errInvalid
			}

			raw.push(char)
		}

		return &JSON{Kind: Null, Raw: raw, Val: nil}, nil
	}

	return nil, errUnmatch
}

func parseBool(r reader) (*JSON, error) {
	char, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	// DELETE
	// fmt.Println("bool", string([]byte{116}))

	// Try to parse true bool
	// t
	if char == 116 {
		raw := new(Raw)
		json := &JSON{Kind: Boolean, Raw: raw}
		raw.push(char)
		for i, c := range trueStr {
			if i == 0 {
				continue
			}

			char, err := r.ReadByte()
			if err != nil {
				return nil, err
			}

			if char != c {
				// DELETE
				// fmt.Println("invalid bool", string([]byte{char}))
				return nil, errInvalid
			}

			raw.push(char)
		}

		json.Val = true
		return json, nil
	}

	// Try to parse false bool
	// f
	if char == 102 {
		raw := new(Raw)
		json := &JSON{Kind: Boolean, Raw: raw}
		raw.push(char)

		for i, c := range falseStr {
			if i == 0 {
				continue
			}

			char, err := r.ReadByte()
			if err != nil {
				return nil, err
			}

			if char != c {
				return nil, errInvalid
			}

			raw.push(char)
		}

		json.Val = false
		return json, nil
	}

	// DELETE
	// fmt.Println("unmatch bool", string([]byte{char}))
	return nil, errUnmatch
}

func parseObj(r reader) (*JSON, error) {
	char, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	if char == openCurlBrack {
		// DELETE
		// fmt.Println("match obj", string([]byte{char}))
		raw := new(Raw)
		raw.push(char)
		json := &JSON{Kind: Object, Raw: raw, KeyVal: make(map[string]*JSON)}

		keyFound := false
		keyEnd := false
		valFound := false
		var currKey string
		for {
			char, err := r.ReadByte()
			if err != nil {
				return nil, err
			}

			if isCharSyntax(char) {
				continue
			} else {
				// DELETE
				// fmt.Println("invalid obj", string([]byte{char}))
			}

			if !keyFound && !keyEnd && !valFound {
				// Try to find key
				r.UnreadByte()
				key, err := parseStr(r)
				if err != nil {
					// DELETE
					// fmt.Println("invalid obj", string([]byte{char}))

					return nil, errInvalid
				}

				// DELETE
				// fmt.Println("key", string(key.Raw.Bytes()))
				keyFound = true
				raw.pushBytes(key.Raw.Bytes())
				currKey = key.Val.(string)
				continue
			}

			if keyFound && !keyEnd && !valFound {
				if char == colon {
					// DELETE
					// fmt.Println("key end", string([]byte{char}))

					keyEnd = true
					raw.push(char)
					continue
				}

				if isCharSyntax(char) {
					continue
				}

				// DELETE
				// fmt.Println("invalid key", string([]byte{char}))
				return nil, errInvalid
			}

			// Try to parse some value
			if keyFound && keyEnd && !valFound {
				if isCharSyntax(char) {
					continue
				}

				r.UnreadByte()

				// Try parse str
				if str, err := parseStr(r); err != nil {
					if err == errInvalid {
						return nil, errInvalid
					}

					if err == io.EOF {
						return nil, err
					}

					if err == errUnmatch {
						r.UnreadByte()
					}
				} else if err != nil {
					return nil, err
				} else {
					json.KeyVal[currKey] = str
					raw.pushBytes(str.Raw.Bytes())
					currKey = ""
					valFound = true
					continue
				}

				// Try to parse numeric
				if num, err := parseNum(r); err != nil {
					if err == errInvalid {
						return nil, errInvalid
					}

					if err == io.EOF {
						return nil, err
					}

					if err == errUnmatch {
						r.UnreadByte()
					}
				} else if err != nil {
					return nil, err
				} else {
					json.KeyVal[currKey] = num
					raw.pushBytes(num.Raw.Bytes())
					currKey = ""
					valFound = true
					continue
				}

				// Try to parse boolean
				if bl, err := parseBool(r); err != nil {
					if err == errInvalid {
						return nil, errInvalid
					}

					if err == io.EOF {
						return nil, err
					}

					if err == errUnmatch {
						r.UnreadByte()
					}
				} else if err != nil {
					return nil, err
				} else {
					json.KeyVal[currKey] = bl
					raw.pushBytes(bl.Raw.Bytes())
					currKey = ""
					valFound = true
					continue
				}

				// Try to parse null
				if nl, err := parseNull(r); err != nil {
					if err == errInvalid {
						return nil, errInvalid
					}

					if err == io.EOF {
						return nil, err
					}

					if err == errUnmatch {
						r.UnreadByte()
					}
				} else if err != nil {
					return nil, err
				} else {
					json.KeyVal[currKey] = nl
					raw.pushBytes(nl.Raw.Bytes())
					currKey = ""
					valFound = true
					continue
				}

				// Try to parse obj
				if obj, err := parseObj(r); err != nil {
					if err == errInvalid {
						return nil, errInvalid
					}

					if err == io.EOF {
						return nil, err
					}

					if err == errUnmatch {
						r.UnreadByte()
					}
				} else if err != nil {
					return nil, err
				} else {
					json.KeyVal[currKey] = obj
					raw.pushBytes(obj.Raw.Bytes())
					currKey = ""
					valFound = true
					continue
				}

				return nil, errInvalid
			}

			if keyFound && keyEnd && valFound {
				// Try to find delimiter
				if char == closeCurlBrack {
					raw.push(char)
					return json, nil
				}

				if char == coma {
					keyFound = false
					keyEnd = false
					valFound = false
					raw.push(char)
					continue
				}

				if isCharSyntax(char) {
					continue
				}

				return nil, errInvalid
			}
		}
	}

	// DELETE
	// fmt.Println("unmatch obj", string([]byte{char}))
	return nil, errUnmatch
}

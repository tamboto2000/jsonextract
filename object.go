package jsonextract

import (
	"io"
	"unicode"
)

func parseObj(r reader) (*JSON, error) {
	json := &JSON{Kind: Object}
	objMap := &objMap{val: make(map[string]*JSON)}
	json.push('{')

	// find first value
	if err := parseKeyVal(r, json, objMap); err != nil {
		return nil, err
	}

	onNext := false
	for {
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return nil, errInvalid
			}

			return nil, err
		}

		if unicode.IsControl(char) || char == ' ' {
			continue
		}

		if char == ',' {
			if onNext {
				return nil, errInvalid
			}

			onNext = true
			json.push(char)
			continue
		}

		if char == '}' {
			if onNext {
				return nil, errInvalid
			}

			json.push(char)
			break
		}

		if onNext {
			r.UnreadRune()
			if err := parseKeyVal(r, json, objMap); err != nil {
				return nil, err
			}

			onNext = false
		}
	}

	json.val = objMap.val
	return json, nil
}

func parseKeyVal(r reader, json *JSON, objMap *objMap) error {
	var id string

	// find key
	for {
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return errInvalid
			}

			return err
		}

		if unicode.IsControl(char) || char == ' ' {
			continue
		}

		if char != '"' {
			return errInvalid
		}

		str, err := parseString(r)
		if err != nil {
			return err
		}

		id = str.val.(string)
		break
	}

	// find key terminator
	for {
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return errInvalid
			}

			return err
		}

		if unicode.IsControl(char) || char == ' ' {
			continue
		}

		if char != ':' {
			return errInvalid
		}

		break
	}

	// find value
	for {
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return errInvalid
			}

			return err
		}

		if unicode.IsControl(char) || char == ' ' {
			continue
		}

		val, err := parse(r, char)
		if err != nil {
			return err
		}

		if val == nil {
			return errInvalid
		}

		objMap.val[id] = val
		json.pushRns([]rune(`"` + id + `":`))
		json.pushRns(val.raw)
		val.parent = json
		break
	}

	return nil
}

type objMap struct {
	val map[string]*JSON
}

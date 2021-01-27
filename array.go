package jsonextract

import (
	"io"
	"unicode"
)

func parseArray(r reader) (*JSON, error) {
	json := &JSON{kind: Array}
	json.push('[')

	vals := make([]*JSON, 0)
	// find first value
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

		if char == ']' {
			json.push(char)
			json.val = vals
			return json, nil
		}

		val, err := parse(r, char)
		if err != nil {
			return nil, err
		}

		if val == nil {
			return nil, errInvalid
		}

		vals = append(vals, val)
		json.pushRns(val.raw)
		val.parent = json
		break
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

		if char == ']' {
			if onNext {
				return nil, errInvalid
			}

			json.push(char)
			json.val = vals
			break
		}

		if char == ',' {
			if onNext {
				return nil, errInvalid
			}

			onNext = true
			json.push(char)
			continue
		}

		if onNext {
			val, err := parse(r, char)
			if err != nil {
				return nil, err
			}

			if val == nil {
				return nil, errInvalid
			}

			vals = append(vals, val)
			json.val = vals
			json.pushRns(val.raw)
			onNext = false
			val.parent = json
			continue
		}
	}

	return json, nil
}

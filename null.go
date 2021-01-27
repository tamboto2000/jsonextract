package jsonextract

import "io"

var nullVal = []rune{'u', 'l', 'l'}

func parseNull(r reader) (*JSON, error) {
	json := &JSON{kind: Null}
	json.push('n')
	for _, c := range nullVal {
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return nil, errInvalid
			}

			return nil, err
		}

		if char != c {
			return nil, errInvalid
		}

		json.push(char)
	}

	return json, nil
}

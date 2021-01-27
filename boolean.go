package jsonextract

import "io"

var (
	trueChars  = []rune{'r', 'u', 'e'}
	falseChars = []rune{'a', 'l', 's', 'e'}
)

func parseBool(r reader, firstC rune) (*JSON, error) {
	json := &JSON{kind: Boolean}
	json.push(firstC)

	// true bool
	if firstC == 't' {
		for _, c := range trueChars {
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

		json.val = true
		return json, nil
	}

	// false bool
	if firstC == 'f' {
		for _, c := range falseChars {
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

		json.val = false
		return json, nil
	}

	return json, nil
}

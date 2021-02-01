package jsonextract

import (
	jsonenc "encoding/json"
	"io"
	"unicode"
)

func parseString(r reader) (*JSON, error) {
	json := &JSON{kind: String}
	json.push('"')

	for {
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return nil, errInvalid
			}

			return nil, err
		}

		// detect escape
		if char == '\\' {
			json.push(char)
			char, _, err := r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return nil, errInvalid
				}

				return nil, err
			}

			if !isCharValidEscape(char) {
				return nil, errInvalid
			}

			json.push(char)

			// detect unicode
			if char == 'u' {
				for i := 0; i < 4; i++ {
					char, _, err := r.ReadRune()
					if err != nil {
						if err == io.EOF {
							return nil, errInvalid
						}

						return nil, err
					}

					if !isCharHex(char) {
						return nil, errInvalid
					}

					json.push(char)
				}
			}

			continue
		}

		// escape control character
		if unicode.IsControl(char) {
			json.pushRns(quoteRune(char))
			continue
		}

		// end string
		if char == '"' {
			json.push(char)
			break
		}

		json.push(char)
	}

	// unmarshal string
	str := new(string)
	if err := jsonenc.Unmarshal(json.Bytes(), str); err != nil {
		return nil, errInvalid
	}

	json.val = *str

	return json, nil
}

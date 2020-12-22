package jsonextract

import (
	jsonenc "encoding/json"
)

func parseStr(r reader) (*JSON, error) {
	char, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	if char == quot {
		raw := new(Raw)
		json := &JSON{Kind: String, Raw: raw}
		raw.push(char)

		for {
			char, err := r.ReadByte()
			if err != nil {
				return nil, err
			}

			if char == backSlash {
				raw.push(char)
				char, err := r.ReadByte()
				if err != nil {
					return nil, err
				}

				if !isCharEscapable(char) {
					return nil, errInvalid
				}

				raw.push(char)
				continue
			}

			if char == quot {
				raw.push(char)
				strVal := new(string)
				if err := jsonenc.Unmarshal(raw.byts, strVal); err != nil {
					return nil, err
				}

				json.Val = *strVal

				return json, nil
			}

			if isCharSyntax(char) {
				// if char is white space, just push
				if char == 32 {
					raw.push(char)
					continue
				}

				escChar := escapeChar(char)
				raw.pushBytes(escChar)
				continue
			}

			raw.push(char)
		}
	}

	return nil, errUnmatch
}

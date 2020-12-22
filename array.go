package jsonextract

func parseArr(r reader) (*JSON, error) {
	char, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	if char == openBrack {
		raw := new(Raw)
		json := &JSON{Kind: Array, Raw: raw}
		raw.push(char)

		for {
			val, err := parseArrayVal(r)
			if err != nil {
				return nil, err
			}

			raw.pushBytes(val.Raw.byts)
			json.Vals = append(json.Vals, val)

			char, err := r.ReadByte()
			if err != nil {
				return nil, err
			}

			if char == coma {
				raw.push(char)
				continue
			}

			if char == closeBrack {
				raw.push(char)
				return json, nil
			}
		}
	}

	return nil, errUnmatch
}

func parseArrayVal(r reader) (*JSON, error) {
	val, err := parse(r)
	if err != nil {
		return nil, err
	}

	// find delimiter (, or ])
	for {
		char, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		if isCharSyntax(char) {
			continue
		}

		if char == coma || char == closeBrack {
			r.UnreadByte()
			return val, nil
		}

		break
	}

	return nil, errInvalid
}

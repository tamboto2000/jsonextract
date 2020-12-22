package jsonextract

func parseObj(r reader) (*JSON, error) {
	char, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	if char == openCurlBrack {
		raw := new(Raw)
		json := &JSON{Kind: Object, Raw: raw, KeyVal: make(map[string]*JSON)}
		raw.push(char)

		for {
			key, val, err := parseObjKeyVal(r)
			if err != nil {
				return nil, err
			}

			json.KeyVal[key] = val
			raw.pushBytes([]byte(`"` + key + `":`))
			raw.pushBytes(val.Raw.byts)

			for {
				char, err := r.ReadByte()
				if err != nil {
					return nil, err
				}

				if char == coma {
					raw.push(char)
					break
				}

				if char == closeCurlBrack {
					raw.push(char)
					return json, nil
				}

				return nil, errInvalid
			}
		}
	}

	return nil, errUnmatch
}

func parseObjKeyVal(r reader) (string, *JSON, error) {
	keyStr, err := parseObjKey(r)
	if err != nil {
		return "", nil, err
	}

	val, err := parseObjVal(r)
	if err != nil {
		return "", nil, err
	}

	for {
		char, err := r.ReadByte()
		if err != nil {
			return "", nil, err
		}

		if char == coma || char == closeCurlBrack {
			r.UnreadByte()
			return keyStr, val, nil
		}

		if isCharSyntax(char) {
			continue
		}

		break
	}

	return "", nil, errInvalid
}

func parseObjKey(r reader) (string, error) {
	for {
		char, err := r.ReadByte()
		if err != nil {
			return "", err
		}

		if isCharSyntax(char) {
			continue
		}

		if char == quot {
			r.UnreadByte()
			val, err := parseStr(r)
			if err != nil {
				return "", err
			}

			keyStr := val.Val.(string)
			// find terminator (:)
			for {
				char, err := r.ReadByte()
				if err != nil {
					return "", err
				}

				if isCharSyntax(char) {
					continue
				}

				if char == colon {
					return keyStr, nil
				}

				break
			}
		}

		break
	}

	return "", errInvalid
}

func parseObjVal(r reader) (*JSON, error) {
	val, err := parse(r)
	if err != nil {
		return nil, err
	}

	for {
		char, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		if isCharSyntax(char) {
			continue
		}

		if char == coma || char == closeCurlBrack {
			r.UnreadByte()
			return val, nil
		}

		break
	}

	return nil, errInvalid
}

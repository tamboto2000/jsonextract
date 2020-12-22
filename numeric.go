package jsonextract

// pre-process numeric
func parseNum(r reader) (*JSON, error) {
	char, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	raw := new(Raw)
	json := &JSON{Kind: Int, Raw: raw}
	// if char is '-', verify if the char after is numeric
	if char == minus {
		raw.push(char)
		char, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		if !isCharNumber(char) {
			return nil, errInvalid
		}

		raw.push(char)
		// if char is num 0, verify if the char after is '.', like 0.3, 0.1, etc.,
		// the only valid numeric format with zero beginning is floating point
		if char == 48 {
			char, err := r.ReadByte()
			if err != nil {
				return nil, err
			}

			if char != dot {
				return nil, errInvalid
			}

			raw.push(char)
			json.Kind = Float
		}

		return parseNumVal(json, r)
	}

	if isCharNumber(char) {
		raw.push(char)

		// if char is num 0, verify if the char after is '.', like 0.3, 0.1, etc.,
		// the only valid numeric format with zero beginning is floating point
		if char == 48 {
			char, err := r.ReadByte()
			if err != nil {
				return nil, err
			}

			if char != dot {
				return nil, errInvalid
			}

			raw.push(char)
			json.Kind = Float
		}

		return parseNumVal(json, r)
	}

	return nil, errUnmatch
}

// parse numeric value
func parseNumVal(num *JSON, r reader) (*JSON, error) {
	for {
		char, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		if !isCharNumber(char) {
			if char == closeBrack || char == closeCurlBrack ||
				char == coma || isCharSyntax(char) {
				r.UnreadByte()
				return num, nil
			}

			if char == dot {
				if num.Kind == Float {
					return nil, errInvalid
				}

				num.Kind = Float
				num.Raw.push(char)
				continue
			}

			break
		}

		num.Raw.push(char)
	}

	return nil, errInvalid
}

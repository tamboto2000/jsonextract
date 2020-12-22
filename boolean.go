package jsonextract

func parseBool(r reader) (*JSON, error) {
	char, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	// t for true
	if char == 116 {
		raw := new(Raw)
		raw.push(char)
		json := &JSON{Kind: Boolean, Raw: raw}
		for i, c := range trueStr {
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

		json.Val = true
		return json, nil
	}

	// f for false
	if char == 102 {
		raw := new(Raw)
		raw.push(char)
		json := &JSON{Kind: Boolean, Raw: raw}
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

	return nil, errInvalid
}

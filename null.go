package jsonextract

func parseNull(r reader) (*JSON, error) {
	char, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	if char == 110 {
		raw := new(Raw)
		json := &JSON{Kind: Null, Raw: raw}
		raw.push(char)

		for i, c := range nullStr {
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

		json.Val = nil
		return json, nil
	}

	return nil, errUnmatch
}

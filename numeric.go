package jsonextract

import (
	"strconv"
)

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
		// the only valid numeric format with zero beginning is floating point, just zero, or exponent
		if char == 48 {
			char, err := r.ReadByte()
			if err != nil {
				return nil, err
			}

			// if char != dot {
			// 	return nil, errInvalid
			// }

			if char == dot {
				raw.push(char)
				char, err := r.ReadByte()
				if err != nil {
					return nil, err
				}

				if !isCharNumber(char) {
					return nil, errInvalid
				}

				r.UnreadByte()
				json.Kind = Float
			} else if isCharExponent(char) {
				raw.push(char)
				for i := 0; i < 2; i++ {
					char, err := r.ReadByte()
					if err != nil {
						return nil, err
					}

					if i == 0 {
						if char != minus && char != plus {
							return nil, errInvalid
						}

						raw.push(char)
					}

					if i == 1 {
						if !isCharNumber(char) {
							return nil, errInvalid
						}

						raw.push(char)
					}
				}

				json.WithExponent = true
			} else {
				return nil, errInvalid
			}
		}

		return parseNumVal(json, r)
	}

	if isCharNumber(char) {
		raw.push(char)

		// if char is num 0, verify if the char after is '.', like 0.3, 0.1, etc.,
		// the only valid numeric format with zero beginning is floating point, just zero, or exponent
		if char == 48 {
			char, err := r.ReadByte()
			if err != nil {
				return nil, err
			}

			if char == closeBrack || char == closeCurlBrack ||
				char == coma || isCharSyntax(char) {
				r.UnreadByte()
				return json, nil
			} else if char == dot {
				raw.push(char)
				char, err := r.ReadByte()
				if err != nil {
					return nil, err
				}

				if !isCharNumber(char) {
					return nil, errInvalid
				}

				r.UnreadByte()
				json.Kind = Float
			} else if isCharExponent(char) {
				raw.push(char)
				for i := 0; i < 2; i++ {
					char, err := r.ReadByte()
					if err != nil {
						return nil, err
					}

					if i == 0 {
						if char != minus && char != plus {
							return nil, errInvalid
						}

						raw.push(char)
					}

					if i == 1 {
						if !isCharNumber(char) {
							return nil, errInvalid
						}

						raw.push(char)
					}
				}

				json.WithExponent = true
			} else {
				return nil, errInvalid
			}
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
				if num.Kind == Int {
					if num.WithExponent {
						byts := convertExponentToParseable(num.Raw.byts)
						i, err := strconv.ParseFloat(string(byts), 64)
						if err != nil {
							return nil, err
						}

						num.Val = int(i)
					} else {
						i, err := strconv.Atoi(string(num.Raw.Bytes()))
						if err != nil {
							return nil, err
						}

						num.Val = i
					}
				}

				if num.Kind == Float {
					i, err := strconv.ParseFloat(string(num.Raw.Bytes()), 64)
					if err != nil {
						return nil, err
					}

					num.Val = i
				}

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

			if isCharExponent(char) {
				if num.WithExponent {
					return nil, errInvalid
				}

				num.Raw.push(char)
				for i := 0; i < 2; i++ {
					char, err := r.ReadByte()
					if err != nil {
						return nil, err
					}

					if i == 0 {
						if char != minus && char != plus {
							return nil, errInvalid
						}

						num.Raw.push(char)
					}

					if i == 1 {
						if !isCharNumber(char) {
							return nil, errInvalid
						}

						num.Raw.push(char)
					}
				}

				num.WithExponent = true
				continue
			}

			break
		}

		num.Raw.push(char)
	}

	return nil, errInvalid
}

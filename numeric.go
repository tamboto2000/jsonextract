package jsonextract

import (
	"io"
	"strconv"
	"unicode"
)

func parseNumeric(r reader, firstC rune) (*JSON, error) {
	json := &JSON{kind: Integer}

	// check if first char is minus sign
	if firstC == '-' {
		json.push(firstC)
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return nil, errInvalid
			}

			return nil, err
		}

		firstC = char
	}

	// if first char is not number, invalid
	if !unicode.IsNumber(firstC) {
		return nil, errInvalid
	}

	json.push(firstC)

	// check if first char is 0
	if firstC == '0' {
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				json.push(firstC)
				json.val = int64(0)
				return json, nil
			}

			return nil, err
		}

		// float indication
		if char == '.' {
			json.kind = Float
			json.push(char)
			char, _, err := r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return nil, errInvalid
				}

				return nil, err
			}

			if !unicode.IsNumber(char) {
				return nil, errInvalid
			}

			json.push(char)
			if err := parseNumFract(r, json); err != nil {
				return nil, err
			}

			return json, nil
		}

		// exponent indication
		if char == 'e' || char == 'E' {
			json.kind = Float
			json.push(char)
			if err := parseExp(r, json); err != nil {
				return nil, err
			}

			return json, nil
		}

		// if char end of number
		if isCharEndNum(char) {
			r.UnreadRune()

			json.val = int64(0)
			return json, nil
		}

		return nil, errInvalid
	}

	if err := parseNumFract(r, json); err != nil {
		return nil, err
	}

	return json, nil
}

func parseNumFract(r reader, json *JSON) error {
	for {
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		// detect float
		if char == '.' {
			if json.kind == Float {
				return errInvalid
			}

			json.push(char)
			json.kind = Float
			continue
		}

		// detect exponent
		if char == 'e' || char == 'E' {
			json.push(char)
			json.kind = Float
			if err := parseExp(r, json); err != nil {
				return err
			}

			break
		}

		if isCharEndNum(char) {
			r.UnreadRune()
			break
		}

		if !unicode.IsNumber(char) {
			return errInvalid
		}

		json.push(char)
	}

	if json.kind == Integer {
		i, err := strconv.ParseInt(string(json.raw), 10, 64)
		if err != nil {
			return errInvalid
		}

		json.val = i
	}

	if json.kind == Float {
		i, err := strconv.ParseFloat(string(json.raw), 64)
		if err != nil {
			return errInvalid
		}

		json.val = i
	}

	return nil
}

func parseExp(r reader, json *JSON) error {
	char, _, err := r.ReadRune()
	if err != nil {
		if err == io.EOF {
			return errInvalid
		}

		return err
	}

	json.push(char)
	if !unicode.IsNumber(char) {
		if isCharMinOrPlus(char) {
			char, _, err := r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return errInvalid
				}

				return err
			}

			if !unicode.IsNumber(char) {
				return errInvalid
			}

			json.push(char)
		} else {
			return errInvalid
		}
	}

	for {
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		if isCharEndNum(char) {
			r.UnreadRune()
			break
		}

		if !unicode.IsNumber(char) {
			return errInvalid
		}

		json.push(char)
	}

	i, err := strconv.ParseFloat(string(json.raw), 64)
	if err != nil {
		return errInvalid
	}

	json.val = i

	return nil
}

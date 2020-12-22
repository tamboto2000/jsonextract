package jsonextract

import (
	"strings"
)

func isCharLetter(char byte) bool {
	for _, c := range letters {
		if char == c {
			return true
		}
	}

	return false
}

func isCharNumber(char byte) bool {
	for _, c := range numbers {
		if char == c {
			return true
		}
	}

	return false
}

func isCharSyntax(char byte) bool {
	for _, c := range syntax {
		if char == c {
			return true
		}
	}

	return false
}

func isCharEscapable(char byte) bool {
	for _, c := range escapable {
		if char == c {
			return true
		}
	}

	return false
}

func isCharValidBeginObj(char byte) bool {
	if char == openBrack {
		return true
	}

	if char == openCurlBrack {
		return true
	}

	if char == quot {
		return true
	}

	if char == minus {
		return true
	}

	// t for true bool
	if char == 116 {
		return true
	}

	// f for false bool
	if char == 102 {
		return true
	}

	// n for null
	if char == 110 {
		return true
	}

	if isCharNumber(char) {
		return true
	}

	return false
}

func escapeChar(char byte) []byte {
	// \a
	if char == 7 {
		return []byte(`\a`)
	}

	// \b
	if char == 8 {
		return []byte(`\b`)
	}

	// \t
	if char == 9 {
		return []byte(`\t`)
	}

	// \n
	if char == 10 {
		return []byte(`\n`)
	}

	// \v
	if char == 11 {
		return []byte(`\v`)
	}

	// \f
	if char == 12 {
		return []byte(`\f`)
	}

	// \r
	if char == 13 {
		return []byte(`\r`)
	}

	return nil
}

func isCharExponent(char byte) bool {
	for _, c := range exponentChar {
		if char == c {
			return true
		}
	}

	return false
}

// apparently in json, something like 1e+1 is valid, but when parse is not,
// the equal valid form is 1.0e+1
func convertExponentToParseable(chars []byte) []byte {
	decimalFound := false
	for _, char := range chars {
		if char == dot {
			decimalFound = true
			return chars
		}

		if isCharExponent(char) && !decimalFound {
			chars = []byte(strings.Replace(string(chars), "e", ".0e", 1))
			return chars
		}
	}

	return chars
}

func isCharHex(char byte) bool {
	for _, c := range hexChars {
		if char == c {
			return true
		}
	}

	return false
}

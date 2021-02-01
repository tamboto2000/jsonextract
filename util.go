package jsonextract

import (
	"bytes"
	"io"
	"strconv"
	"unicode"
	"unicode/utf8"
)

// Convert []rune to []byte
func runesToUTF8(rs []rune) []byte {
	size := 0
	for _, r := range rs {
		size += utf8.RuneLen(r)
	}

	bs := make([]byte, size)

	count := 0
	for _, r := range rs {
		count += utf8.EncodeRune(bs[count:], r)
	}

	return bs
}

var validEsc = []rune{'"', '\\', '/', 'b', 'f', 'n', 'r', 't', 'r', 'u'}

// Check if char is valid escape
func isCharValidEscape(char rune) bool {
	for _, c := range validEsc {
		if char == c {
			return true
		}
	}

	return false
}

var hex = []rune{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'a', 'b', 'c', 'd', 'e', 'f', 'A', 'B', 'C', 'D', 'E', 'F'}

// Check if a rune is a hexa char
func isCharHex(char rune) bool {
	for _, c := range hex {
		if char == c {
			return true
		}
	}

	return false
}

// quote rune to string converted to []rune
func quoteRune(char rune) []rune {
	rns := []rune(strconv.QuoteRune(char))
	return rns[1 : len(rns)-1]
}

var endNum = []rune{'}', ']', ',', ' '}

// check if char numeric ending
func isCharEndNum(char rune) bool {
	for _, c := range endNum {
		if char == c {
			return true
		}
	}

	if unicode.IsControl(char) || char == ' ' {
		return true
	}

	return false
}

// check if exponent valid
func isExpValid(r reader) (rune, error) {
	char, _, err := r.ReadRune()
	if err != nil {
		if err == io.EOF {
			return 0, errInvalid
		}

		return 0, err
	}

	if !unicode.IsNumber(char) && char != '+' && char != '-' {
		return 0, errInvalid
	}

	return char, nil
}

func isCharMinOrPlus(char rune) bool {
	if char == '-' || char == '+' {
		return true
	}

	return false
}

func getParent(json *JSON) *JSON {
	for json.parent != nil {
		json = json.parent
	}

	return json
}

// read runes from bytes.Reader
func readAllRunes(r *bytes.Reader) []rune {
	size := r.Size()
	rns := make([]rune, size)

	for i := int64(0); i < size; i++ {
		r, _, err := r.ReadRune()
		if err != nil {
			break
		}

		rns[i] = r
	}

	return rns
}

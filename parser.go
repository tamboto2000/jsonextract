package jsonextract

import (
	"bytes"
	"io"
)

var (
	letters        = []byte{97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90}
	numbers        = []byte{49, 50, 51, 52, 53, 54, 55, 56, 57, 48}
	symbols        = []byte{126, 96, 33, 64, 35, 36, 37, 94, 38, 42, 40, 41, 95, 45, 43, 61, 123, 91, 125, 93, 124, 92, 58, 59, 34, 39, 60, 44, 62, 46, 63, 47}
	syntax         = []byte{7, 8, 9, 10, 11, 12, 13, 32}
	openCurlBrack  = byte(123)
	closeCurlBrack = byte(125)
	openBrack      = byte(91)
	closeBrack     = byte(93)
	quot           = byte(34)
	colon          = byte(58)
	coma           = byte(44)
)

func parse(data io.Reader) ([][]byte, error) {
	parsed := make([][]byte, 0)
	reader := newReader(data)
	for reader.next() {
		char := reader.get()
		if char == openCurlBrack {
			node := newObjNode(reader)
			if node.parseObj() {
				parsed = append(parsed, node.result())
			} else {
				reader.unread()
				node.flush()
			}

			continue
		}

		if char == openBrack {
			node := newArrNode(reader)
			if node.parseArr() {
				parsed = append(parsed, node.result())
			} else {
				reader.unread()
				node.flush()
			}

			continue
		}
	}

	if err := reader.err; err != io.EOF {
		return nil, err
	}

	return parsed, nil
}

type node struct {
	buff   *bytes.Buffer
	reader *reader
	char   byte
}

func newObjNode(r *reader) *node {
	n := &node{
		buff:   bytes.NewBuffer(make([]byte, 0)),
		reader: r,
	}

	n.buff.WriteByte(openCurlBrack)

	return n
}

func newArrNode(r *reader) *node {
	n := &node{
		buff:   bytes.NewBuffer(make([]byte, 0)),
		reader: r,
	}

	n.buff.WriteByte(openBrack)

	return n
}

func (n *node) next() bool {
	if !n.reader.next() {
		return false
	}

	n.char = n.reader.get()

	return true
}

func (n *node) result() []byte {
	byts := n.buff.Bytes()
	n.flush()

	return byts
}

func (n *node) flush() {
	n.buff.Reset()
	n.buff = nil
}

func (n *node) pushChar() {
	n.buff.WriteByte(n.char)
}

// parse json object key
func (n *node) parseKey() bool {
	endKey := false
	for n.next() {
		if !endKey {
			if n.char == quot {
				endKey = true
			}

			n.pushChar()
			continue
		}

		if endKey {
			if n.char == colon {
				n.pushChar()
				return true
			}

			if isCharSyntax(n.char) {
				n.pushChar()
				continue
			}

			return false
		}
	}

	return false
}

// parse json object
func (n *node) parseObj() bool {
	nextVal := false
	keyFound := false
	valFound := false
	for n.next() {
		if !nextVal && !keyFound && !valFound {
			if n.char == quot {
				n.pushChar()
				if !n.parseKey() {
					return false
				}

				keyFound = true
				continue
			}

			if n.char == closeCurlBrack {
				n.pushChar()
				return true
			}

			if isCharSyntax(n.char) {
				n.pushChar()
				continue
			}

			return false
		}

		if !nextVal && keyFound && !valFound {
			// try parse val
			n.reader.unread()
			if n.parseVal() {
				valFound = true
				continue
			}

			return false
		}

		if !nextVal && keyFound && valFound {
			if n.char == coma {
				n.pushChar()
				keyFound = false
				valFound = false
				nextVal = true
				continue
			}

			if n.char == closeCurlBrack {
				n.pushChar()
				return true
			}

			if isCharSyntax(n.char) {
				n.pushChar()
				continue
			}

			return false
		}

		if nextVal && !keyFound && !valFound {
			if n.char == quot {
				n.pushChar()
				if !n.parseKey() {
					return false
				}

				keyFound = true
				continue
			}

			if isCharSyntax(n.char) {
				n.pushChar()
				continue
			}

			return false
		}

		if nextVal && keyFound && !valFound {
			// try to parse val
			n.reader.unread()
			if n.parseVal() {
				valFound = true
				continue
			}

			return false
		}

		if nextVal && keyFound && valFound {
			if n.char == coma {
				n.pushChar()
				keyFound = false
				valFound = false
				continue
			}

			if n.char == closeCurlBrack {
				n.pushChar()
				return true
			}

			if isCharSyntax(n.char) {
				n.pushChar()
				continue
			}

			return false
		}
	}

	return false
}

func (n *node) parseArr() bool {
	nextVal := false
	foundVal := false
	for n.next() {
		if !nextVal && !foundVal {
			if n.char == closeBrack {
				n.pushChar()
				return true
			}

			if isCharSyntax(n.char) {
				n.pushChar()
				continue
			}

			n.reader.unread()
			if n.parseVal() {
				foundVal = true
				continue
			}

			return false
		}

		if !nextVal && foundVal {
			if n.char == closeBrack {
				n.pushChar()
				return true
			}

			if n.char == coma {
				n.pushChar()
				nextVal = true
				foundVal = false
				continue
			}

			if isCharSyntax(n.char) {
				n.pushChar()
				continue
			}

			return false
		}

		if nextVal && !foundVal {
			if isCharSyntax(n.char) {
				n.pushChar()
				continue
			}

			n.reader.unread()
			if n.parseVal() {
				foundVal = true
				continue
			}

			return false
		}

		if nextVal && foundVal {
			if n.char == closeBrack {
				n.pushChar()
				return true
			}

			if n.char == coma {
				n.pushChar()
				foundVal = false
				continue
			}

			if isCharSyntax(n.char) {
				n.pushChar()
				continue
			}

			return false
		}
	}

	return false
}

func (n *node) parseStrVal() bool {
	var prevc byte
	for n.next() {
		if n.char == quot {
			n.pushChar()

			// back slash (\)
			if prevc == 92 {
				prevc = n.char
				continue
			}

			return true
		}

		// char /
		if n.char == 47 {
			if prevc == 92 {
				n.pushChar()
				prevc = n.char
				continue
			}

			return false
		}

		n.pushChar()
		prevc = n.char
	}

	return false
}

func (n *node) parseNumVal() bool {
	for n.next() {
		if isCharNumber(n.char) {
			n.pushChar()
			continue
		}

		// dot (.)
		if n.char == 46 {
			n.pushChar()
			if n.parseFloatVal() {
				return true
			}

			return false
		}

		if n.char == coma || n.char == closeBrack || n.char == closeCurlBrack {
			n.reader.unread()
			return true
		}

		if isCharSyntax(n.char) {
			n.pushChar()
			continue
		}

		return false
	}

	return false
}

func (n *node) parseFloatVal() bool {
	for n.next() {
		if isCharNumber(n.char) {
			n.pushChar()
			continue
		}

		if n.char == coma || n.char == closeBrack || n.char == closeCurlBrack {
			n.reader.unread()
			return true
		}

		if isCharSyntax(n.char) {
			n.pushChar()
			continue
		}

		return false
	}

	return false
}

func (n *node) parseTrueBool() bool {
	expect := byte(114)
	for n.next() {
		if n.char != expect {
			return false
		}

		n.pushChar()

		if expect == 114 {
			expect = 117
			continue
		}

		if expect == 117 {
			expect = 101
			continue
		}

		if expect == 101 {
			return true
		}
	}

	return false
}

func (n *node) parseFalseBool() bool {
	expect := byte(97)
	for n.next() {
		if n.char != expect {
			return false
		}

		n.pushChar()

		if expect == 97 {
			expect = 108
			continue
		}

		if expect == 108 {
			expect = 115
			continue
		}

		if expect == 115 {
			expect = 101
			continue
		}

		if expect == 101 {
			return true
		}
	}

	return false
}

func (n *node) parseVal() bool {
	for n.next() {
		// try to parse string
		if n.char == quot {
			n.pushChar()
			if n.parseStrVal() {
				return true
			}

			return false
		}

		// try to parse numeric
		if isCharNumber(n.char) {
			n.pushChar()
			if n.parseNumVal() {
				return true
			}

			return false
		}

		// try to parse boolean
		if isCharLetter(n.char) {
			n.pushChar()
			if n.char == 116 {
				if n.parseTrueBool() {
					return true
				}

				return false
			}

			if n.char == 102 {
				if n.parseFalseBool() {
					return true
				}

				return false
			}

			return false
		}

		// try to parse object
		if n.char == openCurlBrack {
			n.pushChar()
			if n.parseObj() {
				return true
			}

			return false
		}

		// try to parse array
		if n.char == openBrack {
			n.pushChar()
			if n.parseArr() {
				return true
			}

			return false
		}

		if isCharSyntax(n.char) {
			n.pushChar()
			continue
		}

		return false
	}

	return false
}

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

func isCharSymbol(char byte) bool {
	for _, c := range symbols {
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

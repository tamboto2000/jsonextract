// Package jsonextract is a small library for extracting JSON from a string, it extract a possible valid JSONs from a string or text
package jsonextract

import (
	"bytes"
	"io"
	"os"
)

var (
	// {
	openCurlBracket = byte(123)
	// }
	closeCurlBracket = byte(125)
	// [
	openBracket = byte(91)
	// ]
	closeBracket = byte(93)
	// "
	doubleTick = byte(34)
	// a-z+A-Z
	letters = []byte{97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90}
	// 0-9
	numbers = []byte{49, 50, 51, 52, 53, 54, 55, 56, 57, 48}
	// symbols
	symbols = []byte{126, 33, 64, 35, 36, 37, 94, 38, 42, 40, 41, 95, 45, 43, 61, 124, 92, 58, 59, 39, 60, 44, 62, 46, 63, 47}
)

// JSONFromStr extract every possible JSONs from a string
func JSONFromStr(data string) ([][]byte, error) {
	byts := []byte(data)
	return JSONFromBytes(byts)
}

// JSONFromBytes extract every possible JSONs in an array of bytes
func JSONFromBytes(str []byte) ([][]byte, error) {
	reader := bytes.NewReader(str)
	r := newReader(reader)
	ex := newExtractor(r)
	ex.run()

	return ex.result(), ex.err()
}

// JSONFromReader extract every possible JSONs in a io.Reader
func JSONFromReader(reader io.Reader) ([][]byte, error) {
	r := newReader(reader)
	ex := newExtractor(r)
	ex.run()
	return ex.result(), ex.err()
}

// JSONFromFile extract every possible JSONs in a file in path
func JSONFromFile(path string) ([][]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return JSONFromReader(f)
}

type extractor struct {
	rest  [][]byte
	temp  []byte
	r     *reader
	prevc byte
	ocb   int
	ccb   int
	ob    int
	cb    int
	er    error
	isObj bool
	isArr bool
}

func newExtractor(r *reader) *extractor {
	return &extractor{r: r}
}

func (ex *extractor) run() {
	for ex.r.next() {
		char := ex.r.get()
		if ex.ocb == 0 && char == openCurlBracket {
			if !ex.isArr {
				ex.prevc = char
				ex.temp = append(ex.temp, char)
				ex.isObj = true
				ex.ocb++
				continue
			}
		}

		if ex.ob == 0 && char == openBracket {
			if !ex.isObj {
				ex.prevc = char
				ex.temp = append(ex.temp, char)
				ex.isArr = true
				ex.ob++
				continue
			}
		}

		if ex.isObj {
			if ex.prevc == openCurlBracket {
				valid := isValidCharAfterOpenCurlBracket(char)
				if valid == 2 {
					for valid == 2 {
						if ex.r.next() {
							char = ex.r.get()
						} else {
							break
						}

						valid = isValidCharAfterOpenCurlBracket(char)
					}
				}

				if valid == 0 {
					ex.reset()
					continue
				}
			}

			if ex.prevc == openBracket {
				valid := isValidCharAfterOpenBracket(char)
				if valid == 2 {
					for valid == 2 {
						if ex.r.next() {
							char = ex.r.get()
						} else {
							break
						}

						valid = isValidCharAfterOpenBracket(char)
					}
				}

				if valid == 0 {
					ex.reset()
					continue
				}
			}

			if char == openCurlBracket {
				ex.ocb++
			}

			if char == closeCurlBracket {
				ex.ccb++
			}

			if char == openBracket {
				ex.ob++
			}

			if char == closeBracket {
				ex.cb++
			}

			ex.prevc = char
			ex.temp = append(ex.temp, char)

			if ex.ocb == ex.ccb && ex.ob == ex.cb {
				ex.merge()
				continue
			}
		}

		if ex.isArr {
			if ex.prevc == openBracket {
				valid := isValidCharAfterOpenBracket(char)
				if valid == 2 {
					for valid == 2 {
						if ex.r.next() {
							char = ex.r.get()
						} else {
							break
						}

						valid = isValidCharAfterOpenBracket(char)
					}
				}

				if valid == 0 {
					ex.reset()
					continue
				}
			}

			if ex.prevc == openCurlBracket {
				valid := isValidCharAfterOpenCurlBracket(char)
				if valid == 2 {
					for valid == 2 {
						if ex.r.next() {
							char = ex.r.get()
						} else {
							break
						}

						valid = isValidCharAfterOpenCurlBracket(char)
					}
				}

				if valid == 0 {
					ex.reset()
					continue
				}
			}

			if char == openBracket {
				ex.ob++
			}

			if char == closeBracket {
				ex.cb++
			}

			if char == openCurlBracket {
				ex.ocb++
			}

			if char == closeCurlBracket {
				ex.ccb++
			}

			ex.prevc = char
			ex.temp = append(ex.temp, char)

			if ex.ob == ex.cb && ex.ocb == ex.ccb {
				ex.merge()
				continue
			}
		}
	}
}

func (ex *extractor) reset() {
	ex.temp = make([]byte, 0)
	ex.prevc = 0
	ex.ocb = 0
	ex.ccb = 0
	ex.ob = 0
	ex.cb = 0
	ex.isObj = false
	ex.isArr = false
}

func (ex *extractor) merge() {
	ex.rest = append(ex.rest, ex.temp)
	ex.reset()
}

func (ex *extractor) err() error {
	return ex.er
}

func (ex *extractor) result() [][]byte {
	return ex.rest
}

type reader struct {
	r   io.Reader
	byt []byte
	err error
}

func newReader(r io.Reader) *reader {
	return &reader{
		byt: make([]byte, 1),
		r:   r,
	}
}

func (r *reader) next() bool {
	_, err := r.r.Read(r.byt)
	if err != nil {
		if err != io.EOF {
			r.err = err
		}

		return false
	}

	return true
}

func (r *reader) get() byte {
	return r.byt[0]
}

// check if char after open curl bracket is valid.
// return 0 if not valid, return 1 if valid, return 2 if not sure
func isValidCharAfterOpenCurlBracket(char byte) int {
	if char == doubleTick {
		return 1
	}

	if char == openCurlBracket {
		return 1
	}

	if char == openBracket {
		return 1
	}

	if char == closeCurlBracket {
		return 1
	}

	for _, c := range letters {
		if char == c {
			return 0
		}
	}

	for _, c := range numbers {
		if char == c {
			return 0
		}
	}

	for _, c := range symbols {
		if char == c {
			return 0
		}
	}

	return 2
}

// check if char after open bracket is valid.
// return 0 if not valid, return 1 if valid, return 2 if not sure
func isValidCharAfterOpenBracket(char byte) int {
	if char == doubleTick {
		return 1
	}

	if char == openCurlBracket {
		return 1
	}

	if char == openBracket {
		return 1
	}

	if char == closeBracket {
		return 1
	}

	for _, c := range letters {
		if char == c {
			return 0
		}
	}

	for _, c := range numbers {
		if char == c {
			return 1
		}
	}

	for _, c := range symbols {
		if char == c {
			return 0
		}
	}

	return 2
}

package jsonextract

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type reader interface {
	ReadRune() (r rune, size int, err error)
	UnreadRune() error
}

func readFromBytes(byts []byte) reader {
	buff := bytes.NewReader(byts)
	return buff
}

func readFromString(str string) reader {
	rdr := strings.NewReader(str)
	return rdr
}

func readFromReader(r io.Reader) (reader, error) {
	buff := bufio.NewReader(r)

	return buff, nil
}

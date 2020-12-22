package jsonextract

import (
	"bytes"
	"io"
	"strings"
)

type reader interface {
	ReadByte() (byte, error)
	UnreadByte() error
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
	init := make([]byte, 0)
	buff := bytes.NewBuffer(init)
	if _, err := buff.ReadFrom(r); err != nil {
		return nil, err
	}

	return buff, nil
}

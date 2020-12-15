package jsonextract

import "io"

type reader struct {
	r        io.Reader
	byt      []byte
	err      error
	isUnread bool
}

// bytes reader
func newReader(r io.Reader) *reader {
	return &reader{
		byt: make([]byte, 1),
		r:   r,
	}
}

func (r *reader) next() bool {
	if r.isUnread {
		return true
	}

	_, err := r.r.Read(r.byt)
	if err != nil {
		r.err = err

		return false
	}

	return true
}

func (r *reader) get() byte {
	if r.isUnread {
		r.isUnread = false
	}

	return r.byt[0]
}

func (r *reader) unread() {
	r.isUnread = true
}

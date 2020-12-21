package jsonextract

import "errors"

// occured when data type already recognized, but its invalid form the recognized type.
// Example, when the first char is numeric or '-' char, it is recognized that the
// data type should be numeric, but if after the first char is not numeric, nor
// JSON delimiter (',', '}' or ']'), nor syntax char, then it must be invalid
var errInvalid = errors.New("invalid")

// occured when data type can't be determined.
// Example, when the first char is a syntax char, like \n, \t, white space. etc. then it
// can't be determined since the next char is not recognizable type
// var errUnrecognizable = errors.New("unrecognizable")

// occured when the called data parser method check the first char and it's not
// the rigth parser to parse the data type.
// Example, when parser method for parsing numeric type called and the first char is
// not numeric, then it's not the parser job to parse the data
var errUnmatch = errors.New("unmatch")

// occured when reader is EOF
// var errReaderEOF = errors.New("EOF")

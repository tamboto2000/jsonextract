package jsonextract

import "errors"

// occured when data type already recognized, but its invalid form the recognized type.
// Example, when the first char is numeric or '-' char, it is recognized that the
// data type should be numeric, but if after the first char is not numeric, nor
// JSON delimiter (',', '}' or ']'), nor syntax char, then it must be invalid
var errInvalid = errors.New("invalid")

package jsonextract

import "encoding/json"

// Valid returns which JSONs is valid and which not
func Valid(data [][]byte) (valid [][]byte, invalid [][]byte) {
	for _, d := range data {
		if json.Valid(d) {
			valid = append(valid, d)
		} else {
			invalid = append(invalid, d)
		}
	}

	return valid, invalid
}

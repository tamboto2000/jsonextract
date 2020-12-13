package jsonextract

import (
	"os"
)

// SaveToPath save extracted JSONs to a file path
func SaveToPath(data [][]byte, path string) error {
	return save(data, path)
}

// Save save extracted JSONs to ./extracted_jsons.json
func Save(data [][]byte) error {
	return save(data, "./extracted_jsons.json")
}

func save(data [][]byte, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	rest := make([]byte, 0)
	rest = append(rest, 91)
	c := len(data)
	for i, d := range data {
		rest = append(rest, d...)
		if i == c-1 {
			rest = append(rest, 93)
		} else {
			rest = append(rest, 44)
		}
	}

	if _, err = f.Write(rest); err != nil {
		return err
	}

	return nil
}

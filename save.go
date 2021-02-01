package jsonextract

import "os"

// SaveToPath save extracted JSONs to a file path
func SaveToPath(data []*JSON, path string) error {
	return save(data, path)
}

// Save save extracted JSONs to ./extracted_jsons.json
func Save(data []*JSON) error {
	return save(data, "./extracted_jsons.json")
}

func save(data []*JSON, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	rest := make([]byte, 0)
	rest = append(rest, 91)
	c := len(data)
	if c > 0 {
		for i, d := range data {
			rest = append(rest, d.Bytes()...)
			if i == c-1 {
				rest = append(rest, 93)
			} else {
				rest = append(rest, 44)
			}
		}
	} else {
		rest = append(rest, 93)
	}

	if _, err = f.Write(rest); err != nil {
		return err
	}

	return nil
}

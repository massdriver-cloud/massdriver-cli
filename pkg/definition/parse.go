package definition

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func Parse(path string) (*Definition, error) {
	def := new(Definition)

	defFile, err := os.Open(path)
	if err != nil {
		return def, err
	}
	defer defFile.Close()

	byteValue, _ := ioutil.ReadAll(defFile)
	if jsonErr := json.Unmarshal(byteValue, &def); err != nil {
		return def, jsonErr
	}

	return def, nil
}

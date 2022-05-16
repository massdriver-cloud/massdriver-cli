package bundle

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// ParseBundle parses a bundle from a YAML file
func Parse(path string) (*Bundle, error) {
	bundle := new(Bundle)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal([]byte(data), &bundle)
	if err != nil {
		return nil, err
	}

	return bundle, nil
}

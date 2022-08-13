package bundle

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Overrides struct {
	Access string
}

// ParseBundle parses a bundle from a YAML file
// overrides allow the CLI to override specific bundle metadata.
// This is useful in a CI/CD scenario when you want to change the `access` if you are deploying to a sandbox org.
func Parse(path string, overrides map[string]interface{}) (*Bundle, error) {
	bundle := new(Bundle)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &bundle)
	if err != nil {
		return nil, err
	}

	if overrideErr := applyOverrides(bundle, overrides); overrideErr != nil {
		return nil, overrideErr
	}

	return bundle, nil
}

func applyOverrides(b *Bundle, overrides map[string]interface{}) error {
	if access, found := overrides["access"]; found {
		if access == "public" || access == "private" {
			var ok bool
			b.Access, ok = access.(string)
			if !ok {
				return errors.New("invalid access override")
			}
		}
	}
	return nil
}

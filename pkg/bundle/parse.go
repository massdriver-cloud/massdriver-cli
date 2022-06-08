package bundle

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type BundleOverrides struct {
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

	err = yaml.Unmarshal([]byte(data), &bundle)
	if err != nil {
		return nil, err
	}

	applyOverrides(bundle, overrides)

	return bundle, nil
}

func applyOverrides(b *Bundle, overrides map[string]interface{}) {
	if access, found := overrides["access"]; found {
		// TODO: we need to add a meta schema for our metadata and validate the bundle
		if access == "public" || access == "private" {
			b.Access = access.(string)
		}
	}
}

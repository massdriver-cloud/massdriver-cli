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

	setDefaultSteps(bundle)
	applyOverrides(bundle, overrides)

	return bundle, nil
}


// Sets the default steps to be a single src dir for terraform
func setDefaultSteps(bundle *Bundle) {
	if len(bundle.Steps) == 0 {
		defaultStep := BundleStep{Path: "src", Provisioner: "terraform"}
		bundle.Steps = []BundleStep{defaultStep}
  }
}

func applyOverrides(b *Bundle, overrides map[string]interface{}) {
	if access, found := overrides["access"]; found {
		if access == "public" || access == "private" {
			b.Access = access.(string)
		}
	}
}

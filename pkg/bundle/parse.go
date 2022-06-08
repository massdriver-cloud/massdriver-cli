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

	setDefaultSteps(bundle)

	return bundle, nil
}

// Sets the default steps to be a single src dir for terraform
func setDefaultSteps(bundle *Bundle) {
	if len(bundle.Steps) == 0 {
		defaultStep := BundleStep{Path: "src", Provisioner: "terraform"}
		bundle.Steps = []BundleStep{defaultStep}
	}
}

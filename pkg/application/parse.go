package application

import (
	"errors"
	"io/ioutil"
	"log"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"gopkg.in/yaml.v3"
)

// TODO: combine with bundle.Parse
func Parse(path string, overrides map[string]interface{}) (*bundle.Bundle, error) {
	app := new(bundle.Bundle)

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, app)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	if overrideErr := applyOverrides(app, overrides); overrideErr != nil {
		return nil, overrideErr
	}

	return app, nil
}

func applyOverrides(a *bundle.Bundle, overrides map[string]interface{}) error {
	if access, found := overrides["access"]; found {
		if access == "public" || access == "private" {
			var ok bool
			a.Access, ok = access.(string)
			if !ok {
				return errors.New("invalid access override")
			}
		}
	}
	return nil
}

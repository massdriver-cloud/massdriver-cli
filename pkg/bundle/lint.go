package bundle

import (
	"embed"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

//go:embed schemas/bundle-schema.json
var bundleFS embed.FS

func Lint(b *Bundle) error {
	err := LintSchema(b)
	if err != nil {
		return err
	}

	err = LintParamsConnectionsNameCollision(b)
	if err != nil {
		return err
	}

	return nil
}

func LintSchema(b *Bundle) error {
	schemaBytes, _ := bundleFS.ReadFile("schemas/bundle-schema.json")
	documentLoader := gojsonschema.NewGoLoader(b)
	schemaLoader := gojsonschema.NewBytesLoader(schemaBytes)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		errors := "massdriver.yaml has schema violations:\n"
		for _, violation := range result.Errors() {
			errors += fmt.Sprintf("\t- %v\n", violation)
		}
		return fmt.Errorf(errors)
	}
	return nil
}

func LintParamsConnectionsNameCollision(b *Bundle) error {
	if b.Params != nil {
		if params, ok := b.Params["properties"]; ok {
			if b.Connections != nil {
				if connections, ok := b.Connections["properties"]; ok {
					for param := range params.(map[string]interface{}) {
						for connection := range connections.(map[string]interface{}) {
							if param == connection {
								return fmt.Errorf("a parameter and connection have the same name: %s", param)
							}
						}
					}
				}
			}
		}
	}
	return nil
}

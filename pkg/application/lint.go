package application

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

func Lint(file string, schema string) (bool, error) {
	schemaLoader := gojsonschema.NewReferenceLoader("file://testdata/schema-application.json")
	documentLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", file))

	return true, nil

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return false, err
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}

	return result.Valid(), nil
}

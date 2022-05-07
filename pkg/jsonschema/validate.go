package jsonschema

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/xeipuuv/gojsonschema"
)

// Validate the input object against the schema
func Validate(schemaPath string, documentPath string) (bool, error) {
	log.Debug().
		Str("schemaPath", schemaPath).
		Str("documentPath", documentPath).Msg("Validating schema.")

	sl := Loader(schemaPath)
	dl := Loader(documentPath)

	result, err := gojsonschema.Validate(sl, dl)
	if err != nil {
		log.Error().Err(err).Msg("Validator failed.")
		return false, err
	}

	if !result.Valid() {
		msg := fmt.Sprintf("The document failed validation:\n\tDocument: %s\n\tSchema: %s\nErrors:\n", documentPath, schemaPath)
		for _, desc := range result.Errors() {
			msg = msg + fmt.Sprintf("\t- %s\n", desc)
		}

		err = errors.New(msg)
		log.Error().Err(err).Msg("Validation failed.")
		return false, err
	}

	return true, nil
}

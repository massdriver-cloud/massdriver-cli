package jsonschema

import (
	"github.com/rs/zerolog/log"
	"github.com/xeipuuv/gojsonschema"
)

// Validate the input object against the schema
func Validate(schemaPath string, documentPath string) (*gojsonschema.Result, error) {
	log.Debug().
		Str("schemaPath", schemaPath).
		Str("documentPath", documentPath).Msg("Validating schema.")

	sl := Loader(schemaPath)
	dl := Loader(documentPath)

	return gojsonschema.Validate(sl, dl)
}

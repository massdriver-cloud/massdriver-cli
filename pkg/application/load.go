package application

import (
	"regexp"

	"github.com/xeipuuv/gojsonschema"
)

const filePrefix = "file://"

var loaderPrefixPattern = regexp.MustCompile(`^(file|http|https)://`)

// Load a JSON Schema with or without a path prefix
func Load(path string) gojsonschema.JSONLoader {
	var ref string
	if loaderPrefixPattern.MatchString(path) {
		ref = path
	} else {
		ref = filePrefix + path
	}

	return gojsonschema.NewReferenceLoader(ref)
}

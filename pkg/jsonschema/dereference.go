package jsonschema

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
)

var schemaTypePattern = regexp.MustCompile(`^.*\/(.*).json$`)

// A RefdSchema is a JSON Schema that may contain $ref
type RefdSchema struct {
	SchemaID   string
	Definition interface{}
}

func WriteDereferencedSchema(schemaFilePath string, outFile io.Writer, c *client.MassdriverClient) error {
	dereferencedSchema := RefdSchema{}
	rawJSON := map[string]interface{}{}
	cwd := filepath.Dir(schemaFilePath)
	data, readErr := ioutil.ReadFile(schemaFilePath)
	if readErr != nil {
		return readErr
	}

	if err := json.Unmarshal(data, &rawJSON); err != nil {
		return err
	}
	definition, err := Hydrate(rawJSON, cwd, c)
	if err != nil {
		return err
	}
	dereferencedSchema.Definition = definition

	json, err := json.Marshal(dereferencedSchema.Definition)
	if err != nil {
		return err
	}

	_, err = outFile.Write(append(json, []byte("\n")...))

	return err
}

func (r *RefdSchema) Type() string {
	bytes := []byte(r.SchemaID)
	match := schemaTypePattern.FindSubmatch(bytes)[1]
	artifactType := string(match)

	return artifactType
}

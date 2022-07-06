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
	SchemaId   string
	Definition interface{}
}

func WriteDereferencedSchema(schemaFilePath string, outFile io.Writer, c *client.MassdriverClient) error {
	dereferencedSchema := RefdSchema{}
	rawJson := map[string]interface{}{}
	cwd := filepath.Dir(schemaFilePath)
	data, err := ioutil.ReadFile(schemaFilePath)
	if err != nil {
		return err
	}

	json.Unmarshal(data, &rawJson)
	definition, err := Hydrate(rawJson, cwd, c)
	if err != nil {
		return err
	}
	dereferencedSchema.Definition = definition

	for k, v := range rawJson {
		if k == "$id" {
			dereferencedSchema.SchemaId = v.(string)
			break
		}
	}

	json, err := json.Marshal(dereferencedSchema.Definition)
	if err != nil {
		return err
	}

	_, err = outFile.Write(append(json, []byte("\n")...))

	return err
}

func (r *RefdSchema) Type() string {
	bytes := []byte(r.SchemaId)
	match := schemaTypePattern.FindSubmatch(bytes)[1]
	artifactType := string(match)

	return artifactType
}

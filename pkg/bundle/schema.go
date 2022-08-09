package bundle

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
)

const idURLPattern = "https://schemas.massdriver.cloud/schemas/bundles/%s/schema-%s.json"
const jsonSchemaURLPattern = "http://json-schema.org/%s/schema"

// Metadata returns common metadata fields for each JSON Schema
func (b *Bundle) Metadata(schemaType string) map[string]string {
	return map[string]string{
		"$schema":     generateSchemaURL(b.Schema),
		"$id":         generateIDURL(b.Name, schemaType),
		"title":       b.Name,
		"description": b.Description,
	}
}

func createFile(dir string, fileName string) (*os.File, error) {
	return os.Create(path.Join(dir, fileName))
}

// Build generates all bundle files in the given bundle
func (b *Bundle) GenerateSchemas(dir string) error {
	err := os.MkdirAll(dir, common.AllRX|common.UserRW)
	if err != nil {
		return err
	}

	paramsSchemaFile, err := createFile(dir, common.ParamsSchemaFilename)
	if err != nil {
		return err
	}

	connectionsSchemaFile, err := createFile(dir, common.ConnectionsSchemaFilename)
	if err != nil {
		return err
	}

	artifactsSchemaFile, err := createFile(dir, common.ArtifactsSchemaFilename)
	if err != nil {
		return err
	}

	uiSchemaFile, err := createFile(dir, common.UISchemaFilename)
	if err != nil {
		return err
	}

	err = GenerateSchema(b.Params, b.Metadata("params"), paramsSchemaFile)
	if err != nil {
		return err
	}
	err = GenerateSchema(b.Connections, b.Metadata("connections"), connectionsSchemaFile)
	if err != nil {
		return err
	}
	err = GenerateSchema(b.Artifacts, b.Metadata("artifacts"), artifactsSchemaFile)
	if err != nil {
		return err
	}

	emptyMetadata := make(map[string]string)
	err = GenerateSchema(b.UI, emptyMetadata, uiSchemaFile)
	if err != nil {
		return err
	}

	err = paramsSchemaFile.Close()
	if err != nil {
		return err
	}
	err = connectionsSchemaFile.Close()
	if err != nil {
		return err
	}
	err = artifactsSchemaFile.Close()
	if err != nil {
		return err
	}

	return nil
}

// generateSchema generates a specific schema-*.json file
func GenerateSchema(schema map[string]interface{}, metadata map[string]string, buffer io.Writer) error {
	var err error
	var mergedSchema = mergeMaps(schema, metadata)

	json, err := json.MarshalIndent(mergedSchema, "", "    ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(buffer, string(json)+"\n")
	if err != nil {
		return err
	}

	return nil
}

func mergeMaps(a map[string]interface{}, b map[string]string) map[string]interface{} {
	for k, v := range b {
		a[k] = v
	}

	return a
}

func generateIDURL(mdName string, schemaType string) string {
	return fmt.Sprintf(idURLPattern, mdName, schemaType)
}

func generateSchemaURL(schema string) string {
	return fmt.Sprintf(jsonSchemaURLPattern, schema)
}

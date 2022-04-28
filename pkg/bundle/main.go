package bundle

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
)

const ArtifactsSchemaFilename = "schema-artifacts.json"
const ConnectionsSchemaFilename = "schema-connections.json"
const ParamsSchemaFilename = "schema-params.json"
const UiSchemaFilename = "schema-ui.json"

const idUrlPattern = "https://schemas.massdriver.cloud/schemas/bundles/%s/schema-%s.json"
const jsonSchemaUrlPattern = "http://json-schema.org/%s/schema"

type BundleStep struct {
	Path        string `json:"path" yaml:"path"`
	Provisioner string `json:"provisioner" yaml:"provisioner"`
}

type Bundle struct {
	Schema      string                 `json:"schema" yaml:"schema"`
	Name        string                 `json:"name" yaml:"name"`
	Description string                 `json:"description" yaml:"description"`
	Ref         string                 `json:"ref" yaml:"ref"`
	Type        string                 `json:"type" yaml:"type"`
	Access      string                 `json:"access" yaml:"access"`
	Steps       []BundleStep           `json:"steps" yaml:"steps"`
	Artifacts   map[string]interface{} `json:"artifacts" yaml:"artifacts"`
	Params      map[string]interface{} `json:"params" yaml:"params"`
	Connections map[string]interface{} `json:"connections" yaml:"connections"`
	Ui          map[string]interface{} `json:"ui" yaml:"ui"`
}

// Metadata returns common metadata fields for each JSON Schema
func (b *Bundle) Metadata(schemaType string) map[string]string {
	return map[string]string{
		"$schema":     generateSchemaUrl(b.Schema),
		"$id":         generateIdUrl(b.Name, schemaType),
		"name":        b.Name,
		"description": b.Description,
	}
}

func createFile(dir string, fileName string) (*os.File, error) {
	return os.Create(path.Join(dir, fileName))
}

// Build generates all bundle files in the given bundle
func (b *Bundle) GenerateSchemas(dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	paramsSchemaFile, err := createFile(dir, ParamsSchemaFilename)
	if err != nil {
		return err
	}

	connectionsSchemaFile, err := createFile(dir, ConnectionsSchemaFilename)
	if err != nil {
		return err
	}

	artifactsSchemaFile, err := createFile(dir, ArtifactsSchemaFilename)
	if err != nil {
		return err
	}

	uiSchemaFile, err := createFile(dir, UiSchemaFilename)
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
	err = GenerateSchema(b.Ui, emptyMetadata, uiSchemaFile)
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

	json, err := json.Marshal(mergedSchema)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(buffer, string(json))
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

func generateIdUrl(mdName string, schemaType string) string {
	return fmt.Sprintf(idUrlPattern, mdName, schemaType)
}

func generateSchemaUrl(schema string) string {
	return fmt.Sprintf(jsonSchemaUrlPattern, schema)
}

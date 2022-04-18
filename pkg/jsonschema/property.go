package jsonschema

// PropertiesMap is a named map of Property
type PropertiesMap map[string]Property

type GenerateAuthFile struct {
	Format   string  `json:"format"`
	Template *string `json:"template,omitempty"`
}

// Property is a single JSON Schema property field
type Property struct {
	AdditionalProperties bool              `json:"additionalProperties"`
	Type                 string            `json:"type"`
	Items                *Property         `json:"items"`
	Properties           PropertiesMap     `json:"properties,omitempty"`
	GenerateAuthFile     *GenerateAuthFile `json:"md.generateAuthFile,omitempty"`
}

// Schema is a flimsy representation of a JSON Schema.
// It provides just enough structure to get type information.
type Schema struct {
	Properties PropertiesMap `json:"properties"`
	Required   []string      `json:"required,omitempty"`
}

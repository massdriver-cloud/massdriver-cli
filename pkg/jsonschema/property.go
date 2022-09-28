package jsonschema

import (
	"encoding/json"
	"fmt"
)

// PropertiesMap is a named map of Property
type PropertiesMap map[string]Property

type GenerateAuthFile struct {
	Format   string  `json:"format"`
	Template *string `json:"template,omitempty"`
}

var Object = "object"

// Property is a single JSON Schema property field
type Property struct {
	AdditionalProperties bool              `json:"additionalProperties"`
	Type                 string            `json:"type"`
	Items                *Property         `json:"items"`
	Properties           PropertiesMap     `json:"properties,omitempty"`
	Default              interface{}       `json:"default,omitempty"`
	GenerateAuthFile     *GenerateAuthFile `json:"md.generateAuthFile,omitempty"`
	Minimum              *float64          `json:"minimum,omitempty"`
}

// Schema is a flimsy representation of a JSON Schema.
// It provides just enough structure to get type information.
type Schema struct {
	Properties PropertiesMap `json:"properties"`
	Required   []string      `json:"required,omitempty"`
	Examples   []Example     `json:"examples,omitempty"`
	Type       string        `json:"type,omitempty"`
}

type Example struct {
	Name   string `json:"__name"`
	Values map[string]interface{}
}

func (e *Example) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}
	name, ok := raw["__name"].(string)
	if !ok {
		return fmt.Errorf("Example name is not a string")
	}
	e.Name = name
	delete(raw, "__name")
	e.Values = raw
	return nil
}

package terraform

import (
	"github.com/massdriver-cloud/massdriver-cli/pkg/jsonschema"
)

// TFVariableFile is a representation of a variables.tf file in JSON format
type TFVariableFile struct {
	Variable map[string]TFVariable `json:"variable"`
}

// In terraform, we need to set "default: null" for non-required fields, however the "default" field should NOT
// be set if the field is required. There isn't a good way to solve this in Golang with a single struct. Thus,
// we need two structs: a TFRequiredVariable which *doesn't* serialize a default value (making it required by TF), and a
// TFOptionalVariable which does serialize a default value (represented as a nil pointer, which serializes to null)
// This interface allows us to work with them interchangably as needed
type TFVariable interface {
	IsTFVariable()
}

// TFRequiredVariable is a representation of a terraform HCL "variable"
type TFRequiredVariable struct {
	Type string `json:"type"`
}

// TFOptionalVariable is a representation of a terraform HCL "variable" with a default of null
type TFOptionalVariable struct {
	Type     string  `json:"type"`
	DoNotSet *string `json:"default"` // DO NOT SET THIS. The intention is to get a nil value for this
}

// Dummy functions to satisfy the interface
func (TFRequiredVariable) IsTFVariable() {}
func (TFOptionalVariable) IsTFVariable() {}

// NewTFVariable creates a new TFVariable from a JSON Schema Property
func NewTFVariable(p jsonschema.Property, required bool) TFVariable {
	t := convertPropertyToType(p)
	if required {
		return TFRequiredVariable{Type: t}
	}
	return TFOptionalVariable{Type: t}
}

func convertPropertyToType(p jsonschema.Property) string {
	switch p.Type {
	case "array":
		return convertArray(p)
	case "object":
		return convertObject(p)
	default:
		return convertScalar(p.Type)
	}
}

func convertObject(prop jsonschema.Property) string {
	objectString := "object({"

	// this comma is used to separate fields. Start it empty (no comma before first element)
	comma := ""
	for key, value := range prop.Properties {
		optional := false
		typeString := convertPropertyToType(value)
		if optional {
			typeString = wrapWithOptional(convertPropertyToType(value))
		}
		objectString += comma + key + "=" + typeString
		comma = ","
	}
	objectString += "})"
	return objectString
}

func convertArray(prop jsonschema.Property) string {
	typeString := "list("
	typeString += convertPropertyToType(*prop.Items)
	typeString += ")"
	return typeString
}

func convertScalar(pType string) string {
	switch pType {
	case "boolean":
		return "bool"
	case "integer":
		return "number"
	default:
		return pType
	}
}

func wrapWithOptional(input string) string {
	return "optional(" + input + ")"
}

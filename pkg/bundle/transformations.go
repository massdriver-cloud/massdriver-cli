package bundle

import (
	"errors"
)

var paramsTransformations = []func(map[string]interface{}) error{AddSetIdToObjectArrays, DisableAdditionalPropertiesInObjects}
var connectionsTransformations = []func(map[string]interface{}) error{DisableAdditionalPropertiesInObjects}
var artifactsTransformations = []func(map[string]interface{}) error{DisableAdditionalPropertiesInObjects}
var uiTransformations = []func(map[string]interface{}) error{}

func ApplyTransformations(schema map[string]interface{}, transformations []func(map[string]interface{}) error) error {

	for _, transformation := range transformations {
		err := transformation(schema)
		if err != nil {
			return err
		}
	}

	for _, v := range schema {
		_, isObject := v.(map[string]interface{})
		if isObject {
			err := ApplyTransformations(v.(map[string]interface{}), transformations)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func AddSetIdToObjectArrays(schema map[string]interface{}) error {
	if schema["type"] == "array" {
		itemsInterface, found := schema["items"]
		if !found {
			return errors.New("found array without items")
		}
		items := itemsInterface.(map[string]interface{})
		if items["type"] == "object" {
			propertiesInterface, found := items["properties"]
			if !found {
				return errors.New("found object without properties")
			}
			properties := propertiesInterface.(map[string]interface{})
			properties["md_set_id"] = map[string]interface{}{"type": "string"}

			requiredInterface, found := items["required"]
			if !found {
				items["required"] = []string{"md_set_id"}
			} else {
				required := requiredInterface.([]interface{})
				items["required"] = append(required, "md_set_id")
			}
		}
	}
	return nil
}

func DisableAdditionalPropertiesInObjects(schema map[string]interface{}) error {
	if schema["type"] == "object" {
		// json schema has a bug where if "anyOf", "allOf" or "oneOf" are used, additionalProperties *MUST* be true
		// we should remove this condition when the bug is fixed
		// https://json-schema.org/understanding-json-schema/reference/combining.html#:~:text=biggest%20surprises
		// https://github.com/massdriver-cloud/xo/issues/53
		_, foundAnyOf := schema["anyOf"]
		_, foundAllOf := schema["allOf"]
		_, foundOneOf := schema["oneOf"]
		if foundAnyOf || foundAllOf || foundOneOf {
			schema["additionalProperties"] = true
		}
		_, found := schema["additionalProperties"]
		if !found {
			schema["additionalProperties"] = false
		}
	}
	return nil
}

package bundle

import (
	"errors"
)

var paramsTransformations = []func(map[string]interface{}) error{}
var connectionsTransformations = []func(map[string]interface{}) error{}
var artifactsTransformations = []func(map[string]interface{}) error{}
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

func AddSetIDToObjectArrays(schema map[string]interface{}) error {
	if schema["type"] == "array" {
		itemsInterface, itemsOK := schema["items"]
		if !itemsOK {
			return errors.New("found array without items")
		}
		items, itemIsObjectOk := itemsInterface.(map[string]interface{})
		if !itemIsObjectOk {
			return errors.New("items is not an object")
		}
		if items["type"] == "object" {
			propertiesInterface, propsOK := items["properties"]
			if !propsOK {
				return errors.New("found object without properties")
			}
			properties, propsIsObjectOk := propertiesInterface.(map[string]interface{})
			if !propsIsObjectOk {
				return errors.New("properties is not an object")
			}
			properties["md_set_id"] = map[string]interface{}{"type": "string"}

			requiredInterface, reqsOK := items["required"]
			if !reqsOK {
				items["required"] = []string{"md_set_id"}
			} else {
				required, ok := requiredInterface.([]interface{})
				if !ok {
					return errors.New("required is not an array")
				}
				items["required"] = append(required, "md_set_id")
			}
		}
	}
	return nil
}

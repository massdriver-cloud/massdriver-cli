package jsonschema

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"regexp"
)

// relativeFilePathPattern only accepts relative file path prefixes "./" and "../"
var relativeFilePathPattern = regexp.MustCompile(`^(\.\/|\.\.\/)`)

func Hydrate(any interface{}, cwd string) (interface{}, error) {
	val := getValue(any)

	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		hydratedList := make([]interface{}, 0)
		for i := 0; i < val.Len(); i++ {
			hydratedVal, err := Hydrate(val.Index(i).Interface(), cwd)
			if err != nil {
				return hydratedList, err
			}
			hydratedList = append(hydratedList, hydratedVal)
		}
		return hydratedList, nil
	case reflect.Map:
		schemaInterface := val.Interface()
		schema := schemaInterface.(map[string]interface{})
		hydratedSchema := map[string]interface{}{}

		// if this part of the schema has a $ref that is a local file, read it and make it
		// the map that we hydrate into. This causes any keys in the ref'ing object to override anything in the ref'd object
		// which adheres to the JSON Schema spec.
		if schemaRefInterface, ok := schema["$ref"]; ok {
			schemaRefPath := schemaRefInterface.(string)
			if relativeFilePathPattern.MatchString(schemaRefPath) {
				// Build up the path from where the dir current schema was read
				schemaRefAbsPath, err := filepath.Abs(filepath.Join(cwd, schemaRefPath))
				if err != nil {
					return hydratedSchema, err
				}

				schemaRefDir := filepath.Dir(schemaRefAbsPath)
				referencedSchema, err := readJsonFile(schemaRefAbsPath)
				if err != nil {
					return hydratedSchema, err
				}

				// Remove it if, so it doesn't get rehydrated below
				delete(schema, "$ref")

				for k, v := range referencedSchema {
					hydratedValue, err := Hydrate(v, schemaRefDir)
					if err != nil {
						return hydratedSchema, err
					}
					hydratedSchema[k] = hydratedValue
				}
			}
		}

		for key, value := range schema {
			var valueInterface = value
			hydratedValue, err := Hydrate(valueInterface, cwd)
			if err != nil {
				return hydratedSchema, err
			}
			hydratedSchema[key] = hydratedValue
		}
		return hydratedSchema, nil
	default:
		return any, nil
	}
}

func getValue(any interface{}) reflect.Value {
	val := reflect.ValueOf(any)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return val
}

func readJsonFile(filepath string) (map[string]interface{}, error) {
	var result map[string]interface{}
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal([]byte(data), &result)

	return result, err
}

// // utility function to extract the "keys" from an OrderedJSONElement Array
// func getKeys(val *reflect.Value) []string {
// 	var keys []string
// 	for i := 0; i < (*val).Len(); i++ {
// 		oje := (*val).Index(i).Interface().(OrderedJSONElement)
// 		keys = append(keys, oje.Key.(string))
// 	}
// 	return keys
// }

// func addAdditionalPropertiesFalseBlock(oj *OrderedJSON) bool {
// 	isObjectBlock := false
// 	additionalPropertiesMissing := true
// 	for _, v := range *oj {
// 		if v.Key == "type" && v.Value == "object" {
// 			isObjectBlock = true
// 		}
// 		if v.Key == "additionalProperties" {
// 			additionalPropertiesMissing = false
// 		}
// 	}
// 	return isObjectBlock && additionalPropertiesMissing
// }

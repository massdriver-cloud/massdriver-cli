package jsonschema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"regexp"

	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
	"github.com/massdriver-cloud/massdriver-cli/pkg/definition"
)

// relativeFilePathPattern only accepts relative file path prefixes "./" and "../"
var relativeFilePathPattern = regexp.MustCompile(`^(\.\/|\.\.\/)`)
var massdriverDefinitionPattern = regexp.MustCompile(`^[a-zA-Z0-9]`)
var httpPattern = regexp.MustCompile(`^(http|https)://`)
var fragmentPattern = regexp.MustCompile(`^#`)

func Hydrate(any interface{}, cwd string, c *client.MassdriverClient) (interface{}, error) {
	val := getValue(any)

	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		hydratedList := make([]interface{}, 0)
		for i := 0; i < val.Len(); i++ {
			hydratedVal, err := Hydrate(val.Index(i).Interface(), cwd, c)
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
			schemaRefValue := schemaRefInterface.(string)
			var referencedSchema map[string]interface{}
			schemaRefDir := cwd
			if relativeFilePathPattern.MatchString(schemaRefValue) {
				// this is a local file ref
				// build up the path from where the dir current schema was read
				schemaRefAbsPath, err := filepath.Abs(filepath.Join(cwd, schemaRefValue))
				if err != nil {
					return hydratedSchema, err
				}

				schemaRefDir = filepath.Dir(schemaRefAbsPath)
				referencedSchema, err = readJsonFile(schemaRefAbsPath)
				if err != nil {
					return hydratedSchema, err
				}

				hydratedSchema, err = replaceRef(schema, referencedSchema, schemaRefDir, c)
				if err != nil {
					return hydratedSchema, err
				}
			} else if httpPattern.MatchString(schemaRefValue) {
				fmt.Println("HTTP/HTTPS refs not supported")
			} else if fragmentPattern.MatchString(schemaRefValue) {
				fmt.Println("Fragment refs not supported")
			} else if massdriverDefinitionPattern.MatchString(schemaRefValue) {
				// this must be a remote schema, so fetch from massdriver
				var err error
				referencedSchema, err = definition.GetDefinition(c, schemaRefValue)
				if err != nil {
					return hydratedSchema, err
				}

				hydratedSchema, err = replaceRef(schema, referencedSchema, schemaRefDir, c)
				if err != nil {
					return hydratedSchema, err
				}
			}
		}

		for key, value := range schema {
			var valueInterface = value
			hydratedValue, err := Hydrate(valueInterface, cwd, c)
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

func replaceRef(base map[string]interface{}, referenced map[string]interface{}, schemaRefDir string, c *client.MassdriverClient) (map[string]interface{}, error) {
	hydratedSchema := map[string]interface{}{}
	delete(base, "$ref")

	for k, v := range referenced {
		hydratedValue, err := Hydrate(v, schemaRefDir, c)
		if err != nil {
			return hydratedSchema, err
		}
		hydratedSchema[k] = hydratedValue
	}
	return hydratedSchema, nil
}

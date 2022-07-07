package jsonschema

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	ctx := context.TODO()

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
		schema := schemaInterface.(map[string]interface{}) //nolint:errcheck
		hydratedSchema := map[string]interface{}{}

		// if this part of the schema has a $ref that is a local file, read it and make it
		// the map that we hydrate into. This causes any keys in the ref'ing object to override anything in the ref'd object
		// which adheres to the JSON Schema spec.
		if schemaRefInterface, ok := schema["$ref"]; ok {
			schemaRefValue := schemaRefInterface.(string) //nolint:errcheck
			var referencedSchema map[string]interface{}
			schemaRefDir := cwd
			if relativeFilePathPattern.MatchString(schemaRefValue) { //nolint:gocritic
				// this is a local file ref
				// build up the path from where the dir current schema was read
				schemaRefAbsPath, err := filepath.Abs(filepath.Join(cwd, schemaRefValue))
				if err != nil {
					return hydratedSchema, err
				}

				schemaRefDir = filepath.Dir(schemaRefAbsPath)
				referencedSchema, err = readJSONFile(schemaRefAbsPath)
				if err != nil {
					return hydratedSchema, err
				}

				hydratedSchema, err = replaceRef(schema, referencedSchema, schemaRefDir, c)
				if err != nil {
					return hydratedSchema, err
				}
			} else if httpPattern.MatchString(schemaRefValue) {
				// HTTP ref. Pull the schema down via HTTP GET and hydrate
				// TODO: this is a security risk as we're blindly doing a get based on a bundle author provided URL
				// see: https://securego.io/docs/rules/g107.html
				// tracked in: https://github.com/massdriver-cloud/massdriver-cli/issues/43
				request, err := http.NewRequestWithContext(ctx, "GET", schemaRefValue, nil)
				if err != nil {
					return hydratedSchema, err
				}
				resp, doErr := c.Client.Do(request)
				if doErr != nil {
					return hydratedSchema, doErr
				}
				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return hydratedSchema, err
				}
				err = json.Unmarshal(body, &referencedSchema)
				if err != nil {
					return hydratedSchema, err
				}

				hydratedSchema, err = replaceRef(schema, referencedSchema, schemaRefDir, c)
				if err != nil {
					return hydratedSchema, err
				}
			} else if fragmentPattern.MatchString(schemaRefValue) {
				fmt.Println("Fragment refs not supported")
			} else if massdriverDefinitionPattern.MatchString(schemaRefValue) {
				// this must be a published schema, so fetch from massdriver
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

func getValue(anyVal interface{}) reflect.Value {
	val := reflect.ValueOf(anyVal)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return val
}

func readJSONFile(filepath string) (map[string]interface{}, error) {
	var result map[string]interface{}
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(data, &result)

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

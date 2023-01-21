package jsonschema

import (
	"context"
	"encoding/json"
	"errors"
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

func Hydrate(ctx context.Context, anyVal interface{}, cwd string, c *client.MassdriverClient) (interface{}, error) {
	val := getValue(anyVal)

	switch val.Kind() { //nolint:exhaustive
	case reflect.Slice, reflect.Array:
		return hydrateList(ctx, c, cwd, val)
	case reflect.Map:
		schemaInterface := val.Interface()
		schema, ok := schemaInterface.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("schema is not an object")
		}
		hydratedSchema := map[string]interface{}{}

		// if this part of the schema has a $ref that is a local file, read it and make it
		// the map that we hydrate into. This causes any keys in the ref'ing object to override anything in the ref'd object
		// which adheres to the JSON Schema spec.
		if schemaRefInterface, refOk := schema["$ref"]; refOk {
			schemaRefValue, refStringOk := schemaRefInterface.(string)
			if !refStringOk {
				return nil, fmt.Errorf("$ref is not a string")
			}
			schemaRefDir := cwd
			var err error
			if relativeFilePathPattern.MatchString(schemaRefValue) { //nolint:gocritic
				// this is a local file ref
				// build up the path from where the dir current schema was read
				hydratedSchema, err = hydrateFilePathRef(ctx, c, cwd, hydratedSchema, schema, schemaRefValue)
			} else if httpPattern.MatchString(schemaRefValue) {
				// HTTP ref. Pull the schema down via HTTP GET and hydrate
				hydratedSchema, err = hydrateHTTPRef(ctx, c, hydratedSchema, schema, schemaRefDir, schemaRefValue)
			} else if fragmentPattern.MatchString(schemaRefValue) {
				fmt.Println("Fragment refs not supported")
			} else if massdriverDefinitionPattern.MatchString(schemaRefValue) {
				// this must be a published schema, so fetch from massdriver
				hydratedSchema, err = hydrateMassdriverRef(ctx, c, hydratedSchema, schema, schemaRefDir, schemaRefValue)
			}
			if err != nil {
				return hydratedSchema, err
			}
		}
		return hydrateMap(ctx, c, cwd, hydratedSchema, schema)
	default:
		return anyVal, nil
	}
}

func hydrateMap(ctx context.Context, c *client.MassdriverClient, cwd string, hydratedSchema map[string]interface{}, schema map[string]interface{}) (map[string]interface{}, error) {
	for key, value := range schema {
		var valueInterface = value
		hydratedValue, err := Hydrate(ctx, valueInterface, cwd, c)
		if err != nil {
			return hydratedSchema, err
		}
		hydratedSchema[key] = hydratedValue
	}
	return hydratedSchema, nil
}

func hydrateList(ctx context.Context, c *client.MassdriverClient, cwd string, val reflect.Value) ([]interface{}, error) {
	hydratedList := make([]interface{}, 0)
	for i := 0; i < val.Len(); i++ {
		hydratedVal, err := Hydrate(ctx, val.Index(i).Interface(), cwd, c)
		if err != nil {
			return hydratedList, err
		}
		hydratedList = append(hydratedList, hydratedVal)
	}
	return hydratedList, nil
}

func hydrateMassdriverRef(ctx context.Context, c *client.MassdriverClient, hydratedSchema map[string]interface{}, schema map[string]interface{}, schemaRefDir string, schemaRefValue string) (map[string]interface{}, error) {
	referencedSchema, err := definition.GetDefinition(c, schemaRefValue)
	if err != nil {
		return hydratedSchema, err
	}

	hydratedSchema, err = replaceRef(ctx, schema, referencedSchema, schemaRefDir, c)
	if err != nil {
		return hydratedSchema, err
	}
	return hydratedSchema, nil
}

func getHTTPRef(ctx context.Context, c *client.MassdriverClient, hydratedSchema map[string]interface{}, schema map[string]interface{}, schemaRefDir string, schemaRefValue string) (map[string]interface{}, error) {
	var referencedSchema map[string]interface{}
	request, err := http.NewRequestWithContext(ctx, "GET", schemaRefValue, nil)
	if err != nil {
		return referencedSchema, err
	}
	resp, doErr := c.Client.Do(request)
	if doErr != nil {
		return referencedSchema, doErr
	}
	if resp.StatusCode != http.StatusOK {
		return referencedSchema, errors.New("received non-200 response getting ref " + resp.Status + " " + schemaRefValue)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return referencedSchema, err
	}
	err = json.Unmarshal(body, &referencedSchema)
	if err != nil {
		return referencedSchema, err
	}
	return referencedSchema, nil
}

func hydrateHTTPRef(ctx context.Context, c *client.MassdriverClient, hydratedSchema map[string]interface{}, schema map[string]interface{}, schemaRefDir string, schemaRefValue string) (map[string]interface{}, error) {
	var referencedSchema map[string]interface{}
	hydratedSchema, errGet := getHTTPRef(ctx, c, hydratedSchema, schema, schemaRefDir, schemaRefValue)
	if errGet != nil {
		return nil, errGet
	}

	// the local refs can be looked up as well
	hydratedSchema, errReplace := optimisticallyReplaceRef(ctx, schema, referencedSchema, schemaRefDir, c)
	if errReplace != nil {
		return hydratedSchema, errReplace
	}

	return hydratedSchema, nil
}

func hydrateFilePathRef(ctx context.Context, c *client.MassdriverClient, cwd string, hydratedSchema map[string]interface{}, schema map[string]interface{}, schemaRefValue string) (map[string]interface{}, error) {
	var referencedSchema map[string]interface{}
	var schemaRefDir string
	schemaRefAbsPath, err := filepath.Abs(filepath.Join(cwd, schemaRefValue))
	if err != nil {
		return hydratedSchema, err
	}

	schemaRefDir = filepath.Dir(schemaRefAbsPath)
	referencedSchema, readErr := readJSONFile(schemaRefAbsPath)
	if readErr != nil {
		return hydratedSchema, readErr
	}

	var replaceErr error
	hydratedSchema, replaceErr = replaceRef(ctx, schema, referencedSchema, schemaRefDir, c)
	if replaceErr != nil {
		return hydratedSchema, replaceErr
	}
	return hydratedSchema, nil
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

func replaceRef(ctx context.Context, base map[string]interface{}, referenced map[string]interface{}, schemaRefDir string, c *client.MassdriverClient) (map[string]interface{}, error) {
	hydratedSchema := map[string]interface{}{}
	delete(base, "$ref")

	for k, v := range referenced {
		hydratedValue, err := Hydrate(ctx, v, schemaRefDir, c)
		if err != nil {
			return hydratedSchema, err
		}
		hydratedSchema[k] = hydratedValue
	}
	return hydratedSchema, nil
}

func optimisticallyReplaceRef(ctx context.Context, base map[string]interface{}, referenced map[string]interface{}, schemaRefDir string, c *client.MassdriverClient) (map[string]interface{}, error) {
	hydratedSchema := map[string]interface{}{}
	delete(base, "$ref")

	for k, v := range referenced {
		panic(v)
		hydratedValue, err := Hydrate(ctx, v, schemaRefDir, c)
		if err != nil {
			return hydratedSchema, err
		}
		hydratedSchema[k] = hydratedValue
	}
	return hydratedSchema, nil
}

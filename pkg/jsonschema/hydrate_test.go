package jsonschema_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/jsonschema"
)

type TestCase struct {
	Name     string
	Input    interface{}
	Expected interface{}
}

func TestHydrate(t *testing.T) {
	cases := []TestCase{
		{
			Name:  "Hydrates a $ref",
			Input: jsonDecode(`{"$ref": "./testdata/artifacts/aws-example.json"}`),
			Expected: map[string]string{
				"id": "fake-schema-id",
			},
		},
		{
			Name:  "Hydrates a $ref alongside arbitrary values",
			Input: jsonDecode(`{"foo": true, "bar": {}, "$ref": "./testdata/artifacts/aws-example.json"}`),
			Expected: map[string]interface{}{
				"foo": true,
				"bar": map[string]interface{}{},
				"id":  "fake-schema-id",
			},
		},
		{
			Name:  "Hydrates a nested $ref",
			Input: jsonDecode(`{"key": {"$ref": "./testdata/artifacts/aws-example.json"}}`),
			Expected: map[string]map[string]string{
				"key": {
					"id": "fake-schema-id",
				},
			},
		},
		{
			Name:  "Does not hydrate HTTPS refs",
			Input: jsonDecode(`{"$ref": "https://elsewhere.com/schema.json"}`),
			Expected: map[string]string{
				"$ref": "https://elsewhere.com/schema.json",
			},
		},
		{
			Name:  "Does not hydrate fragment (#) refs",
			Input: jsonDecode(`{"$ref": "#/its-in-this-file"}`),
			Expected: map[string]string{
				"$ref": "#/its-in-this-file",
			},
		},
		{
			Name:  "Hydrates $refs in a list",
			Input: jsonDecode(`{"list": ["string", {"$ref": "./testdata/artifacts/aws-example.json"}]}`),
			Expected: map[string]interface{}{
				"list": []interface{}{
					"string",
					map[string]interface{}{
						"id": "fake-schema-id",
					},
				},
			},
		},
		{
			Name:  "Hydrates a $ref deterministically (keys outside of ref always win)",
			Input: jsonDecode(`{"conflictingKey": "not-from-ref", "$ref": "./testdata/artifacts/conflicting-keys.json"}`),
			Expected: map[string]string{
				"conflictingKey": "not-from-ref",
				"nonConflictKey": "from-ref",
			},
		},
		{
			Name:  "Hydrates a $ref recursively",
			Input: jsonDecode(`{"$ref": "./testdata/artifacts/ref-aws-example.json"}`),
			Expected: map[string]map[string]string{
				"properties": {
					"id": "fake-schema-id",
				},
			},
		},
		{
			Name:  "Hydrates a $ref recursively",
			Input: jsonDecode(`{"$ref": "./testdata/artifacts/ref-lower-dir-aws-example.json"}`),
			Expected: map[string]map[string]string{
				"properties": {
					"id": "fake-schema-id",
				},
			},
		},
		// {
		// 	Name:  `Adds "additionalProperties":false to object types`,
		// 	Input: jsonDecode(`{"properties": {"a": "b"}, "type": "object"}`),
		// 	Expected: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
		// 		{Key: "properties", Value: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
		// 			{Key: "a", Value: "b"},
		// 		})},
		// 		{Key: "type", Value: "object"},
		// 		{Key: "additionalProperties", Value: false},
		// 	}),
		// },
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			got, _ := jsonschema.Hydrate(test.Input, ".")

			if fmt.Sprint(got) != fmt.Sprint(test.Expected) {
				t.Errorf("got %v, want %v", got, test.Expected)
			}
		})
	}
}

func jsonDecode(data string) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal([]byte(data), &result)
	return result
}

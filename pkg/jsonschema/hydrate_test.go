package jsonschema_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
	"github.com/massdriver-cloud/massdriver-cli/pkg/jsonschema"
)

type TestCase struct {
	Name                string
	Input               interface{}
	Expected            interface{}
	ExpectedErrorSuffix string
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
		{
			Name:  "Reports not found when $ref is not found",
			Input: jsonDecode(`{"$ref": "./testdata/no-type.json"}`),
			// Expected: map[string]map[string]string{
			// 	"properties": {
			// 		"id": "fake-schema-id",
			// 	},
			// },
			ExpectedErrorSuffix: "testdata/no-type.json: no such file or directory",
		},
		{
			Name:  "Hydrates remote (massdriver) ref",
			Input: jsonDecode(`{"$ref": "massdriver/test-schema"}`),
			Expected: map[string]string{
				"foo": "bar",
			},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				urlPath := r.URL.Path
				switch urlPath {
				case "/artifact-definitions/massdriver/test-schema":
					if _, err := w.Write([]byte(`{"foo":"bar"}`)); err != nil {
						t.Fatalf("Failed to write response: %v", err)
					}
				default:
					t.Fatalf("unknown schema: %v", urlPath)
				}
			}))
			defer testServer.Close()

			c := client.NewClient().WithEndpoint(testServer.URL)
			ctx := context.TODO()

			got, gotErr := jsonschema.Hydrate(ctx, test.Input, ".", c)

			if test.ExpectedErrorSuffix != "" {
				if !strings.HasSuffix(gotErr.Error(), test.ExpectedErrorSuffix) {
					t.Errorf("got %v, want %v", gotErr.Error(), test.ExpectedErrorSuffix)
				}
			} else {
				if fmt.Sprint(got) != fmt.Sprint(test.Expected) {
					t.Errorf("got %v, want %v", got, test.Expected)
				}
			}
		})
	}

	// Easier to test HTTP refs separately
	t.Run("HTTP Refs", func(t *testing.T) {
		var recursivePtr *string
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			urlPath := r.URL.Path
			switch urlPath {
			case "/recursive":
				if _, err := w.Write([]byte(*recursivePtr)); err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
				fmt.Println("in recursive")
			case "/endpoint":
				if _, err := w.Write([]byte(`{"foo":"bar"}`)); err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
				fmt.Println("in endpoint")
			default:
				w.WriteHeader(http.StatusNotFound)
				_, err := w.Write([]byte(`404 - not found`))
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}
		}))
		defer testServer.Close()

		c := client.NewClient().WithEndpoint(testServer.URL)
		ctx := context.TODO()

		recursive := fmt.Sprintf(`{"baz":{"$ref":"%s/endpoint"}}`, testServer.URL)
		recursivePtr = &recursive

		input := jsonDecode(fmt.Sprintf(`{"$ref":"%s/recursive"}`, testServer.URL))

		got, _ := jsonschema.Hydrate(ctx, input, ".", c)
		expected := map[string]interface{}{
			"baz": map[string]string{
				"foo": "bar",
			},
		}

		if fmt.Sprint(got) != fmt.Sprint(expected) {
			t.Errorf("got %v, want %v", got, expected)
		}

		input = jsonDecode(fmt.Sprintf(`{"$ref":"%s/not-found"}`, testServer.URL))
		_, gotErr := jsonschema.Hydrate(ctx, input, ".", c)
		expectedErrPrefix := "received non-200 response getting ref 404 Not Found"

		if !strings.HasPrefix(gotErr.Error(), expectedErrPrefix) {
			t.Errorf("got %v, want %v", gotErr.Error(), expectedErrPrefix)
		}
	})
}

func jsonDecode(data string) map[string]interface{} {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		panic(err)
	}
	return result
}

package terraform_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/provisioners/terraform"
)

// Helper function for asserting json serde matches
func doc(str string) string {
	b := []byte(str)

	jsonMap := make(map[string](interface{}))
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		panic(err)
	}

	outBytes, _ := json.MarshalIndent(jsonMap, "", "    ")
	return string(outBytes)
}

func TestGenerateFiles(t *testing.T) {
	type testData struct {
		name       string
		bundlePath string
		srcDir     string
		expected   map[string]string
	}
	tests := []testData{
		{
			name:       "standard",
			bundlePath: "testdata/testbundle/",
			srcDir:     "src",
			expected: map[string]string{
				"_connections_variables.tf.json": `{
    "variable": {
        "foo": {
            "type": "string"
        }
    }
}
`,
				"_params_variables.tf.json": `{
    "variable": {
        "age": {
            "type": "number",
            "default": null
        },
        "name": {
            "type": "string"
        },
        "status": {
            "type": "any",
            "default": null
        }
    }
}
`,
				// Note the age 27 checks that existing values are not overwritten
				"_params.auto.tfvars.json": `{
    "age": 27,
    "md_metadata": {
        "default_tags": {
            "md-manifest": "testbundle",
            "md-package": "local-dev-testbundle-000",
            "md-project": "local",
            "md-target": "dev"
        },
        "name_prefix": "local-dev-testbundle-000"
    },
    "name": "John Doe",
    "status": {
        "alive": "TODO: REPLACE ME",
        "daysSinceLastCrime": 0,
        "knownConvictions": [],
        "relationship": "single",
        "someOtherExistingNestedValue": "foo"
    }
}
`,
				"_md_variables.tf.json": `{
    "variable": {
        "md_metadata": {
            "type": "any"
        }
    }
}
`,
			},
		},
		{
			name:       "missing params",
			bundlePath: "testdata/testbundle-broken/",
			srcDir:     "src",
			expected: map[string]string{
				"_connections_variables.tf.json": `{
    "variable": {
        "foo": {
            "type": "string"
        }
    }
}
`,
				"_params_variables.tf.json": `{
    "variable": {
        "age": {
            "type": "number",
            "default": null
        },
        "name": {
            "type": "string"
        },
        "status": {
            "type": "any",
            "default": null
        }
    }
}
`,
				"_params.auto.tfvars.json": `{
    "age": 25,
    "md_metadata": {
        "default_tags": {
            "md-manifest": "testbundle-broken",
            "md-package": "local-dev-testbundle-broken-000",
            "md-project": "local",
            "md-target": "dev"
        },
        "name_prefix": "local-dev-testbundle-broken-000"
    },
    "name": "John Doe",
    "status": {
        "alive": "TODO: REPLACE ME",
        "daysSinceLastCrime": 0,
        "knownConvictions": [],
        "relationship": "single",
        "someOtherExistingNestedValue": "TODO: REPLACE ME"
    }
}
`,
				"_md_variables.tf.json": `{
    "variable": {
        "md_metadata": {
            "type": "any"
        }
    }
}
`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := terraform.GenerateFiles(tc.bundlePath, tc.srcDir)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			for file, want := range tc.expected {
				got, readErr := os.ReadFile(path.Join(tc.bundlePath, tc.srcDir, file))
				if readErr != nil {
					t.Fatalf("%d, unexpected error", readErr)
				}

				if string(got) != want {
					t.Errorf("got %s want %s", string(got), want)
				}
			}
		})
	}
}
func TestCompile(t *testing.T) {
	type testData struct {
		name       string
		schemaPath string
		expected   string
	}
	tests := []testData{
		{
			name:       "populated schema",
			schemaPath: "file://./testdata/local-schema.json",
			expected: doc(`
{
	"variable": {
		"name": {
			"type": "string"
		},
		"age": {
			"type": "number"
		},
		"active": {
			"type": "bool"
		},
		"height": {
			"type": "number"
		}
	}
}`) + "\n"},
		{
			name:       "empty schema",
			schemaPath: "file://./testdata/empty-schema.json",
			expected:   doc("{}"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got bytes.Buffer
			err := terraform.Compile(tc.schemaPath, &got)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}
			want := tc.expected

			if got.String() != want {
				t.Errorf("got %s want %s", got.String(), want)
			}
		})
	}
}

// https://github.com/xeipuuv/gojsonschema#loading-local-schemas
// This test is failing because the library doesnt automatically
// resolve $refs until a document is validated. You can trick it into
// doing it w/ the last example mentioned in the above link, but
// we will need to have an idea of how we are doing that in bundles
// first. I assume we'll end up treating the bundle's JSON Schema as the main
// and ref loading a single 'compile' JSON Schema that has all of our
// secrets and connections
// func TestCompileRemoteSchemas(t *testing.T) {
// 	got, _  := Compile("file://./testdata/remote-schema.json")
// 	want := doc(`
// 	{
// 		"variable": {
// 			"local": {
// 				"type": "string"
// 			},
// 			"remote": {
// 				"type": "string"
// 			}
// 		}
// 	}
// `)

// 	if got != want {
// 		t.Errorf("got %s want %s", got, want)
// 	}
// }

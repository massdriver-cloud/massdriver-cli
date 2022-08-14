package definition_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/definition"
)

func TestParse(t *testing.T) {
	var got, _ = definition.Parse("./testdata/aws-iam-role.json")
	var want = definition.DefinitionFile{
		// Schema: "",
		Md: definition.Md{
			Access: "public",
			Name:   "aws-iam-role",
			Provisioners: definition.Priovisioners{
				Terraform: "",
			},
		},
		// Type:                 "",
		// Title:                "",
		AdditionalProperties: false,
		// Required:             []string{},
		Properties: map[string]interface{}{
			"data": map[string]interface{}{
				"title": "Artifact Data",
				"type":  "object",
				"required": []string{
					"infrastructure",
				},
			},
		},
	}

	gotBytes, _ := json.Marshal(got)
	wantBytes, _ := json.Marshal(want)

	if !bytes.Equal(gotBytes, wantBytes) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

package definition_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/definition"
)

func TestBuild(t *testing.T) {
	type test struct {
		defPath string
		want    definition.DefinitionFile
	}
	tests := []test{
		{
			defPath: "testdata/simple-def.json",
			want: definition.DefinitionFile{
				Md: definition.Md{
					Access: "public",
					Name:   "simple-def",
					Provisioners: definition.Priovisioners{
						Terraform: "dGVycmFmb3JtIHsKICByZXF1aXJlZF92ZXJzaW9uID0gIj49IDAuMTIuMCIKICByZXF1aXJlZF9wcm92aWRlcnMgewogICAgYXdzID0gewogICAgICBzb3VyY2UgID0gImhhc2hpY29ycC9hd3MiCiAgICAgIHZlcnNpb24gPSAiPj0gMi4wLjAiCiAgICB9CiAgfQp9Cg==",
					},
				},
				Properties: map[string]interface{}{
					"data":  map[string]interface{}{},
					"specs": map[string]interface{}{},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.defPath, func(t *testing.T) {
			got, err := definition.Build(tc.defPath)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			gotBytes, _ := json.Marshal(got)
			wantBytes, _ := json.Marshal(tc.want)

			if !bytes.Equal(gotBytes, wantBytes) {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}

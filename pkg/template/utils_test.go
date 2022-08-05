package template_test

import (
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/template"
)

func TestTypeToName(t *testing.T) {
	type test struct {
		artifactName string
		want         string
	}
	tests := []test{
		{
			artifactName: "massdriver/aws-iam-role",
			want:         "massdriver_aws_iam_role",
		},
	}

	for _, tc := range tests {
		t.Run(tc.artifactName, func(t *testing.T) {
			got := template.TypeToName(tc.artifactName)
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

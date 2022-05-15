package application_test

import (
	"reflect"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/application"
)

func TestParse(t *testing.T) {
	type test struct {
		name    string
		appPath string
		want    application.Application
	}
	tests := []test{
		{
			name:    "simple",
			appPath: "./testdata/appsimple.yaml",
			want: application.Application{
				Schema:      "draft-07",
				Name:        "my-app",
				Description: "An application",
				Ref:         "github.com/user/app",
				Access:      "private",
				Deployment: application.ApplicationDeployment{
					Type: "simple",
				},
				Params: map[string]interface{}{
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":  "string",
							"title": "Name",
						},
					},
					"required": []interface{}{
						"name",
					},
				},
			},
		},
		{
			name:    "custom",
			appPath: "./testdata/appcustom.yaml",
			want: application.Application{
				Schema:      "draft-07",
				Name:        "my-app",
				Description: "An application",
				Ref:         "github.com/user/app",
				Access:      "private",
				Deployment: application.ApplicationDeployment{
					Type: "custom",
					Path: "./test-chart",
				},
				Params: map[string]interface{}{
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":  "string",
							"title": "Name",
						},
					},
					"required": []interface{}{
						"name",
					},
				},
			},
		},
		{
			name:    "deps",
			appPath: "./testdata/appdeps.yaml",
			want: application.Application{
				Schema:      "draft-07",
				Name:        "my-app",
				Description: "An application",
				Ref:         "github.com/user/app",
				Access:      "private",
				Deployment: application.ApplicationDeployment{
					Type: "simple",
				},
				Params: map[string]interface{}{
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":  "string",
							"title": "Name",
						},
						"age": map[string]interface{}{
							"type":  "integer",
							"title": "Age",
						},
					},
					"required": []interface{}{
						"name",
					},
				},
				Dependencies: []application.ApplicationDependencies{
					{
						Type:     "massdriver/rdbms-authentication",
						Field:    "database",
						Required: true,
						Env: []application.ApplicationDependenciesEnvs{
							{
								Name:  "DATABASE_URL",
								Value: "${data.authentication.connection_string}",
							},
						},
					},
					{
						Type:     "massdriver/aws-sqs-queue",
						Field:    "queue",
						Required: false,
						Env: []application.ApplicationDependenciesEnvs{
							{
								Name:  "MY_QUEUE_ARN",
								Value: "${data.infrastructure.arn}",
							},
						},
						Policy: "read",
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			got, err := application.Parse(tc.appPath)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if !reflect.DeepEqual(*got, tc.want) {
				t.Errorf("got %v, want %v", *got, tc.want)
			}
		})
	}
}
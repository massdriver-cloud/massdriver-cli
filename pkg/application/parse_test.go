package application_test

import (
	"reflect"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/application"
)

func TestParseBundle(t *testing.T) {
	var got, _ = application.Parse("./testdata/app.yaml")
	var want = application.Application{
		Schema:      "draft-07",
		Name:        "my-app",
		Description: "An application",
		Ref:         "github.com/user/app",
		Access:      "private",
		Deployment: application.ApplicationDeployment{
			Type: "chart",
		},
		Params: map[string]interface{}{
			"properties": map[interface{}]interface{}{
				"name": map[interface{}]interface{}{
					"type":  "string",
					"title": "Name",
				},
				"age": map[interface{}]interface{}{
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
	}

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("got %v, want %v", *got, want)
	}
}

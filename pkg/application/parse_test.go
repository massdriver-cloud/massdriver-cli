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
			name:    "app-spec",
			appPath: "./testdata/massdriver.yaml",
			want: application.Application{
				Schema:      "draft-07",
				Name:        "my-app",
				Description: "An application",
				Ref:         "github.com/user/app",
				Type:        "application",
				Access:      "private",
				Params: map[string]interface{}{
					"properties": map[string]interface{}{
						"log_level": map[string]interface{}{
							"enum": []interface{}{"warn", "error", "info"},
							"type": "string",
						},
						"name": map[string]interface{}{
							"type": "string",
						},
						"namespace": map[string]interface{}{
							"default": "default",
							"type":    "string",
						},
						"replication": map[string]interface{}{
							"enum": []interface{}{"async", "sync"},
							"type": "string",
						},
					},
				},
				Connections: map[string]interface{}{
					"properties": map[string]interface{}{
						"kubernetes_cluster": map[string]interface{}{
							"$ref": "massdriver/k8s",
						},
						"mongo": map[string]interface{}{
							"$ref": "massdriver/mongo-authentication",
						},
						"sqs": map[string]interface{}{
							"$ref": "massdriver/aws-sqs-pubsub-subscription",
						},
					},
					"required": []interface{}{
						"*",
					},
				},
				App: application.AppBlock{
					Envs: map[string]string{
						"LOG_LEVEL":      "params.log_level",
						"MONGO_USERNAME": "connections.mongo.authentication.username",
						"STRIPE_KEY":     "secrets.ecomm_site_stripe_key",
					},
					Policies: []string{
						"connections.sqs.security.policies.read",
						"connections.s3.security.policies.write",
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := application.Parse(tc.appPath, nil)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if !reflect.DeepEqual(*got, tc.want) {
				t.Errorf("got %v, want %v", *got, tc.want)
			}
		})
	}
}

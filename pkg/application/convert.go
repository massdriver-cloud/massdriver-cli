package application

import (
	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
)

func (app *Application) ConvertToBundle() (*bundle.Bundle, error) {
	b := new(bundle.Bundle)

	b.Schema = app.Schema
	b.Name = app.Title
	b.Description = app.Description
	b.Ref = app.Ref
	b.Type = "application"
	b.Access = app.Access

	b.Steps = []bundle.Step{
		{
			Path:        "src",
			Provisioner: "terraform",
		},
	}
	b.Params = app.Params
	b.Connections = make(map[string]interface{})
	b.Artifacts = make(map[string]interface{})
	b.UI = make(map[string]interface{})

	// default connections are kubernetes and cloud auth
	// connectionsRequired := []string{"kubernetes-cluster", "cloud-authentication"}
	connectionsRequired := []string{"kubernetes_cluster"}
	connectionsProperties := make(map[string]interface{})

	connectionsProperties["kubernetes_cluster"] = map[string]interface{}{
		"$ref": "massdriver/kubernetes-cluster",
	}
	// connectionsProperties["cloud-authentication"] = map[string]interface{}{
	// 	"oneOf": []interface{}{
	// 		map[string]interface{}{"$ref": "massdriver/aws-iam-role"},
	// 		map[string]interface{}{"$ref": "massdriver/azure-service-principal"},
	// 		map[string]interface{}{"$ref": "massdriver/gcp-service-account"},
	// 	},
	// }
	connectionsProperties["aws_authentication"] = map[string]interface{}{
		"$ref": "massdriver/aws-iam-role",
	}
	connectionsProperties["azure_authentication"] = map[string]interface{}{
		"$ref": "massdriver/azure-service-principal",
	}
	connectionsProperties["gcp_authentication"] = map[string]interface{}{
		"$ref": "massdriver/gcp-service-account",
	}

	for depKey, dep := range app.Dependencies {
		if dep.Required {
			connectionsRequired = append(connectionsRequired, depKey)
		}
		connectionsProperties[depKey] = map[string]interface{}{
			"$ref": dep.Type,
		}
	}

	b.Connections["required"] = connectionsRequired
	b.Connections["properties"] = connectionsProperties
	b.Artifacts["properties"] = make(map[string]interface{})
	b.UI = app.UI

	return b, nil
}

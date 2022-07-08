package application

import (
	"encoding/json"
	"fmt"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
)

func (app *Application) ConvertToBundle() (*bundle.Bundle, error) {
	b := new(bundle.Bundle)

	b.Schema = app.Schema
	b.Name = app.Name
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

	if app.Deployment.Type == "simple" {
		b.Params = make(map[string]interface{})
		if err := json.Unmarshal([]byte(SimpleParams), &b.Params); err != nil {
			return b, fmt.Errorf("error parsing simple params: %w", err)
		}
	} else {
		b.Params = app.Params
	}

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

	for _, dep := range app.Dependencies {
		if dep.Required {
			connectionsRequired = append(connectionsRequired, dep.Field)
		}
		connectionsProperties[dep.Field] = map[string]interface{}{
			"$ref": dep.Type,
		}
	}

	b.Connections["required"] = connectionsRequired
	b.Connections["properties"] = connectionsProperties

	// default artifact is kubernetes-application
	// TODO: RE-ENABLE THIS WHEN WE HAVE A WORKING ARTIFACT
	// artifactsRequired := []string{"kubernetes-application"}
	// artifactsProperties := make(map[string]interface{})

	// artifactsProperties["kubernetes-application"] = map[string]interface{}{
	// 	"$ref": "massdriver/kubernetes-application",
	// }

	// b.Artifacts["required"] = artifactsRequired
	// b.Artifacts["properties"] = artifactsProperties

	b.Artifacts["properties"] = make(map[string]interface{})

	// UI
	if app.Deployment.Type == "simple" {
		b.UI = make(map[string]interface{})
		if jsonErr := json.Unmarshal([]byte(SimpleUI), &b.UI); jsonErr != nil {
			return b, fmt.Errorf("error parsing simple ui: %w", jsonErr)
		}
	} else {
		uiOrder := []interface{}{"*"}
		b.UI["ui:order"] = uiOrder
	}

	return b, nil
}

package application

import (
	"encoding/json"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
)

var MASSDRIVER_URL string = "https://api.massdriver.cloud/"

func (app *Application) ConvertToBundle() *bundle.Bundle {
	b := new(bundle.Bundle)

	b.Schema = app.Schema
	b.Name = app.Name
	b.Description = app.Description
	b.Ref = app.Ref
	b.Type = "application"
	b.Access = app.Access

	b.Steps = []bundle.BundleStep{
		{
			Path:        "src",
			Provisioner: "terraform",
		},
	}

	if app.Deployment.Type == "simple" {
		b.Params = make(map[string]interface{})
		json.Unmarshal([]byte(simpleParams), &b.Params)
	} else {
		b.Params = app.Params
	}

	b.Connections = make(map[string]interface{})
	b.Artifacts = make(map[string]interface{})
	b.Ui = make(map[string]interface{})

	// default connections are kubernetes and cloud auth
	//connectionsRequired := []string{"kubernetes-cluster", "cloud-authentication"}
	connectionsRequired := []string{"kubernetes-cluster"}
	connectionsProperties := make(map[string]interface{})

	connectionsProperties["kubernetes-cluster"] = map[string]interface{}{
		"$ref": "massdriver/kubernetes-cluster",
	}
	// connectionsProperties["cloud-authentication"] = map[string]interface{}{
	// 	"oneOf": []interface{}{
	// 		map[string]interface{}{"$ref": "massdriver/aws-iam-role"},
	// 		map[string]interface{}{"$ref": "massdriver/azure-service-principal"},
	// 		map[string]interface{}{"$ref": "massdriver/gcp-service-account"},
	// 	},
	// }
	connectionsProperties["aws-authentication"] = map[string]interface{}{
		"$ref": "massdriver/aws-iam-role",
	}
	connectionsProperties["azure-authentication"] = map[string]interface{}{
		"$ref": "massdriver/azure-service-principal",
	}
	connectionsProperties["gcp-authentication"] = map[string]interface{}{
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

	// NOT NEEDED
	uiOrder := []interface{}{"*"}
	b.Ui["ui:order"] = uiOrder

	return b
}

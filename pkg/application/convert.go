package application

import (
	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
)

var MASSDRIVER_URL string = "https://api.massdriver.cloud/"

func (app *Application) ConvertToBundle() *bundle.Bundle {
	b := new(bundle.Bundle)

	b.Schema = app.Schema
	b.Name = app.Name
	b.Description = app.Description
	b.Ref = app.Ref
	b.Access = app.Access
	b.Type = "application"
	b.Params = app.Params
	b.Connections = make(map[string]interface{})
	b.Artifacts = make(map[string]interface{})

	// default connections are kubernetes and cloud auth
	connectionsRequired := []string{"kubernetes-cluster", "cloud-authentication"}
	connectionsProperties := make(map[string]interface{})

	connectionsProperties["kubernetes-cluster"] = map[string]string{
		"$ref": "massdriver/kubernetes-cluster",
	}
	connectionsProperties["cloud-authentication"] = map[string]string{
		"$ref": "massdriver/cloud-authentication",
	}

	for _, dep := range app.Dependencies {
		if dep.Required {
			connectionsRequired = append(connectionsRequired, dep.Field)
		}
		connectionsProperties[dep.Field] = map[string]string{
			"$ref": dep.Type,
		}
	}

	b.Connections["required"] = connectionsRequired
	b.Connections["properties"] = connectionsProperties

	// default artifact is kubernetes-application
	artifactsRequired := []string{"kubernetes-application"}
	artifactsProperties := make(map[string]interface{})

	artifactsProperties["kubernetes-application"] = map[string]string{
		"$ref": "massdriver/kubernetes-application",
	}

	b.Artifacts["required"] = artifactsRequired
	b.Artifacts["properties"] = artifactsProperties

	return b
}

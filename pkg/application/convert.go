package application

import (
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
	b.Params = app.Params
	b.Connections = make(map[string]interface{})
	b.Artifacts = make(map[string]interface{})
	b.UI = make(map[string]interface{})

	connectionsRequired := []string{}
	connectionsProperties := make(map[string]interface{})

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

	if app.UI != nil {
		b.UI = app.UI
	} else {
		uiOrder := []interface{}{"*"}
		b.UI["ui:order"] = uiOrder
	}

	return b, nil
}

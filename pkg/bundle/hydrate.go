package bundle

import (
	"path/filepath"

	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
	"github.com/massdriver-cloud/massdriver-cli/pkg/jsonschema"
)

func (b *Bundle) Hydrate(path string, c *client.MassdriverClient) error {
	cwd := filepath.Dir(path)

	hydratedArtifacts, err := jsonschema.Hydrate(b.Artifacts, cwd, c)
	if err != nil {
		return err
	}
	b.Artifacts = hydratedArtifacts.(map[string]interface{})
	err = ApplyTransformations(b.Artifacts, artifactsTransformations)
	if err != nil {
		return err
	}

	hydratedParams, err := jsonschema.Hydrate(b.Params, cwd, c)
	if err != nil {
		return err
	}
	b.Params = hydratedParams.(map[string]interface{})
	err = ApplyTransformations(b.Params, paramsTransformations)
	if err != nil {
		return err
	}

	hydratedConnections, err := jsonschema.Hydrate(b.Connections, cwd, c)
	if err != nil {
		return err
	}
	b.Connections = hydratedConnections.(map[string]interface{})
	err = ApplyTransformations(b.Connections, connectionsTransformations)
	if err != nil {
		return err
	}

	hydratedUi, err := jsonschema.Hydrate(b.Ui, cwd, c)
	if err != nil {
		return err
	}
	b.Ui = hydratedUi.(map[string]interface{})
	err = ApplyTransformations(b.Ui, uiTransformations)
	if err != nil {
		return err
	}

	return nil
}

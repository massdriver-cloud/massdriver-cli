package bundle

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
	"github.com/massdriver-cloud/massdriver-cli/pkg/jsonschema"
)

func (b *Bundle) Hydrate(path string, c *client.MassdriverClient) error {
	cwd := filepath.Dir(path)
	ctx := context.TODO()

	hydratedArtifacts, err := jsonschema.Hydrate(ctx, b.Artifacts, cwd, c)
	if err != nil {
		return err
	}
	var artOk bool
	b.Artifacts, artOk = hydratedArtifacts.(map[string]interface{})
	if !artOk {
		return errors.New("hydrated artifacts is not a map")
	}
	err = ApplyTransformations(b.Artifacts, artifactsTransformations)
	if err != nil {
		return err
	}

	hydratedParams, err := jsonschema.Hydrate(ctx, b.Params, cwd, c)
	if err != nil {
		return err
	}
	var paramOk bool
	b.Params, paramOk = hydratedParams.(map[string]interface{})
	if !paramOk {
		return errors.New("hydrated params is not a map")
	}
	err = ApplyTransformations(b.Params, paramsTransformations)
	if err != nil {
		return err
	}

	hydratedConnections, err := jsonschema.Hydrate(ctx, b.Connections, cwd, c)
	if err != nil {
		return err
	}
	var connOk bool
	b.Connections, connOk = hydratedConnections.(map[string]interface{})
	if !connOk {
		return errors.New("hydrated connections is not a map")
	}
	err = ApplyTransformations(b.Connections, connectionsTransformations)
	if err != nil {
		return err
	}

	hydratedUI, err := jsonschema.Hydrate(ctx, b.UI, cwd, c)
	if err != nil {
		return err
	}
	var uiOk bool
	b.UI, uiOk = hydratedUI.(map[string]interface{})
	if !uiOk {
		return errors.New("hydrated UI is not a map")
	}
	err = ApplyTransformations(b.UI, uiTransformations)
	if err != nil {
		return err
	}

	return nil
}

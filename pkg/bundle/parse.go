package bundle

import (
	"io/ioutil"
	"path/filepath"

	"github.com/massdriver-cloud/massdriver-cli/pkg/jsonschema"
	"gopkg.in/yaml.v3"
)

// ParseBundle parses a bundle from a YAML file
// bundle, err := ParseBundle("./bundle.yaml")
// Generate the files in this directory
// bundle.Build(".")
func ParseBundle(path string) (Bundle, error) {
	bundle := Bundle{}
	cwd := filepath.Dir(path)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return bundle, err
	}

	err = yaml.Unmarshal([]byte(data), &bundle)
	if err != nil {
		return bundle, err
	}

	hydratedArtifacts, err := jsonschema.Hydrate(bundle.Artifacts, cwd)
	if err != nil {
		return bundle, err
	}
	bundle.Artifacts = hydratedArtifacts.(map[string]interface{})
	err = ApplyTransformations(bundle.Artifacts, artifactsTransformations)
	if err != nil {
		return bundle, err
	}

	hydratedParams, err := jsonschema.Hydrate(bundle.Params, cwd)
	if err != nil {
		return bundle, err
	}
	bundle.Params = hydratedParams.(map[string]interface{})
	err = ApplyTransformations(bundle.Params, paramsTransformations)
	if err != nil {
		return bundle, err
	}

	hydratedConnections, err := jsonschema.Hydrate(bundle.Connections, cwd)
	if err != nil {
		return bundle, err
	}
	bundle.Connections = hydratedConnections.(map[string]interface{})
	err = ApplyTransformations(bundle.Connections, connectionsTransformations)
	if err != nil {
		return bundle, err
	}

	hydratedUi, err := jsonschema.Hydrate(bundle.Ui, cwd)
	if err != nil {
		return bundle, err
	}
	bundle.Ui = hydratedUi.(map[string]interface{})
	err = ApplyTransformations(bundle.Ui, uiTransformations)
	if err != nil {
		return bundle, err
	}

	return bundle, nil
}

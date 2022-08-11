package application

import (
	"io"
	"os"
	"path"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

// TODO: dedupe w/ build
func Package(appPath string, c *client.MassdriverClient, workingDir string, buf io.Writer) (*Application, error) {
	app, parseErr := Parse(appPath, nil)
	if parseErr != nil {
		return nil, parseErr
	}

	// I yanked this at first, but do we need it?
	// Or, we can just write the "app spec" as the "bundle.yaml"
	// it would have all the app fields, and be type: application
	bytes, errMarshal := yaml.Marshal(*app)
	errWrite := common.WriteFile(path.Join(workingDir, "app.yaml"), bytes, errMarshal)
	if errWrite != nil {
		return nil, errWrite
	}

	b := app.AsBundle()

	// We're using bundle.yaml instead of massdriver.yaml here so we don't overwrite the application config
	sourceDir := path.Dir(appPath)
	bundlePath := path.Join(workingDir, "bundle.yaml")
	bytesB, errMarshalB := yaml.Marshal(*b)
	errWriteB := common.WriteFile(bundlePath, bytesB, errMarshalB)
	if errWriteB != nil {
		return nil, errWriteB
	}

	steps := app.Steps
	if app.Steps == nil {
		steps = []bundle.Step{
			{
				Path:        "src",
				Provisioner: "terraform",
			},
		}
	}

	// COPY FILES.
	// We need this for any modules, charts, etc included in each step
	for _, step := range steps {
		log.Info().Msgf("Copying files for step %s", step.Path)
		err := os.MkdirAll(path.Join(workingDir, step.Path), 0744)
		if err != nil {
			return nil, err
		}
		// ignore these, copy the rest
		ignores := []string{
			".terraform",
			"terraform.tfstate",
			"auto.tfvars.json",
			"connections.auto.tfvars.json",
			"dev.connections.tfvars",
			"dev.params.tfvars",
			"_connections_variables.tf.json",
			"_md_variables.tf.json",
			"_params_variables.tf.json",
		}
		errCopy := common.CopyFolder(path.Join(sourceDir, step.Path), path.Join(workingDir, step.Path), ignores)
		if errCopy != nil {
			return nil, errCopy
		}
	}

	if errBuild := app.Build(c, workingDir); errBuild != nil {
		return nil, errBuild
	}

	err := bundle.PackageBundle(workingDir, buf)
	if err != nil {
		return nil, err
	}

	return app, nil
}

package application

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"github.com/massdriver-cloud/massdriver-cli/pkg/provisioners/terraform"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

func Package(appPath string, c *client.MassdriverClient, workingDir string, buf io.Writer) (*bundle.Bundle, error) {
	app, parseErr := Parse(appPath)
	if parseErr != nil {
		return nil, parseErr
	}

	bytes, errMarshal := yaml.Marshal(*app)
	errWrite := common.WriteFile(path.Join(workingDir, "app.yaml"), bytes, errMarshal)
	if errWrite != nil {
		return nil, errWrite
	}

	// Write bundle.yaml
	b, err := app.ConvertToBundle()
	if err != nil {
		return nil, fmt.Errorf("could not convert app to bundle: %w", err)
	}
	// We're using bundle.yaml instead of massdriver.yaml here so we don't overwrite the application config
	bundlePath := path.Join(workingDir, "bundle.yaml")
	bytesB, errMarshalB := yaml.Marshal(*b)
	errWriteB := common.WriteFile(bundlePath, bytesB, errMarshalB)
	if errWriteB != nil {
		return nil, errWriteB
	}

	steps := b.Steps
	if b.Steps == nil {
		steps = []bundle.Step{
			{
				Path:        "src",
				Provisioner: "terraform",
			},
		}
	}

	err = b.Hydrate(bundlePath, c)
	if err != nil {
		return nil, err
	}

	err = b.GenerateSchemas(workingDir)
	if err != nil {
		return nil, err
	}
	appDir := filepath.Dir(appPath)
	bundleDir := filepath.Dir(bundlePath)
	// Make all directories, generate provisioner-specific files
	for _, step := range steps {
		err = os.MkdirAll(path.Join(workingDir, step.Path), 0744)
		if err != nil {
			return nil, err
		}
		log.Debug().Msgf("copy from: %s", path.Join(appDir, step.Path))
		log.Debug().Msgf("copy to: %s", path.Join(bundleDir, step.Path))
		errCopy := common.CopyFolder(path.Join(appDir, step.Path), path.Join(bundleDir, step.Path))
		if errCopy != nil {
			return nil, errCopy
		}

		if stepErr := generateStep(step, workingDir, bundlePath); stepErr != nil {
			return nil, stepErr
		}
	}

	err = bundle.PackageBundle(bundlePath, buf)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func generateStep(step bundle.Step, workingDir, bundlePath string) error {
	switch step.Provisioner {
	case "terraform":
		err := terraform.GenerateFiles(workingDir, step.Path)
		if err != nil {
			log.Error().Err(err).Str("bundle", bundlePath).Str("provisioner", step.Provisioner).Msg("an error occurred while generating provisioner files")
			return err
		}
	case "exec":
		// No-op (Golang doesn't not fallthrough unless explicitly stated)
	default:
		log.Error().Str("bundle", bundlePath).Msg("unknown provisioner: " + step.Provisioner)
		return fmt.Errorf("unknown provisioner: %v", step.Provisioner)
	}
	return nil
}

package application

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"github.com/massdriver-cloud/massdriver-cli/pkg/provisioners/terraform"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

func PackageApplication(appPath string, c *client.MassdriverClient, workingDir string, buf io.Writer) (*bundle.Bundle, error) {
	app, parseErr := Parse(appPath)
	if parseErr != nil {
		return nil, parseErr
	}

	// Write app.yaml
	appYaml, err := os.Create(path.Join(workingDir, "app.yaml"))
	if err != nil {
		return nil, err
	}
	defer appYaml.Close()
	appYamlBytes, marshalErr := yaml.Marshal(*app)
	if marshalErr != nil {
		return nil, marshalErr
	}

	if _, yamlErr := appYaml.Write(appYamlBytes); yamlErr != nil {
		return nil, yamlErr
	}

	// Write bundle.yaml
	b, err := app.ConvertToBundle()
	if err != nil {
		return nil, fmt.Errorf("could not convert app to bundle: %w", err)
	}
	// We're using bundle.yaml instead of massdriver.yaml here so we don't overwrite the application config
	bundlePath := path.Join(workingDir, "bundle.yaml")
	bundleYaml, err := os.Create(bundlePath)
	if err != nil {
		return nil, err
	}
	defer bundleYaml.Close()
	bundleYamlBytes, marshalErr := yaml.Marshal(*b)
	if marshalErr != nil {
		return nil, marshalErr
	}
	if _, writeErr := bundleYaml.Write(bundleYamlBytes); writeErr != nil {
		return nil, writeErr
	}

	// TODO: move chart into src dir
	if app.Deployment.Type == "custom" {
		// Make chart directory
		err = os.MkdirAll(path.Join(workingDir, "chart"), 0744)
		if err != nil {
			return nil, err
		}
		err = packageChart(path.Join(path.Dir(appPath), app.Deployment.Path), path.Join(workingDir, "chart"))
		if err != nil {
			return nil, err
		}
	}

	// Make src directory
	err = os.MkdirAll(path.Join(workingDir, "src"), 0744)
	if err != nil {
		return nil, err
	}

	err = b.Hydrate(bundlePath, c)
	if err != nil {
		return nil, err
	}

	err = b.GenerateSchemas(workingDir)
	if err != nil {
		return nil, err
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

	for _, step := range steps {
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

func packageChart(chartPath string, destPath string) error {
	err := filepath.Walk(chartPath, func(path string, info os.FileInfo, err error) error {
		var relPath = strings.TrimPrefix(path, chartPath)
		if relPath == "" {
			return nil
		}
		if info.IsDir() {
			return os.Mkdir(filepath.Join(destPath, relPath), common.AllRX|common.UserRW)
		}
		var data, err1 = ioutil.ReadFile(filepath.Join(chartPath, relPath))
		if err1 != nil {
			return err1
		}
		return ioutil.WriteFile(filepath.Join(destPath, relPath), data, common.AllRWX)
	})
	return err
}

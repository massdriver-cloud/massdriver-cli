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
	"github.com/massdriver-cloud/massdriver-cli/pkg/provisioners/terraform"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

func PackageApplication(appPath string, c *client.MassdriverClient, workingDir string, buf io.Writer) (*bundle.Bundle, error) {
	app, err := Parse(appPath)
	if err != nil {
		return nil, err
	}

	// Write app.yaml
	appYaml, err := os.Create(path.Join(workingDir, "app.yaml"))
	if err != nil {
		return nil, err
	}
	defer appYaml.Close()
	appYamlBytes, err := yaml.Marshal(*app)
	if err != nil {
		return nil, err
	}
	appYaml.Write(appYamlBytes)

	// Write bundle.yaml
	b := app.ConvertToBundle()
	// We're using bundle.yaml instead of massdriver.yaml here so we don't overwrite the application
	bundlePath := path.Join(workingDir, "bundle.yaml")
	bundleYaml, err := os.Create(bundlePath)
	if err != nil {
		return nil, err
	}
	defer bundleYaml.Close()
	bundleYamlBytes, err := yaml.Marshal(*b)
	if err != nil {
		return nil, err
	}
	bundleYaml.Write(bundleYamlBytes)

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

	for _, step := range b.Steps {
		switch step.Provisioner {
		case "terraform":
			err = terraform.GenerateFiles(workingDir, step.Path)
			if err != nil {
				log.Error().Err(err).Str("bundle", bundlePath).Str("provisioner", step.Provisioner).Msg("an error occurred while generating provisioner files")
				return nil, err
			}
		case "exec":
			// No-op (Golang doesn't not fallthrough unless explicitly stated)
		default:
			log.Error().Str("bundle", bundlePath).Msg("unknown provisioner: " + step.Provisioner)
			return nil, fmt.Errorf("unknown provisioner: %v", step.Provisioner)
		}
	}

	err = bundle.PackageBundle(bundlePath, buf)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func packageChart(chartPath string, destPath string) error {
	var err error = filepath.Walk(chartPath, func(path string, info os.FileInfo, err error) error {
		var relPath string = strings.TrimPrefix(path, chartPath)
		if relPath == "" {
			return nil
		}
		if info.IsDir() {
			return os.Mkdir(filepath.Join(destPath, relPath), 0755)
		} else {
			var data, err1 = ioutil.ReadFile(filepath.Join(chartPath, relPath))
			if err1 != nil {
				return err1
			}
			return ioutil.WriteFile(filepath.Join(destPath, relPath), data, 0777)
		}
	})
	return err
}

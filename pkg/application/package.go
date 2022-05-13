package application

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
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
	return b, nil
}

func PackageBak(filePath string, buf io.Writer) error {
	app, err := Parse(filePath)
	if err != nil {
		return err
	}
	// tar > gzip > buf
	gzipWriter := gzip.NewWriter(buf)
	tarWriter := tar.NewWriter(gzipWriter)

	// absolutePath, err := filepath.Abs(filePath)
	// if err != nil {
	// 	return err
	// }
	//dirPath := path.Dir(absolutePath)

	// APP.YAML
	appYamlFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer appYamlFile.Close()
	appYamlFileInfo, err := appYamlFile.Stat()
	if err != nil {
		return err
	}
	appYamlHeader, err := tar.FileInfoHeader(appYamlFileInfo, filePath)
	if err != nil {
		return err
	}
	appYamlHeader.Name = "bundle/app.yaml"
	if err := tarWriter.WriteHeader(appYamlHeader); err != nil {
		return err
	}
	if _, err := io.Copy(tarWriter, appYamlFile); err != nil {
		return err
	}

	// BUNDLE.YAML
	b := app.ConvertToBundle()
	bundleBytes, err := yaml.Marshal(b)
	if err != nil {
		return err
	}
	bundleYamlHeader := &tar.Header{
		Name: "bundle/bundle.yaml",
		Mode: 0600,
		Size: int64(len(bundleBytes)),
	}
	if err := tarWriter.WriteHeader(bundleYamlHeader); err != nil {
		return err
	}
	if _, err := io.Copy(tarWriter, bytes.NewReader(bundleBytes)); err != nil {
		return err
	}

	// SCHEMAS
	schemas := map[string]map[string]interface{}{
		// "artifacts": b.Artifacts,
		// "connections": b.Connections,
		"params": b.Params,
	}
	for name, schema := range schemas {
		var buffer bytes.Buffer
		if err = bundle.GenerateSchema(schema, b.Metadata(name), &buffer); err != nil {
			return err
		}
		schemaHeader := &tar.Header{
			Name: "bundle/schema-" + name + ".json",
			Mode: 0600,
			Size: int64(len(buffer.Bytes())),
		}
		if err := tarWriter.WriteHeader(schemaHeader); err != nil {
			return err
		}
		if _, err := io.Copy(tarWriter, bytes.NewReader(buffer.Bytes())); err != nil {
			return err
		}
	}

	// // walk through every file in the folder
	// filepath.Walk(dirPath, func(file string, fi os.FileInfo, err error) error {
	// 	// generate tar header
	// 	header, err := tar.FileInfoHeader(fi, file)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// must provide real name
	// 	// (see https://golang.org/src/archive/tar/common.go?#L626)
	// 	header.Name = path.Join("bundle", strings.TrimPrefix(filepath.ToSlash(file), dirPath))

	// 	// write header
	// 	if err := tarWriter.WriteHeader(header); err != nil {
	// 		return err
	// 	}
	// 	// if not a dir, write file content
	// 	if !fi.IsDir() {
	// 		data, err := os.Open(file)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		if _, err := io.Copy(tarWriter, data); err != nil {
	// 			return err
	// 		}
	// 		data.Close()
	// 	}
	// 	return nil
	// })

	// produce tar
	if err := tarWriter.Close(); err != nil {
		return err
	}
	// produce gzip
	if err := gzipWriter.Close(); err != nil {
		return err
	}

	return nil
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

package application

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"gopkg.in/yaml.v2"
)

func Package(filePath string, buf io.Writer) error {
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

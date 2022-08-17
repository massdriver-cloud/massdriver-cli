package application

import (
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"github.com/rs/zerolog/log"
)

// TODO: dedupe w/ build
// TODO: dedupe w/ bundle Package
// TODO: don't write to disk by default, write to buffer
func Package(massYamlPath string, c *client.MassdriverClient, destinationDir string, buf io.Writer) (*Application, error) {
	// since we don't do any app / bundle yaml transforms
	// we can just copy the entire directory to the destination
	// then run app build, bundle package
	sourceDir := path.Dir(massYamlPath)
	files, errReadDir := ioutil.ReadDir(sourceDir)
	if errReadDir != nil {
		return nil, errReadDir
	}

	for _, fileInfo := range files {
		fileName := fileInfo.Name()
		pathCopyFrom := path.Join(sourceDir, fileName)
		pathCopyTo := path.Join(destinationDir, fileName)

		log.Debug().Msgf("Packaging %s", fileName)
		if fileInfo.IsDir() {
			errMkdir := os.MkdirAll(pathCopyTo, 0744)
			if errMkdir != nil {
				return nil, errMkdir
			}
			errCopy := common.CopyFolder(pathCopyFrom, pathCopyTo, common.FileIgnores)
			if errCopy != nil {
				return nil, errCopy
			}
		} else {
			fileBytes, errRead := ioutil.ReadFile(pathCopyFrom)
			errWrite := common.WriteFile(pathCopyTo, fileBytes, errRead)
			if errWrite != nil {
				return nil, errWrite
			}
		}
	}

	app, errParse := Parse(destinationDir+"/"+common.MassdriverYamlFilename, nil)
	if errParse != nil {
		return nil, errParse
	}

	if errBuild := app.Build(c, destinationDir); errBuild != nil {
		return nil, errBuild
	}

	errPackage := bundle.PackageBundle(destinationDir+"/"+common.MassdriverYamlFilename, buf)
	if errPackage != nil {
		return nil, errPackage
	}

	return app, nil
}

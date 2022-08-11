package application

import (
	"bytes"
	"os"
	"path"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
)

func Publish(c *client.MassdriverClient, appPath string) error {
	workingDir, err := os.MkdirTemp("", "application")
	if err != nil {
		return err
	}
	defer os.RemoveAll(workingDir)

	var buf bytes.Buffer
	_, errPackage := Package(appPath, c, workingDir, &buf)
	if errPackage != nil {
		return errPackage
	}

	// hack, resolve app / bundle publish
	b, parseErr := bundle.Parse(path.Join(workingDir, "bundle.yaml"), nil)
	if parseErr != nil {
		return parseErr
	}

	uploadURL, err := b.Publish(c)
	if err != nil {
		return err
	}

	return bundle.UploadToPresignedS3URL(uploadURL, &buf)
}

package application

import (
	"bytes"
	"os"

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
	b, err := PackageBetter(appPath, c, workingDir, &buf)
	if err != nil {
		return err
	}

	uploadURL, err := b.Publish(c)
	if err != nil {
		return err
	}

	return bundle.UploadToPresignedS3URL(uploadURL, &buf)
}

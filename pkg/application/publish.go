package application

import (
	"bytes"
	"os"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
	"github.com/rs/zerolog/log"
)

func Publish(appPath string, c *client.MassdriverClient) error {
	workingDir, err := os.MkdirTemp("", "application")
	if err != nil {
		return (err)
	}
	defer os.RemoveAll(workingDir)

	var buf bytes.Buffer
	b, err := PackageApplication(appPath, c, workingDir, &buf)
	if err != nil {
		return err
	}

	uploadURL, err := b.Publish(c)
	if err != nil {
		return err
	}

	log.Info().Msg("Application published successfully!")

	return bundle.UploadToPresignedS3URL(uploadURL, &buf)
}

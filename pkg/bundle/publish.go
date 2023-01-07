package bundle

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
)

func Publish(c *client.MassdriverClient, b *Bundle) error {
	if errBuild := b.Build(c, path.Dir(common.MassdriverYamlFilename)); errBuild != nil {
		return errBuild
	}

	var buf bytes.Buffer
	errPackage := Package(b, common.MassdriverYamlFilename, &buf)
	if errPackage != nil {
		return errPackage
	}

	uploadURL, err := b.PublishToMassdriver(path.Dir(common.MassdriverYamlFilename), c)
	if err != nil {
		return err
	}

	return UploadToPresignedS3URL(uploadURL, &buf)
}

func (b *Bundle) PublishToMassdriver(srcDir string, c *client.MassdriverClient) (string, error) {
	body := b.generateBundlePublishBody()

	err := checkForOperatorGuideAndSetValue(srcDir, &body)

	if err != nil {
		return "", err
	}

	bodyBytes, err := json.Marshal(body)

	if err != nil {
		return "", err
	}

	ctx := context.Background()
	req := client.NewRequest("PUT", "bundles", bytes.NewBuffer(bodyBytes))
	resp, err := c.Do(&ctx, req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.Status != "200 OK" {
		fmt.Println(string(respBodyBytes))
		return "", errors.New("received non-200 response from Massdriver: " + resp.Status)
	}

	var respBody PublishResponse
	err = json.Unmarshal(respBodyBytes, &respBody)
	if err != nil {
		return "", err
	}

	return respBody.UploadLocation, nil
}

func (b *Bundle) generateBundlePublishBody() PublishPost {
	var body PublishPost

	body.Name = b.Name
	body.Description = b.Description
	body.Type = b.Type
	body.SourceURL = b.SourceURL
	body.Access = b.Access
	body.ArtifactsSchema = b.Artifacts
	body.ConnectionsSchema = b.Connections
	body.ParamsSchema = b.Params
	body.UISchema = b.UI
	body.AppSpec = b.App

	return body
}

func UploadToPresignedS3URL(url string, object io.Reader) error {
	req, err := http.NewRequestWithContext(context.Background(), "PUT", url, object)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		var respContent S3PresignEndpointResponse
		var respBody bytes.Buffer
		if _, readErr := respBody.ReadFrom(resp.Body); readErr != nil {
			return readErr
		}
		if xmlErr := xml.Unmarshal(respBody.Bytes(), &respContent); xmlErr != nil {
			return fmt.Errorf("enountered non-200 response code, unable to unmarshal xml response body: %v: original error: %w", respBody.String(), xmlErr)
		}

		return errors.New("unable to upload content: " + respContent.Message)
	}
	return nil
}

func checkForOperatorGuideAndSetValue(path string, body *PublishPost) error {
	pathsToCheck := []string{"operator.mdx", "operator.md"}

	for _, fileName := range pathsToCheck {
		_, err := os.Stat(filepath.Join(path, fileName))

		if err != nil {
			continue
		}

		content, err := os.ReadFile(filepath.Join(path, fileName))

		if err != nil {
			return fmt.Errorf("error reading %s", fileName)
		}

		body.OperatorGuide = content
	}

	return nil
}

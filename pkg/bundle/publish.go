package bundle

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
)

var MASSDRIVER_URL string = "https://api.massdriver.cloud/"

func (b *Bundle) Publish(c *client.MassdriverClient) (string, error) {

	body, err := b.generateBundlePublishBody()
	if err != nil {
		return "", err
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req := client.NewRequest("PUT", "bundles", bytes.NewBuffer(bodyBytes))
	resp, err := c.Do(req)
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

	var respBody BundlePublishResponse
	err = json.Unmarshal(respBodyBytes, &respBody)
	if err != nil {
		return "", err
	}

	return respBody.UploadLocation, nil
}

func (b *Bundle) generateBundlePublishBody() (BundlePublishPost, error) {
	var body BundlePublishPost

	body.Name = b.Name
	body.Description = b.Description
	body.Type = "bundle"
	body.Ref = b.Ref
	body.Access = b.Access

	artifactsSchema, err := json.Marshal(b.Artifacts)
	if err != nil {
		return body, err
	}
	body.ArtifactsSchema = string(artifactsSchema)

	connectionsSchema, err := json.Marshal(b.Connections)
	if err != nil {
		return body, err
	}
	body.ConnectionsSchema = string(connectionsSchema)

	paramsSchema, err := json.Marshal(b.Params)
	if err != nil {
		return body, err
	}
	body.ParamsSchema = string(paramsSchema)

	uiSchema, err := json.Marshal(b.Ui)
	if err != nil {
		return body, err
	}
	body.UISchema = string(uiSchema)

	return body, nil
}

func UploadToPresignedS3URL(url string, object io.Reader) error {
	req, err := http.NewRequest("PUT", url, object)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode != 200 {
		var respContent S3PresignEndpointResponse
		var respBody bytes.Buffer
		respBody.ReadFrom(resp.Body)
		xml.Unmarshal(respBody.Bytes(), &respContent)
		return errors.New("unable to upload content: " + respContent.Message)
	}
	return nil
}

package bundle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var MASSDRIVER_URL string = "https://api.massdriver.cloud/"

type BundlePost struct {
	Name              string                 `json:"name"`
	Ref               string                 `json:"ref"`
	ID                string                 `json:"id"`
	Access            string                 `json:"access"`
	ArtifactsSchema   map[string]interface{} `json:"artifacts_schema"`
	ConnectionsSchema map[string]interface{} `json:"connections_schema"`
	ParamsSchema      map[string]interface{} `json:"params_schema"`
	UISchema          map[string]interface{} `json:"ui_schema"`
}

func (b Bundle) Publish(apiKey string) error {
	url, err := url.Parse(MASSDRIVER_URL)
	if err != nil {
		return err
	}
	url.Path = "bundles"

	body, err := json.Marshal(b.generateBundlePublishBody())
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("X-Md-Api-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(respBody))

	return nil
}

func (b Bundle) generateBundlePublishBody() BundlePost {
	var body BundlePost

	body.Name = b.Title
	body.Ref = b.Type
	body.ID = b.Uuid
	body.Access = b.Access
	body.ArtifactsSchema = b.Artifacts
	body.ConnectionsSchema = b.Connections
	body.ParamsSchema = b.Params
	body.UISchema = b.Ui

	return body
}

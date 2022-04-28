package bundle

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var MASSDRIVER_URL string = "https://api.massdriver.cloud/"

type BundlePublishPost struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	Type              string `json:"type"`
	Ref               string `json:"ref"`
	Access            string `json:"access"`
	ArtifactsSchema   string `json:"artifacts_schema"`
	ConnectionsSchema string `json:"connections_schema"`
	ParamsSchema      string `json:"params_schema"`
	UISchema          string `json:"ui_schema"`
}

type BundlePublishResponse struct {
	UploadLocation string `json:"upload_location"`
}

type S3PresignEndpointResponse struct {
	Error                 xml.Name `xml:"Error"`
	Code                  string   `xml:"Code"`
	Message               string   `xml:"Message"`
	AWSAccessKeyId        string   `xml:"AWSAccessKeyId"`
	StringToSign          string   `xml:"StringToSign"`
	SignatureProvided     string   `xml:"SignatureProvided"`
	StringToSignBytes     []byte   `xml:"StringToSignBytes"`
	CanonicalRequest      string   `xml:"CanonicalRequest"`
	CanonicalRequestBytes []byte   `xml:"CanonicalRequestBytes"`
	RequestId             string   `xml:"RequestId"`
	HostId                string   `xml:"HostId"`
}

type S3PresignEndpointResponseError struct {
	Code                  string `xml:"Code"`
	Message               string `xml:"Message"`
	AWSAccessKeyId        string `xml:"AWSAccessKeyId"`
	StringToSign          string `xml:"StringToSign"`
	SignatureProvided     string `xml:"SignatureProvided"`
	StringToSignBytes     []byte `xml:"StringToSignBytes"`
	CanonicalRequest      string `xml:"CanonicalRequest"`
	CanonicalRequestBytes []byte `xml:"CanonicalRequestBytes"`
	RequestId             string `xml:"RequestId"`
	HostId                string `xml:"HostId"`
}

func (b Bundle) Publish(apiKey string) (string, error) {
	mdUrl, err := url.Parse(MASSDRIVER_URL)
	if err != nil {
		return "", err
	}
	mdUrl.Path = "bundles"

	body, err := b.generateBundlePublishBody()
	if err != nil {
		return "", err
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("PUT", mdUrl.String(), bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Md-Api-Key", apiKey)
	req.Header.Set("content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
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

func (b Bundle) generateBundlePublishBody() (BundlePublishPost, error) {
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

func TarGzipBundle(filePath string, buf io.Writer) error {
	// tar > gzip > buf
	gzipWriter := gzip.NewWriter(buf)
	tarWriter := tar.NewWriter(gzipWriter)

	absolutePath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}
	dirPath := path.Dir(absolutePath)
	prefix := path.Dir(dirPath) + "/"

	// walk through every file in the folder
	filepath.Walk(dirPath, func(file string, fi os.FileInfo, err error) error {
		// generate tar header
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		// must provide real name
		// (see https://golang.org/src/archive/tar/common.go?#L626)
		header.Name = strings.TrimPrefix(filepath.ToSlash(file), prefix)

		// write header
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}
		// if not a dir, write file content
		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tarWriter, data); err != nil {
				return err
			}
		}
		return nil
	})

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

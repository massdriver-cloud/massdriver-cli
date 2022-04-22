package bundle

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
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
	Name              string                 `json:"name"`
	Ref               string                 `json:"ref"`
	ID                string                 `json:"id"`
	Access            string                 `json:"access"`
	ArtifactsSchema   map[string]interface{} `json:"artifacts_schema"`
	ConnectionsSchema map[string]interface{} `json:"connections_schema"`
	ParamsSchema      map[string]interface{} `json:"params_schema"`
	UISchema          map[string]interface{} `json:"ui_schema"`
}

type BundlePublishResponse struct {
	UploadLocation string `json:"upload_location"`
}

func (b Bundle) Publish(apiKey string) (string, error) {
	url, err := url.Parse(MASSDRIVER_URL)
	if err != nil {
		return "", err
	}
	url.Path = "bundles"

	body, err := json.Marshal(b.generateBundlePublishBody())
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Md-Api-Key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return "", errors.New("received non-200 response from Massdriver: " + resp.Status)
	}

	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var respBody BundlePublishResponse
	err = json.Unmarshal(respBodyBytes, &respBody)
	if err != nil {
		return "", err
	}

	fmt.Printf("response Body: %v", respBody)

	return respBody.UploadLocation, nil
}

func (b Bundle) generateBundlePublishBody() BundlePublishPost {
	var body BundlePublishPost

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

func UploadBytesToHTTPEndpoint(url string, object io.Reader) error {
	// Use the presigned URL to put a object to S3
	req, err := http.NewRequest("PUT", url, object)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//
	// Print out the response.
	fmt.Println("Status", resp.StatusCode, resp.StatusCode)
	o := &bytes.Buffer{}
	io.Copy(o, resp.Body)
	fmt.Println(o.String())
	return nil
}

func TarGzipDirectory(filePath string, buf io.Writer) error {
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

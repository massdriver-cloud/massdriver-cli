package bundle

import "encoding/xml"

type Step struct {
	Path        string `json:"path" yaml:"path"`
	Provisioner string `json:"provisioner" yaml:"provisioner"`
}

type Bundle struct {
	Schema      string                 `json:"schema" yaml:"schema"`
	Name        string                 `json:"name" yaml:"name"`
	Description string                 `json:"description" yaml:"description"`
	Ref         string                 `json:"ref" yaml:"ref"`
	Type        string                 `json:"type" yaml:"type"`
	Access      string                 `json:"access" yaml:"access"`
	Steps       []Step                 `json:"steps" yaml:"steps"`
	Artifacts   map[string]interface{} `json:"artifacts" yaml:"artifacts"`
	Params      map[string]interface{} `json:"params" yaml:"params"`
	Connections map[string]interface{} `json:"connections" yaml:"connections"`
	UI          map[string]interface{} `json:"ui" yaml:"ui"`
}

type PublishPost struct {
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

type PublishResponse struct {
	UploadLocation string `json:"upload_location"`
}

type S3PresignEndpointResponse struct {
	Error                 xml.Name `xml:"Error"`
	Code                  string   `xml:"Code"`
	Message               string   `xml:"Message"`
	AWSAccessKeyID        string   `xml:"AWSAccessKeyId"`
	StringToSign          string   `xml:"StringToSign"`
	SignatureProvided     string   `xml:"SignatureProvided"`
	StringToSignBytes     []byte   `xml:"StringToSignBytes"`
	CanonicalRequest      string   `xml:"CanonicalRequest"`
	CanonicalRequestBytes []byte   `xml:"CanonicalRequestBytes"`
	RequestID             string   `xml:"RequestId"`
	HostID                string   `xml:"HostId"`
}

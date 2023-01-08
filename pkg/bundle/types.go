package bundle

import (
	"encoding/xml"
)

type Step struct {
	Path        string `json:"path" yaml:"path"`
	Provisioner string `json:"provisioner" yaml:"provisioner"`
}

type Bundle struct {
	Schema      string                 `json:"schema" yaml:"schema"`
	Name        string                 `json:"name" yaml:"name"`
	Description string                 `json:"description" yaml:"description"`
	SourceURL   string                 `json:"source_url" yaml:"source_url"`
	Type        string                 `json:"type" yaml:"type"`
	Access      string                 `json:"access" yaml:"access"`
	Steps       []Step                 `json:"steps" yaml:"steps"`
	Artifacts   map[string]interface{} `json:"artifacts" yaml:"artifacts"`
	Params      map[string]interface{} `json:"params" yaml:"params"`
	Connections map[string]interface{} `json:"connections" yaml:"connections"`
	UI          map[string]interface{} `json:"ui" yaml:"ui"`
	App         *AppBlock              `json:"app" yaml:"app"`
}

type AppBlock struct {
	Envs     map[string]string `json:"envs" yaml:"envs"`
	Policies []string          `json:"policies" yaml:"policies"`
	Secrets  map[string]Secret `json:"secrets" yaml:"secrets"`
}

type Secret struct {
	Required    bool   `json:"required" yaml:"required"`
	Json        bool   `json:"json" yaml:"json"`
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
}

type PublishPost struct {
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	Type              string                 `json:"type"`
	SourceURL         string                 `json:"source_url"`
	Access            string                 `json:"access"`
	ArtifactsSchema   map[string]interface{} `json:"artifacts_schema"`
	ConnectionsSchema map[string]interface{} `json:"connections_schema"`
	ParamsSchema      map[string]interface{} `json:"params_schema"`
	UISchema          map[string]interface{} `json:"ui_schema"`
	OperatorGuide     []byte                 `json:"operator_guide,omitempty"`
	AppSpec           *AppBlock              `json:"app,omitempty"`
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

func (b *Bundle) IsInfrastructure() bool {
	// a Deprecation warning is printed in the bundle parse function
	return b.Type == "bundle" || b.Type == "infrastructure"
}

func (b *Bundle) IsApplication() bool {
	return b.Type == "application"
}

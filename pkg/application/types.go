package application

import "github.com/massdriver-cloud/massdriver-cli/pkg/bundle"

type Application struct {
	Schema      string `json:"schema" yaml:"schema"`
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Ref         string `json:"ref" yaml:"ref"`
	Type        string `json:"type" yaml:"type"`
	Access      string `json:"access" yaml:"access"`
	// TODO: deprecate
	Deployment Deployment             `json:"deployment" yaml:"deployment"`
	Steps      []bundle.Step          `json:"steps" yaml:"steps"`
	Params     map[string]interface{} `json:"params" yaml:"params"`
	// TODO: deprecate
	Dependencies map[string]Dependencies `json:"dependencies" yaml:"dependencies"`
	Connections  map[string]interface{}  `json:"connections" yaml:"connections"`
	UI           map[string]interface{}  `json:"ui" yaml:"ui"`
	App          AppBlock                `json:"app" yaml:"app"`
}

type AppBlock struct {
	Envs     map[string]string `json:"envs" yaml:"envs"`
	Policies []string          `json:"policies" yaml:"policies"`
}

// TODO: deprecate
type Deployment struct {
	Type       string `json:"type" yaml:"type"`
	Path       string `json:"path,omitempty" yaml:"path,omitempty"`
	Chart      string `json:"chart,omitempty" yaml:"chart,omitempty"`
	Repository string `json:"repository,omitempty" yaml:"repository,omitempty"`
}

// TODO: deprecate
type Dependencies struct {
	Type     string             `json:"type" yaml:"type"`
	Required bool               `json:"required,omitempty" yaml:"required,omitempty"`
	Envs     []DependenciesEnvs `json:"envs" yaml:"envs"`
	Policies []string           `json:"policies,omitempty" yaml:"policies,omitempty"`
}

// TODO: deprecate
type DependenciesEnvs struct {
	Name string `json:"name" yaml:"name"`
	Path string `json:"path" yaml:"path"`
}

// TODO: deprecate
type ChartYAML struct {
	APIVersion  string `yaml:"apiVersion"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Type        string `yaml:"type"`
	Version     string `yaml:"version"`
}

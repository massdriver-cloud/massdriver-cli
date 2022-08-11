package application

import "github.com/massdriver-cloud/massdriver-cli/pkg/bundle"

type Application struct {
	Schema      string                 `json:"schema" yaml:"schema"`
	Name        string                 `json:"name" yaml:"name"`
	Description string                 `json:"description" yaml:"description"`
	Ref         string                 `json:"ref" yaml:"ref"`
	Type        string                 `json:"type" yaml:"type"`
	Access      string                 `json:"access" yaml:"access"`
	Steps       []bundle.Step          `json:"steps" yaml:"steps"`
	Params      map[string]interface{} `json:"params" yaml:"params"`
	Connections map[string]interface{} `json:"connections" yaml:"connections"`
	UI          map[string]interface{} `json:"ui" yaml:"ui"`
	App         AppBlock               `json:"app" yaml:"app"`
}

type AppBlock struct {
	Envs     map[string]string `json:"envs" yaml:"envs"`
	Policies []string          `json:"policies" yaml:"policies"`
}

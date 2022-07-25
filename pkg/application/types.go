package application

type TemplateData struct {
	// kubernetes-deployment, kubernetes-job
	TemplateName string
	// application name
	Name        string
	Description string
	Access      string
	OutputDir   string
}

type Application struct {
	Schema       string                  `json:"schema" yaml:"schema"`
	Title        string                  `json:"title" yaml:"title"`
	Description  string                  `json:"description" yaml:"description"`
	Ref          string                  `json:"ref" yaml:"ref"`
	Access       string                  `json:"access" yaml:"access"`
	Metadata     Metadata                `json:"metadata" yaml:"metadata"`
	Params       map[string]interface{}  `json:"params" yaml:"params"`
	Dependencies map[string]Dependencies `json:"dependencies" yaml:"dependencies"`
	UI           map[string]interface{}  `json:"ui" yaml:"ui"`
}

type Metadata struct {
	Template           string `json:"template" yaml:"template"`
	TemplateRepository string `json:"repository,omitempty" yaml:"repository,omitempty"`
	Path               string `json:"path,omitempty" yaml:"path,omitempty"`
}

type Dependencies struct {
	Type     string             `json:"type" yaml:"type"`
	Required bool               `json:"required,omitempty" yaml:"required,omitempty"`
	Envs     []DependenciesEnvs `json:"envs" yaml:"envs"`
	Policies []string           `json:"policies,omitempty" yaml:"policies,omitempty"`
}

type DependenciesEnvs struct {
	Name string `json:"name" yaml:"name"`
	Path string `json:"path" yaml:"path"`
}

type ChartYAML struct {
	APIVersion  string `yaml:"apiVersion"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Type        string `yaml:"type"`
	Version     string `yaml:"version"`
}

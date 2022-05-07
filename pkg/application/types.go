package application

type Application struct {
	Schema       string                    `json:"schema" yaml:"schema"`
	Name         string                    `json:"name" yaml:"name"`
	Description  string                    `json:"description" yaml:"description"`
	Ref          string                    `json:"ref" yaml:"ref"`
	Access       string                    `json:"access" yaml:"access"`
	Deployment   ApplicationDeployment     `json:"deployment" yaml:"deployment"`
	Params       map[string]interface{}    `json:"params" yaml:"params"`
	Dependencies []ApplicationDependencies `json:"dependencies" yaml:"dependencies"`
}

type ApplicationDeployment struct {
	Type  string `json:"type" yaml:"type"`
	Chart string `json:"chart" yaml:"chart"`
}

type ApplicationDependencies struct {
	Type     string                        `json:"type" yaml:"type"`
	Field    string                        `json:"field" yaml:"field"`
	Required bool                          `json:"required" yaml:"required,omitempty"`
	Env      []ApplicationDependenciesEnvs `json:"env" yaml:"env"`
	Policy   string                        `json:"policy" yaml:"policy"`
}

type ApplicationDependenciesEnvs struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

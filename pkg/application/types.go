package application

type Application struct {
	Bundle       string
	Title        string
	Description  string
	Deployment   ApplicationDeployment
	Params       map[string]interface{}
	Dependencies []ApplicationDependencies
}

type ApplicationDeployment struct {
	Type string
}

type ApplicationParams struct {
}

type ApplicationDependencies struct {
	Type     string
	Field    string
	Required *bool `yaml:"required,omitempty"`
	Env      []ApplicationDependenciesEnvs
	Policy   string
}

type ApplicationDependenciesEnvs struct {
	Name  string ""
	Value string
}

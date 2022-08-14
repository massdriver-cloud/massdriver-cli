package definition

type Definition map[string]interface{}

type DefinitionFile struct {
	Schema string `json:"$schema" yaml:"$schema"`
	Md     Md     `json:"$md" yaml:"$md"`
	Type   string `json:"type" yaml:"type"`
	Title  string `json:"title" yaml:"title"`
	// can be bool or object
	AdditionalProperties bool                   `json:"additionalProperties" yaml:"additionalProperties"`
	Required             []string               `json:"required" yaml:"required"`
	Properties           map[string]interface{} `json:"properties" yaml:"properties"`
}

type Md struct {
	Name         string        `json:"name" yaml:"name"`
	Access       string        `json:"access" yaml:"access"`
	Provisioners Priovisioners `json:"provisioners" yaml:"provisioners"`
}

type Priovisioners struct {
	Terraform string `json:"terraform" yaml:"terraform"`
}

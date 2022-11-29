package api2

type Artifact struct {
	Name string
	ID   string
}

type ArtifactDefinition struct {
	Name string
}

type Project struct {
	ID            string
	Slug          string
	DefaultParams map[string]interface{}
}

type Deployment struct {
	ID     string
	Status string
}

type Environment struct {
	ID   string
	Slug string
}

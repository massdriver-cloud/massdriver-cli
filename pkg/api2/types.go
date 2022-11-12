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

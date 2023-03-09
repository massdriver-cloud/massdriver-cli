package template

type Data struct {
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	Access          string                 `json:"access"`
	Location        string                 `json:"location"`
	TemplateName    string                 `json:"templateName"`
	TemplateSource  string                 `json:"templateSource"`
	OutputDir       string                 `json:"outputDir"`
	Type            string                 `json:"type"`
	Connections     map[string]interface{} `json:"connections"`
	CloudPrefix     string                 `json:"cloudPrefix"`
	RepoName        string                 `json:"repoName"`
	RepoNameEncoded string                 `json:"repoNameEncoded"`
}

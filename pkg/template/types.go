package template

type Data struct {
	Name           string
	Description    string
	Access         string
	Chart          string
	Location       string
	TemplateName   string
	TemplateSource string
	OutputDir      string
	Type           string
	CloudProvider  string
	Dependencies   map[string]string
}

package template

type Data struct {
	Name           string
	Description    string
	Access         string
	Location       string
	TemplateName   string
	TemplateSource string
	OutputDir      string
	Type           string
	Connections    map[string]interface{}
}

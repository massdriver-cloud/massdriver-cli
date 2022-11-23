package api

type JSONScalar string

func (p JSONScalar) GetGraphQLType() string {
	return "JSON"
}

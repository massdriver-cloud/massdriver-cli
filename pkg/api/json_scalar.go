package api

type JSONScalar string

func (p JSONScalar) GetGraphQLType() string {
	return "JSON"
}

// TODO: implement marshaling
// func (p JSONScalar) UnmarshalJSON(data []byte) error { panic("mock implementation") }
// func (p JSONScalar) MarshalJSON() ([]byte, error)    { panic("mock implementation") }

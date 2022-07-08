package jsonschema

import "encoding/json"

func GetJSONSchema(path string) (Schema, error) {
	schema := Schema{}
	sl := Loader(path)

	schemaSrc, err := sl.LoadJSON()
	if err != nil {
		return schema, err
	}

	byteData, err := json.Marshal(schemaSrc)
	if err != nil {
		return schema, err
	}

	if marshalErr := json.Unmarshal(byteData, &schema); marshalErr != nil {
		return schema, marshalErr
	}
	return schema, nil
}

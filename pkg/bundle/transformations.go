package bundle

var paramsTransformations = []func(map[string]interface{}) error{}
var connectionsTransformations = []func(map[string]interface{}) error{}
var artifactsTransformations = []func(map[string]interface{}) error{}
var uiTransformations = []func(map[string]interface{}) error{}

func ApplyTransformations(schema map[string]interface{}, transformations []func(map[string]interface{}) error) error {
	for _, transformation := range transformations {
		err := transformation(schema)
		if err != nil {
			return err
		}
	}

	for _, v := range schema {
		_, isObject := v.(map[string]interface{})
		if isObject {
			err := ApplyTransformations(v.(map[string]interface{}), transformations)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

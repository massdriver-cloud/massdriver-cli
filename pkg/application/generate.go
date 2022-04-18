package application

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"

	"gopkg.in/yaml.v2"
)

type PropertiesMap map[string]*Property

// Property is a single JSON Schema property field
type Property struct {
	AdditionalProperties *bool         `json:"additionalProperties,omitempty"`
	Title                string        `json:"title,omitempty"`
	Type                 string        `json:"type,omitempty"`
	Items                *Property     `json:"items,omitempty"`
	Properties           PropertiesMap `json:"properties,omitempty"`
	Required             []string      `json:"required,omitempty"`
}

// Schema is a flimsy representation of a JSON Schema.
// It provides just enough structure to get type information.
// type Schema struct {
// 	Properties PropertiesMap `json:"properties"`
// 	Required   []string      `json:"required,omitempty"`
// }

func Generate() {
	var config map[interface{}]interface{}

	yamlFile, err := ioutil.ReadFile("./src/application/testdata/values.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	schema, err := convertObject(config)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	schemaJson, err := json.Marshal(&schema)
	if err != nil {
		panic("Error encoding yaml")
	}

	err = ioutil.WriteFile("schema-values.json", schemaJson, 0644)
	if err != nil {
		panic("Unable to write data into the file")
	}
}

func convertObject(values map[interface{}]interface{}) (*Property, error) {
	schema := new(Property)

	schema.Properties = PropertiesMap{}

	for key, value := range values {
		keyString := key.(string)
		schema.Required = append(schema.Required, keyString)
		property, err := generateProperty(keyString, value)
		if err != nil {
			return schema, err
		}
		schema.Properties[keyString] = property
	}

	return schema, nil
}

func convertArray(array []interface{}) (*Property, error) {
	property := new(Property)
	var items *Property
	var err error

	if len(array) == 0 {
		// default to a string
		items, err = generateProperty("elem", "")

	} else {
		items, err = generateProperty("elem", array[0])
	}
	if err != nil {
		return nil, err
	}
	property.Items = items

	return property, nil
}

func generateProperty(name string, any interface{}) (*Property, error) {
	property := new(Property)

	val := getValue(any)

	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		property, err := convertArray(any.([]interface{}))
		if err != nil {
			return property, err
		}
		property.Title = name
		property.Type = "array"
		return property, nil
	case reflect.Map:
		property, err := convertObject(any.(map[interface{}]interface{}))
		if err != nil {
			return property, err
		}
		addlProp := false
		property.Title = name
		property.Type = "object"
		property.AdditionalProperties = &addlProp
		return property, nil
	case reflect.Int, reflect.Float64:
		property.Title = name
		property.Type = "number"
		return property, nil
	case reflect.String:
		property.Title = name
		property.Type = "string"
		return property, nil
	case reflect.Bool:
		property.Title = name
		property.Type = "boolean"
		return property, nil
	default:
		fmt.Printf("%s: %v\n", name, val.Kind())
		return property, nil
	}
}

func getValue(any interface{}) reflect.Value {
	val := reflect.ValueOf(any)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return val
}

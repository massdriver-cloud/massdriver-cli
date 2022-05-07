package application

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func Parse(path string) (*Application, error) {
	app := new(Application)

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, app)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return app, nil
}

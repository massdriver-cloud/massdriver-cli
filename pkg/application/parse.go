package application

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

// deployment:
//   type: chart
//   repo: # nil, "./path", "repo/chart"
//     name: retool
//     url: https://charts.retool.com

// params:
//   license_key:
//     title: License Key
//     type: string
//     required: true
//     chartValuePath: config.licenseKey
//   domain_name:
//     title: DNS Name
//     type: string
//     required: true
//     chartValuePath: ingress.hostName

// dependencies:
//   - type: massdriver/rdbms-connection
//     field: database
//     env:
//       - name: DATABASE_URL
//         value: ${data.authentication.connection_string}
//     port: 5432
//     protocol: tcp
//   - type: massdriver/kubernetes-cluster
//     field: kubernetes_cluster
//   - type: massdriver/aws-sqs-queue
//     field: queue
//     env:
//       - name: MY_QUEUE_ARN
//         value: ${data.infrastructure.arn}
//     policy: read
//   - type: massdriver/aws-sns-topic
//     field: topic
//     env:
//       - name: MY_TOPIC_ARN
//         value: ${data.infrastructure.arn}
//     policy: write

// type Application struct {
// 	Bundle       string
// 	Title        string
// 	Description  string
// 	Deployment   ApplicationDeployment
// 	Params       map[string]interface{}
// 	Dependencies []ApplicationDependencies
// }

// type ApplicationDeployment struct {
// 	Type string
// }

// type ApplicationParams struct {
// }

// type ApplicationDependencies struct {
// 	Type     string
// 	Field    string
// 	Required *bool `yaml:"required,omitempty"`
// 	Env      []ApplicationDependenciesEnvs
// 	Policy   string
// }

// type ApplicationDependenciesEnvs struct {
// 	Name  string ""
// 	Value string
// }

func Parse(path string) {
	var app Application

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &app)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	fmt.Printf("app: %v\n", app)

	var bundle Bundle

	bundle.Schema = "draft-07"
	bundle.UUID = "deadbeef-need-to-generate"
	bundle.Title = app.Title
	bundle.Description = app.Description
	bundle.Provisioner = "terraform"
	bundle.Access = "private"
	// Trim whitespace, all lowercase, spaces to hyphens. Maybe random suffix for globally unique? seems gross
	bundle.Type = strings.ReplaceAll(strings.ToLower(strings.Trim(app.Title, " ")), " ", "-")

	// Params
	bundle.Params = app.Params

	// Connections
	bundle.Connections.Properties = make(map[string]map[string]string)
	for idx := range app.Dependencies {
		conn := app.Dependencies[idx]
		if conn.Required == nil || *conn.Required {
			bundle.Connections.Required = append(bundle.Connections.Required, conn.Field)
		}
		bundle.Connections.Properties[conn.Field] = map[string]string{"ref": conn.Type}
	}

	bundleYaml, err := yaml.Marshal(&bundle)
	if err != nil {
		panic("Error encoding yaml")
	}

	err = ioutil.WriteFile("apptest.yaml", bundleYaml, 0644)
	if err != nil {
		panic("Unable to write data into the file")
	}
}

// # This file will be used to generate all of the schema-*.json files in a bundle
// schema: draft-07
// uuid: "0e4d2694-8f58-11ec-93e1-1e007b349036"
// title: "Azure Regional Cloud"
// description: "Massdriver Regional Cloud for Azure provides a secure, production-ready VPC with observability, alerting, and CI/CD."
// provisioner: "terraform"
// access: "private"
// type: "arm-regional-cloud"

// # schema-params.json
// # JSON Schema sans-fields above
// params:
//   examples:
//     - __name: default
//       azure_region: "Central US"
//       cidr: 10.0.0.0/16
//   required:
//     - azure_region
//     - cidr
//   properties:
//     azure_region:
//       $ref: ../../definitions/types/azure-region.json
//     cidr:
//       $ref: ../../definitions/types/cidr.json

// connections:
//   required:
//     - azure_service_principal
//   properties:
//     azure_service_principal:
//       $ref: ../../definitions/artifacts/arm-service-principal.json

// # schema-artifacts.json
// # Named list of output artifacts  (map[name]artifact)
// artifacts:
//   required:
//     - mrc
//   properties:
//     mrc:
//       $ref: ../../definitions/artifacts/azure-regional-cloud.json

// # schema-ui.json
// # List of form customizations for params-schema
// ui:

type Bundle struct {
	Schema      string
	UUID        string
	Title       string
	Description string
	Provisioner string
	Access      string
	Type        string
	Params      map[string]interface{}
	Connections BundleConnections
}

type BundleConnections struct {
	Required   []string
	Properties map[string]map[string]string
}

package application

type Application struct {
	Schema       string                    `json:"schema" yaml:"schema"`
	Name         string                    `json:"name" yaml:"name"`
	Description  string                    `json:"description" yaml:"description"`
	Ref          string                    `json:"ref" yaml:"ref"`
	Access       string                    `json:"access" yaml:"access"`
	Deployment   ApplicationDeployment     `json:"deployment" yaml:"deployment"`
	Params       map[string]interface{}    `json:"params" yaml:"params"`
	Dependencies []ApplicationDependencies `json:"dependencies" yaml:"dependencies"`
}

type ApplicationDeployment struct {
	Type       string `json:"type" yaml:"type"`
	Path       string `json:"path,omitempty" yaml:"path,omitempty"`
	Chart      string `json:"chart,omitempty" yaml:"chart,omitempty"`
	Repository string `json:"repository,omitempty" yaml:"repository,omitempty"`
}

type ApplicationDependencies struct {
	Type     string                        `json:"type" yaml:"type"`
	Field    string                        `json:"field" yaml:"field"`
	Required bool                          `json:"required,omitempty" yaml:"required,omitempty"`
	Envs     []ApplicationDependenciesEnvs `json:"envs" yaml:"envs"`
	Policies   []string                    `json:"policies,omitempty" yaml:"policies,omitempty"`
}

type ApplicationDependenciesEnvs struct {
	Name string `json:"name" yaml:"name"`
	Path string `json:"path" yaml:"path"`
}

var SimpleUi = `{
	"ui:order": [
		"name",
		"namespace",
		"image",
		"resource_requests",
		"autoscaling",
		"envs",
		"port",
		"ingress"
	],
	"autoscaling": {
		"items": {
			"ui:order": [
				"enabled",
				"minReplicas",
				"maxReplicas",
				"targetCPUUtilizationPercentage"
			]
		}
	},
	"ingress": {
		"items": {
			"ui:order": [
				"enabled",
				"host",
				"path"
			]
		}
	}
}`
var SimpleParams = `{
	"required": [
		"name",
		"namespace",
		"image",
		"resource_requests"
	],
	"properties": {
		"name": {
			"title": "Name",
			"type": "string",
			"description": "Name of the application"
		},
		"namespace": {
			"title": "Namespace",
			"type": "string",
			"description": "Kubernetes namespace to run application within"
		},
		"image": {
			"type": "object",
			"title": "Image",
			"description": "Container image to use",
			"required": [
				"repository",
				"tag"
			],
			"properties": {
				"repository": {
					"type": "string",
					"title": "Image Repository",
					"description": "Docker image run"
				},
				"tag": {
					"type": "string",
					"title": "Image Tag",
					"description": "If you are using continuous delivery, this should be a mutable tag that points to the most recent version (such as \"latest\")."
				}
			}
		},
		"resource_requests": {
			"type": "object",
			"title": "Resources",
			"required": [
				"cpu",
				"memory"
			],
			"properties": {
				"cpu": {
					"type": "string",
					"title": "CPU"
				},
				"memory": {
					"type": "string",
					"title": "Memory"
				}
			}
		},
		"port": {
			"type": "integer",
			"title": "Port",
			"description": "If this is set, a kubernetes service will be created exposing this port. If ingress is enabled, it will be connected to this port.",
			"minimum": 1,
			"maximum": 65535
		},
		"autoscaling": {
			"type": "object",
			"title": "Autoscaling",
			"required": [
				"minReplicas",
				"maxReplicas",
				"targetCPUUtilizationPercentage"
			],
			"properties": {
				"enabled": {
					"type": "boolean",
					"title": "Enabled",
					"description": "Enable pod autoscaling"
				},
				"minReplicas": {
					"type": "integer",
					"title": "Minimum Replicas",
					"minimum": 1
				},
				"maxReplicas": {
					"type": "integer",
					"title": "Maximum Replicas",
					"minimum": 1
				},
				"targetCPUUtilizationPercentage": {
					"type": "integer",
					"title": "Target CPU Utilization Percentage",
					"minimum": 1,
					"maximum": 100
				}
			}
		},
		"envs": {
			"title": "Environment Variables",
			"type": "array",
			"description": "Additional environment variables to set",
			"default": [],
			"items": {
				"type": "object",
				"required": [
					"name",
					"path"
				],
				"properties": {
					"name": {
						"type": "string",
						"title": "Name"
					},
					"path": {
						"type": "string",
						"title": "JSON Path"
					}
				}
			}
		},
		"ingress": {
			"title": "Ingress",
			"type": "object",
			"properties": {
				"enabled": {
					"title": "Enable Ingress",
					"type": "boolean",
					"description": "Enabling this will create ingress configurations to allow internet traffic to reach our application on the specified host and path",
					"default": false
				},
				"host": {
					"type": "string",
					"title": "Host",
					"pattern": "^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$",
					"description": "Endpoint the application should be accessible on. Omit the protocol (https://). Include the path, if necessary (app.mydomain.com/path)"
				},
				"path": {
					"type": "string",
					"title": "Path",
					"pattern": "^(\\/[-a-zA-Z0-9()@:%_\\+.~#?&\\/=]*)$",
					"description": "Endpoint the application should be accessible on. Omit the protocol (https://). Include the path, if necessary (app.mydomain.com/path)"
				}
			}
		}
	}
}`

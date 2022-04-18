package massdriver

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/kelseyhightower/envconfig"
)

type MassdriverClient struct {
	Specification Specification
}

type Specification struct {
	DeploymentID  string `envconfig:"DEPLOYMENT_ID"`
	EventTopicARN string `envconfig:"EVENT_TOPIC_ARN"`
	Provisioner   string `envconfig:"PROVISIONER"`
}

func InitializeMassdriverClient() (*MassdriverClient, error) {
	client := new(MassdriverClient)
	err := envconfig.Process("massdriver", &client.Specification)
	if err != nil {
		return nil, err
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	_ = cfg

	return client, nil
}

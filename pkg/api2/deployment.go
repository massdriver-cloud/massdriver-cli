package api2

import (
	"context"

	"github.com/Khan/genqlient/graphql"
)

func GetDeployment(client graphql.Client, orgID string, id string) (Deployment, error) {
	response, err := getDeploymentById(context.Background(), client, orgID, id)

	return response.Deployment.toDeployment(), err
}

func (d *getDeploymentByIdDeployment) toDeployment() Deployment {
	return Deployment{
		ID:     d.Id,
		Status: d.Status,
	}
}

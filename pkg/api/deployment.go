package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
	"github.com/rs/zerolog/log"
)

type Deployment struct {
	ID     string
	Status string
}

func GetDeployment(client *graphql.Client, orgID string, id string) (*Deployment, error) {
	log.Debug().Str("deploymentId", id).Msg("Getting deployment")

	var q struct {
		Deployment struct {
			ID     graphql.String
			Status graphql.String
		} `graphql:"deployment(id: $id, organizationId: $organizationId)"`
	}

	variables := map[string]interface{}{
		"id":             graphql.ID(id),
		"organizationId": graphql.ID(orgID),
	}

	err := client.Query(context.Background(), &q, variables)

	if err != nil {
		return nil, err
	}

	deployment := Deployment{
		ID:     string(q.Deployment.ID),
		Status: string(q.Deployment.Status),
	}

	log.Debug().Str("deploymentId", string(q.Deployment.ID)).Str("status", string(q.Deployment.Status)).Msg("Got deployment")

	return &deployment, nil
}

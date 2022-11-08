package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
	"github.com/rs/zerolog/log"
)

type Project struct {
	ID            string
	DefaultParams interface{}
	Slug          string
}

func GetProject(client *graphql.Client, orgID string, id string) (*Project, error) {
	log.Debug().Str("projectID", id).Msg("Getting project")

	var q struct {
		Project struct {
			ID            graphql.String
			DefaultParams interface{} `scalar:"true"`
			Slug          graphql.String
		} `graphql:"project(id: $id, organizationId: $organizationId)"`
	}

	variables := map[string]interface{}{
		"id":             graphql.ID(id),
		"organizationId": graphql.ID(orgID),
	}

	err := client.Query(context.Background(), &q, variables)

	if err != nil {
		return nil, err
	}

	project := Project{
		ID:            string(q.Project.ID),
		Slug:          string(q.Project.Slug),
		DefaultParams: q.Project.DefaultParams,
	}

	log.Debug().Str("projectId", string(q.Project.ID)).Msg("Got project")
	return &project, nil
}

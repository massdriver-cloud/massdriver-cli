package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
	"github.com/massdriver-cloud/massdriver-cli/pkg/jsonschema"
	"github.com/rs/zerolog/log"
)

type RepositoryAuth struct {
	ID         string
	Name       string
	NamePrefix string
	ProjectID  string
	ManifestID string
	TargetID   string
	// TODO: implement caching + expiration
	Token            string
	ActiveDeployment Deployment
	ParamsSchema     jsonschema.Schema
}

func GetToken(client *graphql.Client, orgID string, name string) (*RepositoryAuth, error) {
	log.Debug().Str("packageName", name).Msg("Getting token")

	// TODO: update the query after syncing w/ Andreas
	var q struct {
		GetPackageByNamingConvention struct {
			ID           graphql.String
			ParamsSchema graphql.String `scalar:"true"`
		} `graphql:"getPackageByNamingConvention(name: $name, organizationId: $organizationId)"`
	}

	variables := map[string]interface{}{
		"name":           graphql.String(name),
		"organizationId": graphql.ID(orgID),
	}

	err := client.Query(context.Background(), &q, variables)

	if err != nil {
		return nil, err
	}

	auth := RepositoryAuth{
		ID:   string(q.GetPackageByNamingConvention.ID),
		Name: name,
	}

	return &auth, nil
}

// TODO: consider: https://github.com/Khan/genqlient (need to look into testing w/ it, but looks nice for a lot of queries)
// TODO: websocket or longpoll gql subscription - there isnt a phoenix socket impl for golang I could find, so we'll probably have to longpoll
package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/rs/zerolog/log"
)

type Package struct {
	ID               string
	Name             string
	NamePrefix       string
	ManifestID       string
	TargetID         string
	ActiveDeployment Deployment
}

const deploymentTimeout time.Duration = time.Duration(5) * time.Minute
const deploymentStatusSleep time.Duration = time.Duration(10) * time.Second

func DeployPackage(client *graphql.Client, orgID string, name string) (*Deployment, error) {
	log.Debug().Str("packageName", name).Msg("Deploying package")
	pkg, err := GetPackage(client, orgID, name)
	if err != nil {
		return nil, err
	}
	var m struct {
		DeployPackage struct {
			Successful graphql.Boolean
			Result     struct {
				ID graphql.ID
			}
		} `graphql:"deployPackage(organizationId: $organizationId, manifestID: $manifestID, targetID: $targetID)"`
	}

	variables := map[string]interface{}{
		"manifestID":     graphql.ID(pkg.ManifestID),
		"targetID":       graphql.ID(pkg.TargetID),
		"organizationId": graphql.ID(orgID),
	}

	err = client.Mutate(context.Background(), &m, variables)

	if err != nil {
		return nil, err
	}

	did := fmt.Sprintf("%s", m.DeployPackage.Result.ID)
	log.Info().Str("packageName", name).Str("deploymentId", did).Msg("Deployment enqueued")
	deployment, err := checkDeploymentStatus(client, orgID, did, deploymentTimeout)

	if err != nil {
		return nil, err
	}

	return deployment, nil
}

func GetPackage(client *graphql.Client, orgID string, name string) (*Package, error) {
	log.Debug().Str("packageName", name).Msg("Getting package")

	var q struct {
		GetPackageByNamingConvention struct {
			ID         graphql.String
			NamePrefix graphql.String
			Manifest   struct {
				ID graphql.String
			}
			ActiveDeployment struct {
				ID     graphql.String
				Status graphql.String
			}
			Target struct {
				ID graphql.String
			}
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

	pkg := Package{
		ID:         string(q.GetPackageByNamingConvention.ID),
		Name:       name,
		NamePrefix: string(q.GetPackageByNamingConvention.NamePrefix),
		ManifestID: string(q.GetPackageByNamingConvention.Manifest.ID),
		TargetID:   string(q.GetPackageByNamingConvention.Target.ID),
		// NOTE: this is any _previous_ ActiveDeployment that is running
		ActiveDeployment: Deployment{
			ID:     string(q.GetPackageByNamingConvention.ActiveDeployment.ID),
			Status: string(q.GetPackageByNamingConvention.ActiveDeployment.Status),
		},
	}

	log.Debug().
		Str("packageName", name).
		Msg("Got package")

	return &pkg, nil
}

func checkDeploymentStatus(client *graphql.Client, orgID string, id string, timeout time.Duration) (*Deployment, error) {
	deployment, err := GetDeployment(client, orgID, id)

	if err != nil {
		return nil, err
	}

	timeout -= deploymentStatusSleep

	switch deployment.Status {
	case "COMPLETED":
		log.Debug().Str("deploymentId", id).Msg("Deployment completed")
		return deployment, nil
	case "FAILED":
		log.Debug().Str("deploymentId", id).Msg("Deployment failed")
		return nil, errors.New("Deployment failed")
	default:
		time.Sleep(deploymentStatusSleep)
		return checkDeploymentStatus(client, orgID, id, timeout)
	}
}

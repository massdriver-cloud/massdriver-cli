// TODO: consider: https://github.com/Khan/genqlient (need to look into testing w/ it, but looks nice for a lot of queries)
// TODO: websocket or longpoll gql subscription - there isnt a phoenix socket impl for golang I could find, so we'll probably have to longpoll
package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hasura/go-graphql-client"
	"github.com/rs/zerolog/log"
)

var ErrDeploymentFailed = errors.New("deployment failed")
var ErrDeploymentComplete = errors.New("deployment succeeded")

type Package struct {
	ID               string
	Name             string
	NamePrefix       string
	ManifestID       string
	TargetID         string
	ActiveDeployment Deployment
}

// const deploymentTimeout time.Duration = time.Duration(5) * time.Minute

// const deploymentStatusSleep time.Duration = time.Duration(10) * time.Second

func DeployPackage(client *graphql.Client, subClient *graphql.SubscriptionClient, orgID string, name string) (*Deployment, error) {
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
	var s struct {
		DeploymentLifecycleEvent struct {
		} `graphql:"subscription deploymentProgress($packageId: ID!, $organizationId: ID!)"`
	}
	subVariables := map[string]interface{}{
		"organizationId": graphql.ID(orgID),
		"packageId":      pkg.ID,
	}

	// TODO use deploymentTimeout
	// subClient = subClient.WithTimeout(deploymentTimeout)
	subID, err := subClient.Subscribe(s, subVariables, rawMessageHandler)
	if err != nil {
		return nil, err
	}
	defer subClient.Unsubscribe(subID) // nolint:errcheck
	defer subClient.Close()
	log.Debug().Str("subscription", subID).Str("deploymentId", did).Msg("subscribed to deployment progress")
	// deployment, err := checkDeploymentStatus(client, orgID, did, deploymentTimeout)
	deployment, err := GetDeployment(client, orgID, did)
	if err != nil {
		return nil, err
	}

	subErr := subClient.Run()
	if errors.Is(subErr, ErrDeploymentComplete) {
		log.Info().Str("deploymentId", did).Msg("Deployment succeeded")
		return deployment, nil
	}
	if errors.Is(subErr, ErrDeploymentFailed) {
		log.Error().Str("deploymentId", did).Msg("Deployment failed")
		return deployment, ErrDeploymentFailed
	}

	return deployment, nil
}

func rawMessageHandler(message []byte, err error) error {
	if err != nil {
		return fmt.Errorf("error from server in subscription: %w", err)
	}
	rawMessage := make(map[string]interface{})
	err = json.Unmarshal(message, &rawMessage)
	// TODO handle various types of messages in the pheonix protocol: join / leave room / message / server error
	// TODO log something prettier / more structured.  This is just for debugging
	log.Info().Msg(string(message))
	if err != nil {
		return fmt.Errorf("error unmarshalling json message: %w", err)
	}
	// TODO should throw ErrDeploymentFailed if the deployment failed
	// TODO should throw ErrDeploymentComplete if the deployment succeeded
	// TODO should be ablet to indicate to subscriber to stop listening when we get throwing one of above.
	return nil
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
		"organizationId": graphql.String(orgID),
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

// func checkDeploymentStatus(client *graphql.Client, orgID string, id string, timeout time.Duration) (*Deployment, error) {
// 	deployment, err := GetDeployment(client, orgID, id)

// 	if err != nil {
// 		return nil, err
// 	}

// 	timeout -= deploymentStatusSleep

// 	switch deployment.Status {
// 	case "COMPLETED":
// 		log.Debug().Str("deploymentId", id).Msg("Deployment completed")
// 		return deployment, nil
// 	case "FAILED":
// 		log.Debug().Str("deploymentId", id).Msg("Deployment failed")
// 		return nil, errors.New("Deployment failed")
// 	default:
// 		time.Sleep(deploymentStatusSleep)
// 		return checkDeploymentStatus(client, orgID, id, timeout)
// 	}
// }

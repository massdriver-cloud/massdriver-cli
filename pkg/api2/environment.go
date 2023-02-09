package api2

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Khan/genqlient/graphql"
	"github.com/rs/zerolog/log"
)

func DeployPreviewEnvironment(client graphql.Client, orgID string, projectID string, credentials []Credential, packageParams map[string]interface{}, ciContext map[string]interface{}) (Environment, error) {
	ctx := context.Background()
	env := Environment{}

	input := PreviewEnvironmentInput{
		Credentials:   credentials,
		PackageParams: packageParams,
		CiContext:     ciContext,
	}

	response, err := deployPreviewEnvironment(ctx, client, orgID, projectID, input)

	if err != nil {
		return env, err
	}

	if response.DeployPreviewEnvironment.Successful {
		// TODO: is there a less obnoxious way to do this...
		env = Environment{
			ID:   response.DeployPreviewEnvironment.Result.Id,
			Slug: response.DeployPreviewEnvironment.Result.Slug,
			Project: Project{
				ID: response.DeployPreviewEnvironment.Result.Project.Id,
			},
		}
		return env, nil
	}

	log.Error().Str("project", projectID).Msg("Preview environment deployment failed.")
	msgs, err := json.Marshal(response.DeployPreviewEnvironment.Messages)
	if err != nil {
		return env, fmt.Errorf("failed to deploy preview environment and couldn't marshal error messages: %w", err)
	}

	return env, fmt.Errorf("failed to deploy environment: %v", string(msgs))
}

func DecommissionPreviewEnvironment(client graphql.Client, orgID string, targetID string) (Environment, error) {
	ctx := context.Background()
	env := Environment{}

	response, err := decommissionPreviewEnvironment(ctx, client, orgID, targetID)

	if err != nil {
		return env, err
	}

	if response.DecommissionPreviewEnvironment.Successful {
		env = Environment{
			ID:   response.DecommissionPreviewEnvironment.Result.Id,
			Slug: response.DecommissionPreviewEnvironment.Result.Slug,
		}
		return env, nil
	}

	log.Error().Str("target", targetID).Msg("Preview environment decommission failed.")
	msgs, err := json.Marshal(response.DecommissionPreviewEnvironment.Messages)
	if err != nil {
		return env, fmt.Errorf("failed to decommission preview environment and couldn't marshal error messages: %w", err)
	}

	return env, fmt.Errorf("failed to decommission environment: %v", string(msgs))
}

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
		env.ID = response.DeployPreviewEnvironment.Result.Id
		env.Slug = response.DeployPreviewEnvironment.Result.Slug
		return env, nil
	}

	log.Error().Str("project", projectID).Msg("Preview environment deployment failed.")
	msgs, err := json.Marshal(response.DeployPreviewEnvironment.Messages)
	if err != nil {
		return env, fmt.Errorf("failed to deploy preview environment and couldn't marshal error messages: %w", err)
	}

	return env, fmt.Errorf("failed to deploy environment: %v", string(msgs))
}
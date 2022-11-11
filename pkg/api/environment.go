// Environment (target) management
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/rs/zerolog/log"
)

const urlTemplate = "https://app.massdriver.cloud/projects/%s/targets/%v"

type Environment struct {
	ConfigTemplate string
	Config         string
}

func DeployPreviewEnvironment(client *graphql.Client, orgID string, id string, templateData io.Reader) (*Environment, error) {
	log.Info().Str("project", id).Msg("Deploying preview environment.")

	buf := new(strings.Builder)
	_, err := io.Copy(buf, templateData)

	if err != nil {
		return nil, err
	}

	envVars := getOsEnv()
	template := buf.String()
	config := os.Expand(template, func(s string) string { return envVars[s] })

	environment := Environment{
		ConfigTemplate: template,
		Config:         config,
	}

	var m struct {
		EnvironmentPayload struct {
			Successful graphql.Boolean
			Result     struct {
				ID graphql.ID
			}
			Messages []struct {
				Code    graphql.String
				Field   graphql.String
				Message graphql.String
				Options []struct {
					Key   graphql.String
					Value graphql.String
				}
			}
		} `graphql:"createPreviewEnvironment(description: $description, name: $name, slug: $slug, organizationId: $organizationId, projectId: $projectId, envParams: $envParams, prMetadata: $prMetadata)"`
	}

	// TODO: move slug & name gen into the domain function.
	// TODO: pull description from CI Context
	envSlug := generateEnvSlug()
	envName := fmt.Sprintf("✨ Preview environment: %s", envSlug)

	variables := map[string]interface{}{
		"description":    graphql.String("Test description."),
		"name":           graphql.String(envName),
		"slug":           graphql.String(envSlug),
		"organizationId": graphql.ID(orgID),
		"projectId":      graphql.ID(id),
		"envParams":      JSONScalar(environment.Config),
		"prMetadata":     JSONScalar("{}"),
	}

	err = client.Mutate(context.Background(), &m, variables, graphql.OperationName("createPreviewEnvironment"))

	if err != nil {
		return nil, err
	}

	if m.EnvironmentPayload.Successful {
		url := fmt.Sprintf(urlTemplate, id, m.EnvironmentPayload.Result.ID)
		log.Info().
			Str("project", id).
			Str("url", url).
			Interface("environment", m.EnvironmentPayload.Result.ID).
			Msg("Preview environment deploying.")
		exec.Command("open", url).Run()

		return &environment, nil
	}

	log.Error().Str("project", id).Msg("Preview environment deployment failed.")
	msgs, err := json.Marshal(m.EnvironmentPayload.Messages)
	if err != nil {
		return &environment, fmt.Errorf("failed to deploy preview environment and couldn't marshal error messages: %w", err)
	}

	return &environment, fmt.Errorf("failed to deploy environment: %v", string(msgs))
}

func getOsEnv() map[string]string {
	getenvironment := func(data []string, getkeyval func(item string) (key, val string)) map[string]string {
		items := make(map[string]string)
		for _, item := range data {
			key, val := getkeyval(item)
			items[key] = val
		}
		return items
	}

	osEnv := getenvironment(os.Environ(), func(item string) (key, val string) {
		splits := strings.Split(item, "=")
		key = splits[0]
		val = splits[1]
		return
	})

	return osEnv
}

func generateEnvSlug() string {
	rand.Seed(time.Now().Unix())
	charset := "bcdfghjklmnpqrstvwxz0123456789"
	length := 7

	ran_str := make([]byte, length)
	ran_str[0] = charset[0]

	for i := 1; i < length; i++ {
		ran_str[i] = charset[rand.Intn(len(charset))]
	}

	str := string(ran_str)
	return str
}

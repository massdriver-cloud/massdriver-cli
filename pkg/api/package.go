// TODO: consider: https://github.com/Khan/genqlient (need to look into testing w/ it, but looks nice for a lot of queries)
// TODO: websocket or longpoll gql subscription - there isnt a phoenix socket impl for golang I could find, so we'll probably have to longpoll
package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/hasura/go-graphql-client"
	"github.com/jackdelahunt/survey-json-schema/pkg/surveyjson"
	"github.com/massdriver-cloud/massdriver-cli/pkg/jsonschema"
	"github.com/rs/zerolog/log"
)

type Package struct {
	ID               string
	Name             string
	NamePrefix       string
	ManifestID       string
	TargetID         string
	ActiveDeployment Deployment
	ParamsSchema     jsonschema.Schema
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
		ParamsSchema: *deserializeSchema(q.GetPackageByNamingConvention.ParamsSchema),
	}

	log.Debug().
		Str("packageName", name).
		Msg("Got package")

	return &pkg, nil
}

func deserializeSchema(schema graphql.String) *jsonschema.Schema {
	var s jsonschema.Schema
	// byteData, err := json.Marshal(schema)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("Failed to marshal schema")
	// }
	byteData := []byte(schema)

	// do a little json dance to get the schema into our structured go type
	if marshalErr := json.Unmarshal(byteData, &s); marshalErr != nil {
		log.Fatal().Err(marshalErr).Msg("Failed to unmarshal schema from API")
	}
	return &s
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

type Params string

func (p Params) GetGraphQLType() string {
	return "JSON"
}

func promptForConfigurableVariables(pkg *Package) (Params, error) {
	ret := Params("{}")
	// start by prompting for which set of presets to use
	exampleNames := make([]string, len(pkg.ParamsSchema.Examples))
	exampleMap := make(map[string]jsonschema.Example)
	for i, example := range pkg.ParamsSchema.Examples {
		exampleMap[example.Name] = example
		exampleNames[i] = example.Name
	}
	initialValues := make(map[string]interface{})
	if len(exampleNames) > 0 {
		var qs = []*survey.Question{
			{
				Name: "Presets",
				Prompt: &survey.Select{
					Message: "Choose a guided configuration for this package:",
					Options: exampleNames,
					Description: func(value string, index int) string {
						bytes, err := json.MarshalIndent(exampleMap[value].Values, "", "  ")
						if err != nil {
							log.Debug().Err(err).Msg("Failed to get example description")
							return ""
						}
						return string(bytes)
					},
				},
				Validate: survey.Required,
			},
		}

		var answers struct {
			Presets string
		}

		err := survey.Ask(qs, &answers)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to prompt for presets")
			return ret, err
		}
		initialValues = exampleMap[answers.Presets].Values
	}

	options := surveyjson.JSONSchemaOptions{
		Out:                 os.Stdout,
		In:                  os.Stdin,
		OutErr:              os.Stderr,
		AskExisting:         true,
		AutoAcceptDefaults:  false,
		NoAsk:               false,
		IgnoreMissingValues: false,
	}

	schemaBytes, err := json.Marshal(pkg.ParamsSchema)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to marshal params schema")
		return ret, err
	}

	result, err := options.GenerateValues(schemaBytes, initialValues)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to collect param values")
		return ret, err
	}
	unmarshalErr := json.Unmarshal(result, &ret)
	return ret, unmarshalErr
}

func ConfigurePackage(client *graphql.Client, orgID string, name string) (*Package, error) {
	log.Debug().Str("packageName", name).Msg("Configuring package")
	pkg, err := GetPackage(client, orgID, name)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get package")
		return nil, err
	}
	params, err := promptForConfigurableVariables(pkg)
	if err != nil {
		return nil, err
	}
	var m struct {
		PackagePayload struct {
			Successful graphql.Boolean `json:"successful"`
			Result     struct {
				ID graphql.ID `json:"id"`
			} `json:"result"`
			Messages []struct {
				Code    graphql.String `json:"code"`
				Field   graphql.String `json:"field"`
				Message graphql.String `json:"message"`
				Options []struct {
					Key   graphql.String `json:"key"`
					Value graphql.String `json:"value"`
				} `json:"options"`
			}
		} `graphql:"configurePackage(manifestID: $manifestID, organizationId: $organizationId, params: $params, targetID: $targetID)"`
	}

	// TODO not sure if graphql library expects a JSON type field as a string or a map
	paramString, err := json.Marshal(params)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to marshal params")
		return nil, err
	}
	log.Debug().Str("params", fmt.Sprintf("%v", params)).Msg("Params")
	log.Debug().Str("paramString", string(paramString)).Msg("Params Marshaled")
	variables := map[string]interface{}{
		"manifestID":     graphql.ID(pkg.ManifestID),
		"targetID":       graphql.ID(pkg.TargetID),
		"organizationId": graphql.ID(orgID),
		"params":         Params(paramString),
	}

	err = client.Mutate(context.Background(), &m, variables)

	if err != nil {
		return nil, err
	}

	if m.PackagePayload.Successful {
		log.Info().Str("packageName", name).Interface("packageID", m).Msg("Package configured successfully")
		return pkg, nil
	}
	log.Error().Str("packageName", name).Interface("packageID", m).Msg("Package configure failed")
	msgs, err := json.Marshal(m.PackagePayload.Messages)
	if err != nil {
		return pkg, errors.New(fmt.Sprintf("failed to configure package and couldn't marshal error messages: %v", err)) //nolint: revive,gosimple
	}
	return pkg, errors.New(fmt.Sprintf("failed to configure package: %v", string(msgs))) //nolint: revive,gosimple
}

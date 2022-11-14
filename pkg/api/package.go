package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	gql "github.com/Khan/genqlient/graphql"
	"github.com/hasura/go-graphql-client"
	"github.com/massdriver-cloud/massdriver-cli/pkg/api2"
	"github.com/massdriver-cloud/massdriver-cli/pkg/jsonschema"
	"github.com/rs/zerolog/log"
	yaml "gopkg.in/yaml.v3"
)

type Package struct {
	ID               string
	Name             string
	NamePrefix       string
	ProjectID        string
	ManifestID       string
	TargetID         string
	ActiveDeployment api2.Deployment
	ParamsSchema     jsonschema.Schema
}

func (p *Package) GetMDMetadata() map[string]interface{} {
	if p != nil {
		return map[string]interface{}{
			"name_prefix": p.NamePrefix,
			"default_tags": map[string]string{
				"md-project":  "local",
				"md-target":   p.TargetID,
				"md-manifest": p.ManifestID,
				"md-package":  p.ID,
			},
			"observability": map[string]interface{}{
				"alarm_webhook_url": "https://placeholder.com",
			},
		}
	}
	return map[string]interface{}{}
}

const deploymentTimeout time.Duration = time.Duration(5) * time.Minute
const deploymentStatusSleep time.Duration = time.Duration(10) * time.Second

func DeployPackage(client *graphql.Client, client2 *gql.Client, orgID string, name string) (*api2.Deployment, error) {
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
	deployment, err := checkDeploymentStatus(client2, orgID, did, deploymentTimeout)

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
				ID      graphql.String
				Project struct {
					ID graphql.String
				}
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
		ProjectID:  string(q.GetPackageByNamingConvention.Target.Project.ID),
		ManifestID: string(q.GetPackageByNamingConvention.Manifest.ID),
		TargetID:   string(q.GetPackageByNamingConvention.Target.ID),
		// NOTE: this is any _previous_ ActiveDeployment that is running
		ActiveDeployment: api2.Deployment{
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
	byteData := []byte(schema)

	// do a little json dance to get the schema into our structured go type
	if marshalErr := json.Unmarshal(byteData, &s); marshalErr != nil {
		log.Fatal().Err(marshalErr).Msg("Failed to unmarshal schema from API")
	}
	return &s
}

func checkDeploymentStatus(client *gql.Client, orgID string, id string, timeout time.Duration) (*api2.Deployment, error) {
	deployment, err := api2.GetDeployment(*client, orgID, id)

	if err != nil {
		return nil, err
	}

	timeout -= deploymentStatusSleep

	switch deployment.Status {
	case "COMPLETED":
		log.Debug().Str("deploymentId", id).Msg("Deployment completed")
		return &deployment, nil
	case "FAILED":
		log.Debug().Str("deploymentId", id).Msg("Deployment failed")
		return nil, errors.New("Deployment failed")
	default:
		time.Sleep(deploymentStatusSleep)
		return checkDeploymentStatus(client, orgID, id, timeout)
	}
}

func ReadParamsFromFile(path string, pkg *Package) (map[string]interface{}, error) {
	ret := map[string]interface{}{}
	file, err := os.Open(path)
	if err != nil {
		return ret, err
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	switch filepath.Ext(path) {
	case ".json":
		err = json.Unmarshal(byteValue, &ret)
		if _, ok := ret["md_metadata"]; !ok {
			ret["md_metadata"] = pkg.GetMDMetadata()
		}
	case ".yaml", ".yml":
		err = yaml.Unmarshal(byteValue, &ret)
		if _, ok := ret["md_metadata"]; !ok {
			ret["md_metadata"] = pkg.GetMDMetadata()
		}
	default:
		return ret, errors.New("invalid file type, must me json or yaml")
	}

	return ret, err
}

func ConfigurePackage(client *graphql.Client, orgID, name, paramValuePath string) (*Package, error) {
	log.Debug().Str("packageName", name).Msg("Configuring package")
	pkg, err := GetPackage(client, orgID, name)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get package")
		return nil, err
	}
	var params map[string]interface{}
	if len(paramValuePath) > 0 {
		params, err = ReadParamsFromFile(paramValuePath, pkg)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to read params from file")
			return nil, err
		}
	}
	var m struct {
		PackagePayload struct {
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
		} `graphql:"configurePackage(manifestID: $manifestID, organizationId: $organizationId, params: $params, targetID: $targetID)"`
	}

	paramString, err := json.Marshal(params)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to marshal params")
		return nil, err
	}
	variables := map[string]interface{}{
		"manifestID":     graphql.ID(pkg.ManifestID),
		"targetID":       graphql.ID(pkg.TargetID),
		"organizationId": graphql.ID(orgID),
		"params":         JSONScalar(paramString),
	}

	err = client.Mutate(context.Background(), &m, variables, graphql.OperationName("configurePackage"))

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
		return pkg, fmt.Errorf("failed to configure package and couldn't marshal error messages: %w", err)
	}
	return pkg, fmt.Errorf("failed to configure package: %v", string(msgs))
}

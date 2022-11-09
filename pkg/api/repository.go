package api

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

	auth := RepositoryAuth{
		ID:         string(q.GetPackageByNamingConvention.ID),
		Name:       name,
		NamePrefix: string(q.GetPackageByNamingConvention.NamePrefix),
		ProjectID:  string(q.GetPackageByNamingConvention.Target.Project.ID),
		ManifestID: string(q.GetPackageByNamingConvention.Manifest.ID),
		TargetID:   string(q.GetPackageByNamingConvention.Target.ID),
	}

	log.Debug().
		Str("packageName", name).
		Msg("Got package")

	return &auth, nil
}

package api2

import (
	"context"

	"github.com/Khan/genqlient/graphql"
)

func GetProject(client graphql.Client, orgID string, idOrSlug string) (Project, error) {
	response, err := getProjectById(context.Background(), client, orgID, idOrSlug)

	return response.Project.toProject(), err
}

func (p *getProjectByIdProject) toProject() Project {
	return Project{
		ID:            p.Id,
		Slug:          p.Slug,
		DefaultParams: p.DefaultParams,
	}
}

// Environment (target) management
package api

import (
	"io"
	"os"
	"strings"

	"github.com/hasura/go-graphql-client"
)

type Environment struct {
	ConfigTemplate string
	Config         string
}

func CreatePreviewEnvironment(client *graphql.Client, orgID string, id string, templateData io.Reader) (*Environment, error) {
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

	return &environment, nil
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

// import (
// 	"context"

// 	"github.com/hasura/go-graphql-client"
// 	"github.com/rs/zerolog/log"
// )

// type Project struct {
// 	ID            string
// 	DefaultParams interface{}
// 	Slug          string
// }

// func GetProject(client *graphql.Client, orgID string, id string) (*Project, error) {
// 	log.Debug().Str("projectID", id).Msg("Getting project")

// 	var q struct {
// 		Project struct {
// 			ID            graphql.String
// 			DefaultParams interface{} `scalar:"true"`
// 			Slug          graphql.String
// 		} `graphql:"project(id: $id, organizationId: $organizationId)"`
// 	}

// 	variables := map[string]interface{}{
// 		"id":             graphql.ID(id),
// 		"organizationId": graphql.ID(orgID),
// 	}

// 	err := client.Query(context.Background(), &q, variables)

// 	if err != nil {
// 		return nil, err
// 	}

// 	project := Project{
// 		ID:            string(q.Project.ID),
// 		Slug:          string(q.Project.Slug),
// 		DefaultParams: q.Project.DefaultParams,
// 	}

// 	log.Debug().Str("id", string(q.Project.ID)).Msg("Got project")
// 	return &project, nil
// }

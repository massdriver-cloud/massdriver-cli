package application

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
)

var MASSDRIVER_HELM_CHART_REPOSITORY = "https://massdriver-cloud.github.io/helm-charts"

func Generate(data *ApplicationTemplateData) error {

	cpo := action.ChartPathOptions{
		InsecureSkipTLSverify: true,
		RepoURL:               MASSDRIVER_HELM_CHART_REPOSITORY,
		Version:               ">0.0.0-0",
	}

	chartLocation := path.Join(data.Location, data.Name)

	if _, err := os.Stat(chartLocation); !os.IsNotExist(err) {
		return errors.New("specified directory already exists")
	}

	err := os.MkdirAll(data.Location, 0755)
	if err != nil {
		return err
	}

	tempDir, err := ioutil.TempDir(data.Location, "helm-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	client := action.NewPullWithOpts(action.WithConfig(&action.Configuration{}))
	client.RepoURL = MASSDRIVER_HELM_CHART_REPOSITORY
	client.ChartPathOptions = cpo
	client.Settings = cli.New()

	client.Untar = true
	client.UntarDir = tempDir //data.Location

	_, err = client.Run(data.Chart)
	if err != nil {
		return err
	}

	// regenerate Chart.yaml to match their config
	chart := ChartYaml{
		ApiVersion:  "v2",
		Name:        data.Name,
		Description: data.Description,
		Type:        "application",
		Version:     "1.0.0",
	}
	chartBytes, err := yaml.Marshal(chart)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(tempDir, data.Chart, "Chart.yaml"), chartBytes, 0644)
	if err != nil {
		return err
	}

	err = os.Rename(path.Join(tempDir, data.Chart), chartLocation)
	if err != nil {
		return err
	}

	return nil
}

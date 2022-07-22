package application

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
)

const MassdriverHelmChartRepository = "https://massdriver-cloud.github.io/helm-charts"

func Generate(data *TemplateData) error {
	cpo := action.ChartPathOptions{
		InsecureSkipTLSverify: true,
		RepoURL:               MassdriverHelmChartRepository,
		Version:               ">0.0.0-0",
	}

	// TODO: write template source as top-level in app.yaml
	// TODO: write cli-version as top-level in app.yaml

	chartLocation := path.Join(data.OutputDir, data.Name)

	if _, err := os.Stat(chartLocation); !os.IsNotExist(err) {
		return errors.New("specified directory already exists")
	}

	err := os.MkdirAll(data.OutputDir, common.AllRX|common.UserRW)
	if err != nil {
		return err
	}

	tempDir, err := ioutil.TempDir(data.OutputDir, "helm-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	client := action.NewPullWithOpts(action.WithConfig(&action.Configuration{}))
	client.RepoURL = MassdriverHelmChartRepository
	client.ChartPathOptions = cpo
	client.Settings = cli.New()

	client.Untar = true
	client.UntarDir = tempDir // data.OutputDir

	// _, err = client.Run(data.Chart)
	// if err != nil {
	// 	return err
	// }

	// regenerate Chart.yaml to match their config
	chart := ChartYAML{
		APIVersion:  "v2",
		Name:        data.Name,
		Description: data.Description,
		Type:        "application",
		Version:     "1.0.0",
	}
	chartBytes, err := yaml.Marshal(chart)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(tempDir, "", "Chart.yaml"), chartBytes, common.AllRead|common.UserRW)
	if err != nil {
		return err
	}

	err = os.Rename(path.Join(tempDir, ""), chartLocation)
	if err != nil {
		return err
	}

	return nil
}

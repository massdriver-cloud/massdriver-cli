package application

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

const MassdriverHelmChartRepository = "https://massdriver-cloud.github.io/helm-charts"
const MassdriverApplicationTemplatesRepository = "https://github.com/massdriver-cloud/application-templates"

func Generate(data *TemplateData) error {
	tempDir, err := ioutil.TempDir("/tmp/foo", "app-")
	if err != nil {
		return err
	}
	// defer os.RemoveAll(tempDir)

	_, cloneErr := git.PlainClone(tempDir, false, &git.CloneOptions{
		URL:      MassdriverApplicationTemplatesRepository,
		Progress: os.Stdout,
	})
	if cloneErr != nil {
		return errors.New("failed to clone application templates repository")
	}
	log.Info().Msg("starting copy")

	templateFiles, _ := fs.Sub(os.DirFS(tempDir), "kubernetes-deployment")
	// templateFiles := path.Join(tempDir, "kubernetes-deployment")
	err = fs.WalkDir(templateFiles, ".", func(filePath string, info fs.DirEntry, err error) error {
		log.Info().Msg("filePath" + filePath)
		if err != nil {
			return err
		}
		// outputPath := path.Join(data.Location, filePath)
		outputPath := path.Join(".", filePath)
		if info.IsDir() {
			if filePath == "." {
				// return os.MkdirAll(data.Location, common.AllRWX
				return os.MkdirAll(".", common.AllRWX)
			}

			return os.Mkdir(outputPath, common.AllRWX)
		}

		// var tmpl *template.Template
		// var outputFile *os.File
		// tmpl, _ = template.ParseFS(templateFiles, filePath)
		// outputFile, err = os.Create(outputPath)

		file, readErr := os.ReadFile(tempDir + "/kubernetes-deployment/" + filePath)
		if readErr != nil {
			panic(readErr)
		}

		writeErr := os.WriteFile(outputPath, file, common.AllRWX)
		if writeErr != nil {
			panic(writeErr)
			// return err
		}

		// defer outputFile.Close()
		// return tmpl.Execute(outputFile, data)
		return nil
	})
	if err != nil {
		return errors.New("failed to copy files")
	}

	// chartLocation := path.Join("app/chart/Chart.yaml")

	// if _, err := os.Stat(chartLocation); !os.IsNotExist(err) {
	// 	return errors.New("specified directory already exists")
	// }

	// err = os.MkdirAll(data.Location, common.AllRX|common.UserRW)
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
		panic("marshal error")
		return err
	}
	err = ioutil.WriteFile(path.Join(tempDir, "kubernetes-deployment/app/chart", "Chart.yaml"), chartBytes, common.AllRead|common.UserRW)
	if err != nil {
		return err
	}

	err = os.Rename(path.Join(tempDir, "kubernetes-deployment/app/chart/Chart.yaml"), "app/chart/Chart.yaml")
	if err != nil {
		return err
	}

	return nil
}

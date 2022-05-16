package generator

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
)

type TemplateData struct {
	Name        string
	Description string
	Access      string
	Type        string
	TemplateDir string
	BundleDir   string
}

func Generate(data *TemplateData) error {
	bundleDir := fmt.Sprintf("%s/%s", data.BundleDir, data.Name)
	currentDirectory := ""

	err := filepath.WalkDir(data.TemplateDir, func(path string, info fs.DirEntry, err error) error {

		if info.IsDir() {
			if isRootPath(path, data.TemplateDir) {
				os.MkdirAll(bundleDir, 0777)
				return nil
			}

			subDirectory := fmt.Sprintf("%s/%s", bundleDir, info.Name())
			os.Mkdir(subDirectory, 0777)
			currentDirectory = fmt.Sprintf("%s/", info.Name())
			return nil
		}

		renderPath := fmt.Sprintf("%s/%s%s", bundleDir, currentDirectory, info.Name())
		err = renderTemplate(path, renderPath, data)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func renderTemplate(path, renderPath string, data *TemplateData) error {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return err
	}

	fileToWrite, err := os.Create(renderPath)
	if err != nil {
		return err
	}

	tmpl.Execute(fileToWrite, data)

	fileToWrite.Close()

	return nil
}

func isRootPath(rootPath, currentPath string) bool {
	return rootPath == currentPath
}

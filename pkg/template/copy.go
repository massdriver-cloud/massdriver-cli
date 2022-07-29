package template

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"github.com/rs/zerolog/log"
)

func Copy(templateDir string, data *Data) error {
	templateName := data.TemplateName
	templatePath := path.Join(templateDir, templateName)

	return copyTemplate(templatePath, data, func(filePath string) ([]byte, error) {
		return os.ReadFile(path.Join(templatePath, filePath))
	})
}

func CopyFS(templateFiles fs.FS, data *Data) error {
	return copyTemplateFS(templateFiles, data, func(filePath string) ([]byte, error) {
		file, err := templateFiles.Open(filePath)
		if err != nil {
			return nil, err
		}
		return io.ReadAll(file)
	})
}

func copyTemplate(templateDir string, data *Data, tmplFunc func(path string) ([]byte, error)) error {
	outputDir := data.OutputDir

	return filepath.Walk(templateDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return walkDir(outputDir, filePath, info.IsDir(), data, tmplFunc)
		// outputPath := path.Join(outputDir, filePath)
		// if info.IsDir() {
		// 	if filePath == "." {
		// 		return os.MkdirAll(".", common.AllRWX)
		// 	}

		// 	return os.Mkdir(outputPath, common.AllRWX)
		// }

		// var tmpl *template.Template
		// file, errTmpl := tmplFunc(filePath)
		// if errTmpl != nil {
		// 	return errTmpl
		// }

		// tmpl, err = template.New("tmpl").Delims(OpenPattern, ClosePattern).Parse(string(file))

		// if _, err = os.Stat(outputPath); err == nil {
		// 	fmt.Printf("%s exists. Overwrite? (y|N): ", outputPath)
		// 	var response string
		// 	fmt.Scanln(&response)

		// 	if response == "y" || response == "Y" || response == "yes" {
		// 		return writeFile(outputPath, tmpl, data)
		// 	}
		// }
		// log.Info().Msg("writing file")

		// return checkAndWrite(file, tmpl, data, outputPath)
	})
}

func copyTemplateFS(templateFiles fs.FS, data *Data, tmplFunc func(path string) ([]byte, error)) error {
	// templateName := data.TemplateName
	outputDir := data.OutputDir

	return fs.WalkDir(templateFiles, ".", func(filePath string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// outputPath := path.Join(outputDir, filePath)
		return walkDir(outputDir, filePath, info.IsDir(), data, tmplFunc)

		// if info.IsDir() {
		// 	if filePath == "." {
		// 		return os.MkdirAll(".", common.AllRWX)
		// 	}

		// 	return os.Mkdir(outputPath, common.AllRWX)
		// }

		// var tmpl *template.Template
		// file, errTmpl := tmplFunc(filePath)
		// if errTmpl != nil {
		// 	return errTmpl
		// }
		// return checkAndWrite(file, tmpl, data, outputPath)
		// tmpl, err = template.New("tmpl").Delims(OpenPattern, ClosePattern).Parse(string(file))

		// if _, err = os.Stat(outputPath); err == nil {
		// 	fmt.Printf("%s exists. Overwrite? (y|N): ", outputPath)
		// 	var response string
		// 	fmt.Scanln(&response)

		// 	if response == "y" || response == "Y" || response == "yes" {
		// 		return writeFile(outputPath, tmpl, data)
		// 	}
		// }
		// log.Info().Msg("writing file")

		// return writeFile(outputPath, tmpl, data)
	})
}

func walkDir(outputDir string, filePath string, isDir bool, data *Data, tmplFunc func(path string) ([]byte, error)) error {
	outputPath := path.Join(outputDir, filePath)

	if isDir {
		if filePath == "." {
			return os.MkdirAll(".", common.AllRWX)
		}

		return os.Mkdir(outputPath, common.AllRWX)
	}

	var tmpl *template.Template
	file, errTmpl := tmplFunc(filePath)
	if errTmpl != nil {
		return errTmpl
	}
	return checkAndWrite(file, tmpl, data, outputPath)
}

func checkAndWrite(file []byte, tmpl *template.Template, data *Data, outputPath string) error {
	tmpl, err := template.New("tmpl").Delims(OpenPattern, ClosePattern).Parse(string(file))

	if _, err = os.Stat(outputPath); err == nil {
		fmt.Printf("%s exists. Overwrite? (y|N): ", outputPath)
		var response string
		fmt.Scanln(&response)

		if response == "y" || response == "Y" || response == "yes" {
			return writeFile(outputPath, tmpl, data)
		}
	}
	log.Info().Msg("writing file")

	return writeFile(outputPath, tmpl, data)
}

func writeFile(outputPath string, tmpl *template.Template, data *Data) error {
	outputFile, err := os.Create(outputPath)

	if err != nil {
		return err
	}

	defer outputFile.Close()
	err = tmpl.Execute(outputFile, data)
	if err != nil {
		log.Info().Msgf("error executing template: %s", err)
	}
	return err
}

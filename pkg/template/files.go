package template

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"text/template"

	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
)

func WriteToFile(filePath string, tmpl *template.Template, data *Data) error {
	outputFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	return tmpl.Execute(outputFile, data)
}

func readFileFromEmbededFunc(templateFS fs.FS) func(filePath string) ([]byte, error) {
	return func(filePath string) ([]byte, error) {
		file, err := templateFS.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		return io.ReadAll(file)
	}
}

func readFileFunc(dir string) func(filePath string) ([]byte, error) {
	return func(filePath string) ([]byte, error) {
		return os.ReadFile(path.Join(dir, filePath))
	}
}

func mkDirOrWriteFile(outputDir string, filePath string, info fs.DirEntry, data *Data, fileReadFunc func(path string) ([]byte, error)) error {
	outputPath := path.Join(outputDir, filePath)

	if info.IsDir() {
		if filePath == "." {
			return os.MkdirAll(".", common.AllRWX)
		}

		return os.Mkdir(outputPath, common.AllRWX)
	}

	file, err := fileReadFunc(filePath)
	if err != nil {
		return err
	}
	return promptAndWrite(file, data, outputPath)
}

func promptAndWrite(file []byte, data *Data, outputPath string) error {
	tmpl, errTmpl := template.New("tmpl").Delims(OpenPattern, ClosePattern).Parse(string(file))
	if errTmpl != nil {
		return errTmpl
	}

	if _, err := os.Stat(outputPath); err == nil {
		fmt.Printf("%s exists. Overwrite? (y|N): ", outputPath)
		var response string
		fmt.Scanln(&response)

		if response == "y" || response == "Y" || response == "yes" {
			return WriteToFile(outputPath, tmpl, data)
		}
	}

	return WriteToFile(outputPath, tmpl, data)
}

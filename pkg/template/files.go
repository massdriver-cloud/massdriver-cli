package template

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"

	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"github.com/osteele/liquid"
)

func WriteToFile(filePath string, template []byte, data *Data) error {
	engine := liquid.NewEngine()

	var bindings map[string]interface{}
	inrec, _ := json.Marshal(data)
	json.Unmarshal(inrec, &bindings)

	out, renderErr := engine.ParseAndRender(template, bindings)

	if renderErr != nil {
		fmt.Printf("hey its here %s", filePath)
		return renderErr
	}

	return os.WriteFile(filePath, out, 0600)

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
	// create the output dir if it doesn't exist
	if outputDir != "" && outputDir != "." {
		if _, err := os.Stat(outputDir); err != nil {
			return os.Mkdir(outputDir, common.AllRWX)
		}
	}
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

func promptAndWrite(template []byte, data *Data, outputPath string) error {
	if _, err := os.Stat(outputPath); err == nil {
		fmt.Printf("%s exists. Overwrite? (y|N): ", outputPath)
		var response string
		fmt.Scanln(&response)

		if response == "y" || response == "Y" || response == "yes" {
			return WriteToFile(outputPath, template, data)
		} else {
			return nil
		}
	}

	return WriteToFile(outputPath, template, data)
}

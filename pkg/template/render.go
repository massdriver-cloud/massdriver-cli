package template

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	cp "github.com/otiai10/copy"
)

// When we generate code from templates, we only do a liquid pass on top level text files (massdriver.yaml, readme, etc)
// and copy everything else. A users 'src' is either helm or terraform, they can programmatically do whatever they want in there
// the generator should just be configuring input values that lets them do their own logic.

// Example: We want to have a template that has database migrations and an 'enableMigration' field.
// We want to render the migration code (k8s job) either way, but we'll set the enableMigration job to false so
// helm won't provision it. This allows the user to enable migrations later w/o having to re-render.

func RenderDirectory(templateDir string, data *Data) error {
	files, readDirErr := ioutil.ReadDir(templateDir)
	if readDirErr != nil {
		return readDirErr
	}

	if _, checkDirExistsErr := os.Stat(data.OutputDir); errors.Is(checkDirExistsErr, os.ErrNotExist) {
		mkdirErr := os.MkdirAll(data.OutputDir, os.ModePerm)
		if mkdirErr != nil {
			return mkdirErr
		}
	}

	for _, file := range files {
		srcPath := templateDir + "/" + file.Name()
		destPath := data.OutputDir + "/" + file.Name()

		if _, err := os.Stat(destPath); err == nil {
			fmt.Printf("%s exists. Overwrite? (y|N): ", destPath)
			var response string
			fmt.Scanln(&response)

			if response != "y" && response != "Y" && response != "yes" {
				continue
			}
		}

		if file.IsDir() {
			err := cp.Copy(srcPath, destPath)
			if err != nil {
				fmt.Printf("Error!: %v", err)
				return err
			}
		} else {
			contents, err := os.ReadFile(srcPath)
			if err != nil {
				return err
			}

			err = WriteToFile(destPath, contents, data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

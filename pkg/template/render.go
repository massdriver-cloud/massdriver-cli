package template

import (
	"fmt"
	"io/fs"
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
	files, err := ioutil.ReadDir(templateDir)

	if err != nil {
		return err
	}

	for _, file := range files {
		srcPath := templateDir + "/" + file.Name()
		destPath := data.OutputDir + "/" + file.Name()

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

func RenderEmbededDirectory(templateFS fs.FS, data *Data) error {
	return renderEmbededDirectory(templateFS, data, readFileFromEmbededFunc(templateFS))
}

func renderEmbededDirectory(templateFiles fs.FS, data *Data, tmplFunc func(path string) ([]byte, error)) error {
	return fs.WalkDir(templateFiles, ".", func(filePath string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		return mkDirOrWriteFile(data.OutputDir, filePath, info, data, tmplFunc)
	})
}

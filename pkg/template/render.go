package template

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func RenderDirectory(templateDir string, data *Data) error {
	return renderDirectory(templateDir, data, readFileFunc(templateDir))
}

func RenderEmbededDirectory(templateFS fs.FS, data *Data) error {
	return renderEmbededDirectory(templateFS, data, readFileFromEmbededFunc(templateFS))
}

func renderDirectory(templateDir string, data *Data, tmplFunc func(path string) ([]byte, error)) error {
	return filepath.WalkDir(templateDir, func(filePath string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		newPath := strings.Replace(filePath, templateDir, "", 1)
		if newPath == "" {
			newPath = "."
		}
		return mkDirOrWriteFile(data.OutputDir, newPath, info, data, tmplFunc)
	})
}

func renderEmbededDirectory(templateFiles fs.FS, data *Data, tmplFunc func(path string) ([]byte, error)) error {
	return fs.WalkDir(templateFiles, ".", func(filePath string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		return mkDirOrWriteFile(data.OutputDir, filePath, info, data, tmplFunc)
	})
}

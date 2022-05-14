package compress

import (
	"archive/tar"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func TarDirectory(dirPath string, prefix string, tarWriter *tar.Writer) error {

	absolutePath, err := filepath.Abs(dirPath)
	if err != nil {
		return err
	}

	// walk through every file in the folder
	filepath.Walk(absolutePath, func(file string, fi os.FileInfo, err error) error {
		relativeDestinationFilePath := path.Join(prefix, strings.TrimPrefix(filepath.ToSlash(file), absolutePath))
		// generate tar header
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}
		if ShouldIgnore(relativeDestinationFilePath) {
			return nil
		}

		// must provide real name
		// (see https://golang.org/src/archive/tar/common.go?#L626)
		header.Name = relativeDestinationFilePath

		// write header
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}
		// if not a dir, write file content
		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tarWriter, data); err != nil {
				return err
			}
		}
		return nil
	})

	return nil
}

func TarFile(filePath string, prefix string, tarWriter *tar.Writer) error {

	filePtr, err := os.Open(filePath)
	if err != nil {
		return err
	}

	fileInfo, err := filePtr.Stat()
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return errors.New("specified path is not a file")
	}

	header, err := tar.FileInfoHeader(fileInfo, filePath)
	if err != nil {
		return err
	}

	// must provide real name
	// (see https://golang.org/src/archive/tar/common.go?#L626)
	header.Name = path.Join(prefix, filepath.Base(filePath))

	// write header
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}
	// if not a dir, write file content

	if _, err := io.Copy(tarWriter, filePtr); err != nil {
		return err
	}

	return nil
}

func ShouldIgnore(relativeFilePath string) bool {
	// .terraform, .terraform.lock.hcl
	return strings.HasPrefix(relativeFilePath, "bundle/src/.terraform") ||
		  // terraform.tfstate, terraform.tfstate.backup, etc...
			strings.Contains(relativeFilePath, ".tfstate") ||
			// Massdriver gives the vars
			strings.Contains(relativeFilePath, ".tfvars") ||
			strings.Contains(relativeFilePath, ".md")
}

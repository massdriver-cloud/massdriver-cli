package compress

import (
	"archive/tar"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
)

func TarDirectory(dirPath string, prefix string, tarWriter *tar.Writer) error {
	absolutePath, pathErr := filepath.Abs(dirPath)
	if pathErr != nil {
		return pathErr
	}

	// walk through every file in the folder
	if err := filepath.Walk(absolutePath, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
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
		if writeErr := tarWriter.WriteHeader(header); writeErr != nil {
			return writeErr
		}
		// if not a dir, write file content
		if !fi.IsDir() {
			data, openErr := os.Open(file)
			if openErr != nil {
				return openErr
			}
			if _, copyErr := io.Copy(tarWriter, data); copyErr != nil {
				return copyErr
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func TarFile(filePath string, prefix string, tarWriter *tar.Writer) error {
	filePtr, errOpen := os.Open(filePath)
	if errOpen != nil {
		return errOpen
	}

	fileInfo, errStat := filePtr.Stat()
	if errStat != nil {
		return errStat
	}

	if fileInfo.IsDir() {
		return errors.New("specified path is not a file")
	}

	header, errHeader := tar.FileInfoHeader(fileInfo, filePath)
	if errHeader != nil {
		return errHeader
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
	// full filenames to ignore
	if common.Contains(common.FileIgnores, relativeFilePath) {
		return true
	}

	// partial filenames to ignore
	// .terraform, .terraform.lock.hcl
	return strings.Contains(relativeFilePath, ".terraform") ||
		// terraform.tfstate, terraform.tfstate.backup, etc...
		strings.Contains(relativeFilePath, ".tfstate") ||
		// Massdriver gives the vars
		strings.Contains(relativeFilePath, ".tfvars") ||
		strings.Contains(relativeFilePath, ".md")
}

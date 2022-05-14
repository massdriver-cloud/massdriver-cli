package compress

import (
	"archive/tar"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"fmt"
	"strings"
)

func TarDirectory(relativeSourcePath string, destinationPrefix string, tarWriter *tar.Writer) error {
	absoluteSourcePath, err := filepath.Abs(relativeSourcePath)
	if err != nil {
		return err
	}

	// walk through every absolueSourceFilePath in the folder
	filepath.Walk(absoluteSourcePath, func(absolueSourceFilePath string, info os.FileInfo, err error) error {
		relativeDestinationFilePath := path.Join(destinationPrefix, strings.TrimPrefix(filepath.ToSlash(absolueSourceFilePath), absoluteSourcePath))

		// this feels a little weird, "ideally" we pass in the relativeSourceFilePath
		// however, we'd only use that here and it would only make this test slightly more sane
		if info.IsDir() || ShouldIgnore(relativeDestinationFilePath) {
			return nil
		}

		// Open the source file which will be written into the archive
		data, err := os.Open(absolueSourceFilePath)
		if err != nil {
			return err
		}
		defer data.Close()

		// Create a tar Header from the FileInfo data
		// info.Name() is the file name without any directory information
		// aka chart.yaml _not_ chart/chart.yaml
		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		// Use full path as name (FileInfoHeader only takes the basename)
		// If we don't do this the directory strucuture would
		// not be preserved
		// https://golang.org/src/archive/tar/common.go?#L626
		header.Name = relativeDestinationFilePath

		err = tarWriter.WriteHeader(header)
		if err != nil {
			return err
		}

		// write soure file content into tar
		_, err = io.Copy(tarWriter, data)
		if err != nil {
			return err
		}

		return nil
	})

	return nil
}

func ShouldIgnore(relativeFilePath string) bool {
	return strings.HasPrefix(relativeFilePath, "bundle/src/.terraform") ||
			strings.HasPrefix(relativeFilePath, "bundle/src/terraform") ||
			strings.HasPrefix(relativeFilePath, "bundle/src/connections.tfvars") ||
			strings.HasPrefix(relativeFilePath, "bundle/src/params.tfvars") ||
			strings.HasPrefix(relativeFilePath, "bundle/schema") ||
			strings.HasSuffix(relativeFilePath, ".md")
}

func TarFile(filePath string, destinationPrefix string, tarWriter *tar.Writer) error {

	filePtr, err := os.Open(filePath)
	if err != nil {
		return err
	}

	fileInfo, err := filePtr.Stat()
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return errors.New("specified path is not a absolueSourceFilePath")
	}

	header, err := tar.FileInfoHeader(fileInfo, filePath)
	if err != nil {
		return err
	}

	// must provide real name
	// (see https://golang.org/src/archive/tar/common.go?#L626)
	header.Name = path.Join(destinationPrefix, filepath.Base(filePath))

	// write header
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}
	// if not a dir, write absolueSourceFilePath content

	if _, err := io.Copy(tarWriter, filePtr); err != nil {
		return err
	}

	return nil
}

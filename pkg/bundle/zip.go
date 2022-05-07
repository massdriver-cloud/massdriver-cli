package bundle

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func TarGzipBundle(filePath string, buf io.Writer) error {
	// tar > gzip > buf
	gzipWriter := gzip.NewWriter(buf)
	tarWriter := tar.NewWriter(gzipWriter)

	absolutePath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}
	dirPath := path.Dir(absolutePath)

	// walk through every file in the folder
	filepath.Walk(dirPath, func(file string, fi os.FileInfo, err error) error {
		// generate tar header
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		// must provide real name
		// (see https://golang.org/src/archive/tar/common.go?#L626)
		header.Name = path.Join("bundle", strings.TrimPrefix(filepath.ToSlash(file), dirPath))

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

	// produce tar
	if err := tarWriter.Close(); err != nil {
		return err
	}
	// produce gzip
	if err := gzipWriter.Close(); err != nil {
		return err
	}

	return nil
}

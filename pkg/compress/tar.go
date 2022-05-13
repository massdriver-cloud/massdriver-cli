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
		// generate tar header
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		// must provide real name
		// (see https://golang.org/src/archive/tar/common.go?#L626)
		header.Name = path.Join(prefix, strings.TrimPrefix(filepath.ToSlash(file), absolutePath))

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
		} else if fi.IsDir() && fi.Name() == "chart" {
			packageChart(path.Join(path.Dir(appPath), app.Deployment.Path), path.Join(workingDir, "chart"))

			filepath.Walk(absolutePath + "/chart", func(file string, fi os.FileInfo, err error) error {
				if !fi.IsDir() {
					data, err := os.Open(file)
					if err != nil {
						return err
					}
					if _, err := io.Copy(tarWriter, data); err != nil {
						return err
					}
				}
			}
		}
		return nil
	})

	return nil
}

// func packageChart(chartPath string, destPath string) error {
// 	var err error = filepath.Walk(chartPath, func(path string, info os.FileInfo, err error) error {
// 		var relPath string = strings.TrimPrefix(path, chartPath)
// 		if relPath == "" {
// 			return nil
// 		}
// 		if info.IsDir() {
// 			return os.Mkdir(filepath.Join(destPath, relPath), 0755)
// 		} else {
// 			var data, err1 = ioutil.ReadFile(filepath.Join(chartPath, relPath))
// 			if err1 != nil {
// 				return err1
// 			}
// 			return ioutil.WriteFile(filepath.Join(destPath, relPath), data, 0777)
// 		}
// 	})
// 	return err
// }

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

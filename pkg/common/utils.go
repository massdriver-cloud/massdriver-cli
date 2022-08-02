package common

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func CopyFolder(sourcePath string, destPath string) error {
	err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		var relPath = strings.TrimPrefix(path, sourcePath)
		if relPath == "" {
			return nil
		}
		if info.IsDir() {
			return os.Mkdir(filepath.Join(destPath, relPath), AllRX|UserRW)
		}
		var data, err1 = ioutil.ReadFile(filepath.Join(sourcePath, relPath))
		if err1 != nil {
			return err1
		}
		return ioutil.WriteFile(filepath.Join(destPath, relPath), data, AllRWX)
	})
	return err
}

func WriteFile(filePath string, data []byte, errMarshal error) error {
	if errMarshal != nil {
		return errMarshal
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, errWrite := file.Write(data); errWrite != nil {
		return errWrite
	}
	return nil
}

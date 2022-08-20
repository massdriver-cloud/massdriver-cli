package common

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
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
	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		var relPath = strings.TrimPrefix(path, sourcePath)
		if relPath == "" {
			return nil
		}

		if shouldNotInclude(info) {
			return nil
		}

		if info.IsDir() {
			log.Info().Msgf("mkdir: %s", relPath)
			return os.Mkdir(filepath.Join(destPath, relPath), AllRX|UserRW)
		}

		log.Info().Msgf("copying: %s", relPath)
		var data, err1 = ioutil.ReadFile(filepath.Join(sourcePath, relPath))
		if err1 != nil {
			return err1
		}

		return ioutil.WriteFile(filepath.Join(destPath, relPath), data, AllRWX)
	})
}

func WriteFile(filePath string, data []byte, errToBytes error) error {
	if errToBytes != nil {
		return errToBytes
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

func shouldNotInclude(info os.FileInfo) bool {
	fileName := info.Name()
	for _, ignore := range FileIgnores {
		if strings.Contains(fileName, ignore) {
			return true
		}
	}

	bytes := info.Size()
	kilobytes := (bytes / 1024)
	var megabytes float64
	megabytes = (float64)(kilobytes / 1024)
	if megabytes > 10 {
		return true
	}

	return false
}

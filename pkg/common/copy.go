package common

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

type CopyConfig struct {
	Allows  []string
	Ignores []string
}

// public private pair dog
func CopyFolder(src, dst string, config *CopyConfig) error {
	files, errReadDir := ioutil.ReadDir(src)
	if errReadDir != nil {
		return errReadDir
	}

	for _, fileInfo := range files {
		fileOrDirName := fileInfo.Name()
		log.Info().Msgf("fileOrDirName	: %s", fileOrDirName)

		if !shouldInclude(fileOrDirName, config) {
			continue
		}

		// recurse
		errCopy := copyFolder(path.Join(src, fileOrDirName), path.Join(dst, fileOrDirName), config)
		if errCopy != nil {
			return errCopy
		}
	}

	return nil
}

// if it's a file
//   base case
// make the folder in dest
//   pass the files in that folder to this function
// if file, write file
func copyFolder(sourcePath string, destPath string, config *CopyConfig) error {
	fileInfo, _ := os.Stat(sourcePath)
	if shouldIgnore(fileInfo, config) {
		return nil
	}

	if !fileInfo.IsDir() {
		// base case, write the file
		data, err1 := ioutil.ReadFile(sourcePath)
		if err1 != nil {
			return err1
		}

		return ioutil.WriteFile(destPath, data, AllRWX)
	}

	errMkdir := os.Mkdir(destPath, AllRX|UserRW)
	if errMkdir != nil {
		return errMkdir
	}

	files, errReadDir := ioutil.ReadDir(sourcePath)
	if errReadDir != nil {
		return errReadDir
	}
	for _, subDirFileInfo := range files {
		fileOrDirName := subDirFileInfo.Name()

		if shouldIgnore(subDirFileInfo, config) {
			continue
		}

		// recurse
		errCopy := copyFolder(filepath.Join(sourcePath, fileOrDirName), filepath.Join(destPath, fileOrDirName), config)
		if errCopy != nil {
			return errCopy
		}
	}

	return nil
}

func shouldInclude(fileOrDirName string, conf *CopyConfig) bool {
	for _, allow := range conf.Allows {
		if strings.Contains(fileOrDirName, allow) {
			log.Info().Msgf("including: %s", fileOrDirName)
			return true
		}
	}
	return false
}

const MaxFileSizeMB = 10
const tenTwentyFour = 1024

func shouldIgnore(info fs.FileInfo, config *CopyConfig) bool {
	filePath := info.Name()

	for _, ignore := range config.Ignores {
		if strings.Contains(filePath, ignore) {
			log.Info().Msgf("ignoring	: %s", filePath)
			return true
		}
	}

	bytes := info.Size()
	kilobytes := (bytes / tenTwentyFour)
	megabytes := (float64)(kilobytes / tenTwentyFour)

	if megabytes > MaxFileSizeMB {
		log.Info().Msgf("megabytes: %v", megabytes)
		return true
	}

	return false
}

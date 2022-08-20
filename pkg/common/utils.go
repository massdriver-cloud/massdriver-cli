package common

import (
	"io/fs"
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

func CopyFolder(sourcePath string, destPath string, allowList []string) error {
	// if it's a folder
	// make the folder in dest
	//   pass the files in that folder to this function
	// if file, write file
	file, _ := os.Stat(sourcePath)

	if shouldIgnore(file) {
		return nil
	}

	if file.IsDir() {
		errMkdir := os.Mkdir(destPath, AllRX|UserRW)
		if errMkdir != nil {
			return errMkdir
		}

		files, errReadDir := ioutil.ReadDir(sourcePath)
		if errReadDir != nil {
			return errReadDir
		}
		for _, fileInfo := range files {
			fileOrDirName := fileInfo.Name()

			log.Info().Msgf("fileOrDirName: %s", fileOrDirName)
			if shouldIgnore(fileInfo) {
				continue
			}

			errCopy := CopyFolder(filepath.Join(sourcePath, fileOrDirName), filepath.Join(destPath, fileOrDirName), allowList)
			if errCopy != nil {
				return errCopy
			}
		}
	} else {
		fileOrDirName := file.Name()
		log.Info().Msgf("writing: %s", fileOrDirName)
		var data, err1 = ioutil.ReadFile(sourcePath)
		if err1 != nil {
			return err1
		}

		return ioutil.WriteFile(destPath, data, AllRWX)
	}

	return nil
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

func shouldIgnore(info fs.FileInfo) bool {
	filePath := info.Name()
	log.Info().Msgf("shouldIgnore: %s", filePath)

	for _, ignore := range FileIgnores {
		if strings.Contains(filePath, ignore) {
			log.Info().Msgf("ignoring	: %s", filePath)
			return true
		}
	}

	bytes := info.Size()
	kilobytes := (bytes / 1024)
	var megabytes float64
	megabytes = (float64)(kilobytes / 1024)

	if megabytes > 10 {
		log.Info().Msgf("megabytes: %v", megabytes)
		return true
	}

	return false
}

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

type CopyStats struct {
	FolderSize int64
}

func CopyFolder(src, dst string, config *CopyConfig) (CopyStats, error) {
	stats := CopyStats{
		FolderSize: 0,
	}
	files, errReadDir := ioutil.ReadDir(src)
	if errReadDir != nil {
		return stats, errReadDir
	}

	for _, fileInfo := range files {
		name := fileInfo.Name()
		if !shouldInclude(name, config) {
			continue
		}

		// recursively copy the allowed files / folders
		errCopy := copyFolder(path.Join(src, name), path.Join(dst, name), config, &stats)
		if errCopy != nil {
			return stats, errCopy
		}
	}

	return stats, nil
}

// Written because path.Walk traverses the entire directory tree.
// If a step has a folder like .terraform, path.Walk will still traverse it
//
// base case: a file or folder we should ignore, skip
// base case: a file we should include, write file
// resurce when: it's a folder we should include
func copyFolder(src, dest string, config *CopyConfig, stats *CopyStats) error {
	info, _ := os.Stat(src)
	// a file or folder we should ignore, skip
	if shouldIgnore(info, config) {
		return nil
	}

	// a file we should include, write file
	if !info.IsDir() {
		stats.FolderSize += info.Size()
		data, err1 := ioutil.ReadFile(src)
		if err1 != nil {
			return err1
		}

		return ioutil.WriteFile(dest, data, AllRWX)
	}

	// a folder we should include
	// so we create the folder, then iterate through
	// the files in that folder.
	errMkdir := os.Mkdir(dest, AllRX|UserRW)
	if errMkdir != nil {
		return errMkdir
	}

	files, errReadDir := ioutil.ReadDir(src)
	if errReadDir != nil {
		return errReadDir
	}
	for _, subDirFileInfo := range files {
		// recurse
		name := subDirFileInfo.Name()
		errCopy := copyFolder(filepath.Join(src, name), filepath.Join(dest, name), config, stats)
		if errCopy != nil {
			return errCopy
		}
	}

	return nil
}

func shouldInclude(fileOrDirName string, conf *CopyConfig) bool {
	for _, allow := range conf.Allows {
		if strings.Contains(fileOrDirName, allow) {
			return true
		}
	}
	return false
}

func shouldIgnore(info fs.FileInfo, config *CopyConfig) bool {
	fileName := info.Name()

	for _, ignore := range config.Ignores {
		if strings.Contains(fileName, ignore) {
			return true
		}
	}

	mbs := FileSizeMB(info.Size())
	if mbs > MaxFileSizeMB {
		log.Debug().Msgf("File: %s is larger than limit of %vMB.", fileName, MaxFileSizeMB)
		return true
	}

	return false
}

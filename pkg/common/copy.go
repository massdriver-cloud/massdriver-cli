package common

import (
	"io/fs"
	"io/ioutil"
	"os"
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

func CopyFolder(src, dstDir string, config *CopyConfig) (CopyStats, error) {
	stats := CopyStats{
		FolderSize: 0,
	}

	err := filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == "." {
			return nil
		}

		relPath := strings.Replace(path, src, "", 1)
		if src == "." {
			relPath = "/" + path
		}

		depth := strings.Count(relPath, string(os.PathSeparator))
		skip, errSkip := shouldSkip(info, depth, config)
		if errSkip != nil {
			return errSkip
		}
		if skip {
			return nil
		}

		writePath := filepath.Abs(filepath.Join(dstDir, relPath))
		if info.IsDir() {
			errMkdir := os.Mkdir(writePath, info.Mode())
			if errMkdir != nil {
				return errMkdir
			}
		} else {
			stats.FolderSize += info.Size()
			data, err1 := ioutil.ReadFile(path)
			if err1 != nil {
				return err1
			}

			return ioutil.WriteFile(writePath, data, info.Mode())
		}

		return nil
	})

	return stats, err
}

func shouldSkip(info fs.FileInfo, depth int, config *CopyConfig) (bool, error) {
	name := info.Name()

	// the parent folder, aka k8s-bundle
	// we only want what's inside
	if depth == 0 {
		return true, nil
	}

	// subfolders, aka src, core-services, etc
	// if we're at the root of the bundle
	// we only want to honor the include list
	if depth == 1 && !shouldInclude(name, config) {
		if info.IsDir() {
			return true, filepath.SkipDir
		}
		return true, nil
	}

	// inside bundle directories like src, core-services, etc
	// we want to include every file _except_ the ones
	// that match the ignore _criteria_. File names, sizes, etc...
	if shouldIgnore(info, config) {
		if info.IsDir() {
			return true, filepath.SkipDir
		}
		return true, nil
	}
	return false, nil
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
		log.Error().Msgf("File: %s is larger than limit of %vMB.", fileName, MaxFileSizeMB)
		return true
	}

	return false
}

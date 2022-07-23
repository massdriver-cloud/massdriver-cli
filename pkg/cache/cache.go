package cache

import (
	"os"
	"os/user"
	"path/filepath"
)

func Dir() string {
	dir, _ := cacheDir()
	return dir
}

func cacheDir() (string, error) {
	usr, _ := user.Current()
	dir := usr.HomeDir
	cacheDir := filepath.Join(dir, ".massdriver")
	if _, err := os.Stat(cacheDir); !os.IsNotExist(err) {
		return cacheDir, err
	}

	if errMkdir := os.Mkdir(cacheDir, 0755); errMkdir != nil {
		return cacheDir, errMkdir
	}
	return cacheDir, nil
}

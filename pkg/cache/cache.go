package cache

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
)

const (
	cacheDir         = "/tmp/massdriver-cli"
	TemplateCacheDir = "/tmp/massdriver-cli/templates"
)

func checkCacheDir() error {
	if _, err := os.Stat(cacheDir); !os.IsNotExist(err) {
		return nil
	}

	err := os.Mkdir(cacheDir, 0755)
	if err != nil {
		return err
	}
	return nil
}

func GetMassdriverTemplates() error {
	err := checkCacheDir()
	if _, err = os.Stat(TemplateCacheDir); !os.IsNotExist(err) {
		return nil
	}

	err = os.Mkdir(TemplateCacheDir, 0755)
	if err != nil {
		return err
	}
	_, cloneErr := git.PlainClone(TemplateCacheDir, false, &git.CloneOptions{
		URL:      common.MassdriverApplicationTemplatesRepository,
		Progress: os.Stdout,
		Depth:    1,
	})
	if cloneErr != nil {
		return cloneErr
	}
	return nil
}

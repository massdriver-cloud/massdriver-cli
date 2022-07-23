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
		return err
	}

	if errMkdir := os.Mkdir(cacheDir, 0755); errMkdir != nil {
		return errMkdir
	}
	return nil
}

func GetMassdriverTemplates() error {
	if err := checkCacheDir(); err != nil {
		return err
	}
	if _, errTemplateDir := os.Stat(TemplateCacheDir); !os.IsNotExist(errTemplateDir) {
		return nil
	}

	errMkdir := os.Mkdir(TemplateCacheDir, 0755)
	if errMkdir != nil {
		return errMkdir
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

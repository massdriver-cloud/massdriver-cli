package cache

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/rs/zerolog/log"
)

// TODO: pull from massrc
const MassdriverApplicationTemplatesRepository = "https://github.com/massdriver-cloud/application-templates"

func AppTemplateCacheDir() string {
	dir, _ := templateCacheDir()
	return dir
}

func ApplicationTemplates() ([]string, error) {
	templates := []string{}
	templateDirs, err := ioutil.ReadDir(AppTemplateCacheDir())
	if err != nil {
		return templates, err
	}
	for _, f := range templateDirs {
		// all directories that aren't .git
		// cheap way of listing templates
		if f.IsDir() && f.Name() != ".git" {
			templates = append(templates, f.Name())
		}
	}
	return templates, nil
}

func RefreshAppTemplates() error {
	if err := clearAppTemplateCache(); err != nil {
		return err
	}
	return downloadAppTemplates()
}

func templateCacheDir() (string, error) {
	cacheDir, err := cacheDir()
	if err != nil {
		return "", err
	}

	templateDir := filepath.Join(cacheDir, "templates")
	if _, errDir := os.Stat(templateDir); !os.IsNotExist(errDir) {
		return templateDir, errDir
	}

	if errMkdir := os.Mkdir(templateDir, 0755); errMkdir != nil {
		return templateDir, errMkdir
	}
	return templateDir, nil
}

func clearAppTemplateCache() error {
	if err := os.RemoveAll(AppTemplateCacheDir()); err != nil {
		return err
	}
	return nil
}

func downloadAppTemplates() error {
	templateCacheDir := AppTemplateCacheDir()
	log.Debug().Msgf("Downloading templates to %s", templateCacheDir)
	// log.Debug().Msgf("Cloning templates from %s", common.Config().Application.Templates.Repository)
	_, cloneErr := git.PlainClone(templateCacheDir, false, &git.CloneOptions{
		// URL:      common.Config().Application.Templates.Repository,
		URL:      MassdriverApplicationTemplatesRepository,
		Progress: os.Stdout,
		Depth:    1,
	})
	if cloneErr != nil {
		return cloneErr
	}
	return nil
}

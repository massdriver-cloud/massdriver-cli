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

// AppTemplateCacheDir is a reader function to access the local cache of templates.
// When developing templates, the cache source can be overwritten for reads by setting `MD_DEV_TEMPLATES_PATH`
func AppTemplateCacheDir() string {
	var templatesPath string
	localDevTemplatesPath := os.Getenv("MD_DEV_TEMPLATES_PATH")
	if localDevTemplatesPath == "" {
		dir, _ := appTemplateCacheDir()
		templatesPath = dir
	} else {
		log.Info().Msgf("Reading templates for local development path: %s", localDevTemplatesPath)
		templatesPath = localDevTemplatesPath
	}

	return templatesPath
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

// appTemplateCacheDir is the actual cache directory. This should be used internally when managing
// files so that development template directories aren't overwritten on accident.
func appTemplateCacheDir() (string, error) {
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
	templateCacheDir, _ := appTemplateCacheDir()
	if err := os.RemoveAll(templateCacheDir); err != nil {
		return err
	}
	return nil
}

func downloadAppTemplates() error {
	templateCacheDir, _ := appTemplateCacheDir()
	log.Debug().Msgf("Downloading templates to %s", templateCacheDir)
	// log.Debug().Msgf("Cloning templates from %s", common.Config().Application.Templates.Repository)
	_, cloneErr := git.PlainClone(templateCacheDir, false, &git.CloneOptions{
		// URL:      common.Config().Application.Templates.Repository,
		URL: MassdriverApplicationTemplatesRepository,
		// Progress: os.Stdout,
		Depth: 1,
	})
	if cloneErr != nil {
		return cloneErr
	}
	return nil
}

package bundle

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"github.com/massdriver-cloud/massdriver-cli/pkg/compress"
	"github.com/rs/zerolog/log"
)

var FileAllows []string = []string{
	"massdriver.yaml",
	"schema-params.json",
	"schema-connections.json",
	"schema-artifacts.json",
	"schema-ui.json",
	"src",
	// do we want this inside a step aka "src" or top-level
	"chart",
}

func PackageBundle(b *Bundle, filePath string, buf io.Writer) error {
	log.Info().Msg("build")
	buildDir := "_build"
	// buildDir, err := os.MkdirTemp("", "bundle-build")
	// if err != nil {
	// 	return err
	// }
	// defer os.RemoveAll(buildDir)

	allowList := []string{}
	allowList = addAllowed(b, allowList)
	log.Info().Msgf("allowed %s", allowList)

	sourceDir := path.Dir(filePath)
	files, errReadDir := ioutil.ReadDir(sourceDir)
	if errReadDir != nil {
		return errReadDir
	}

	for _, fileInfo := range files {
		fileOrDirName := fileInfo.Name()

		// at the top-level directory
		// only copy files or folders in the allowList
		if !common.Contains(allowList, fileOrDirName) {
			continue
		}

		log.Info().Msgf("Copying to build dir: %s", fileOrDirName)
		pathCopyFrom := path.Join(sourceDir, fileOrDirName)
		pathCopyTo := path.Join(buildDir, fileOrDirName)

		if fileInfo.IsDir() {
			log.Info().Msgf("mkdir: %s", fileOrDirName)
			errMkdir := os.Mkdir(filepath.Join(buildDir, fileOrDirName), common.AllRX|common.UserRW)
			if errMkdir != nil {
				return errMkdir
			}

			errCopy := common.CopyFolder(pathCopyFrom, pathCopyTo)
			if errCopy != nil {
				return errCopy
			}
		} else {
			log.Info().Msgf("copying: %s", fileOrDirName)
			var data, err1 = ioutil.ReadFile(pathCopyFrom)
			if err1 != nil {
				return err1
			}

			errWrite := ioutil.WriteFile(pathCopyTo, data, common.AllRWX)
			if errWrite != nil {
				return errWrite
			}
		}
	}

	return tarFolder(buildDir+"/"+common.MassdriverYamlFilename, buf)
}

func addAllowed(b *Bundle, allowList []string) []string {
	allAllows := []string{}
	allAllows = append(allAllows, FileAllows...)
	allAllows = append(allAllows, allowList...)
	log.Info().Msgf("allAllows: %v", allAllows)
	if b.Steps != nil {
		for _, step := range b.Steps {
			allAllows = append(allAllows, step.Path)
		}
	}
	return allAllows
}

func tarFolder(filePath string, buf io.Writer) error {
	// tar > gzip > buf
	gzipWriter := gzip.NewWriter(buf)
	tarWriter := tar.NewWriter(gzipWriter)

	if err := compress.TarDirectory(filepath.Dir(filePath), "bundle", tarWriter); err != nil {
		return err
	}

	// produce tar
	if err := tarWriter.Close(); err != nil {
		return err
	}
	// produce gzip
	if err := gzipWriter.Close(); err != nil {
		return err
	}

	return nil
}

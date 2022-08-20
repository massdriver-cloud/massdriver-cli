package bundle

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"github.com/massdriver-cloud/massdriver-cli/pkg/compress"
	"github.com/rs/zerolog/log"
)

func PackageBundle(b *Bundle, filePath string, buf io.Writer) error {
	// buildDir, err := os.MkdirTemp("", "bundle-build")
	// if err != nil {
	// 	return err
	// }
	// defer os.RemoveAll(buildDir)
	buildDir := "_build"

	allowList := []string{}
	allowList = addAllowed(b)

	sourceDir := path.Dir(filePath)
	files, errReadDir := ioutil.ReadDir(sourceDir)
	if errReadDir != nil {
		return errReadDir
	}

	for _, fileInfo := range files {
		fileOrDirName := fileInfo.Name()

		shouldInclude := false
		for _, allow := range allowList {
			if strings.Contains(fileOrDirName, allow) {
				log.Info().Msgf("including: %s", fileOrDirName)
				shouldInclude = true
			}
		}
		if !shouldInclude {
			continue
		}

		log.Debug().Msgf("Copying to build dir: %s", fileOrDirName)
		pathCopyFrom := path.Join(sourceDir, fileOrDirName)
		pathCopyTo := path.Join(buildDir, fileOrDirName)
		errCopy := common.CopyFolder(pathCopyFrom, pathCopyTo, allowList)
		if errCopy != nil {
			return errCopy
		}
	}

	return tarFolder(buildDir+"/"+common.MassdriverYamlFilename, buf)
}

func addAllowed(b *Bundle) []string {
	allAllows := []string{}
	allAllows = append(allAllows, common.FileAllows...)

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

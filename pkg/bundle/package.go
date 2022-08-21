package bundle

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"github.com/massdriver-cloud/massdriver-cli/pkg/compress"
)

func Package(b *Bundle, filePath string, buf io.Writer) error {
	buildDir, err := os.MkdirTemp("", "bundle-build")
	if err != nil {
		return err
	}
	defer os.RemoveAll(buildDir)

	allowList := getAllowList(b)
	copyConfig := common.CopyConfig{
		Allows:  allowList,
		Ignores: common.FileIgnores,
	}
	srcDir := filepath.Dir(filePath)

	stats, errCopy := common.CopyFolder(srcDir, buildDir, &copyConfig)
	if errCopy != nil {
		return errCopy
	}

	mbs := common.FileSizeMB(stats.FolderSize)
	if mbs > common.MaxBundleSizeMB {
		return fmt.Errorf("Bundle size exceeds maximum allowed size of %vMB", common.MaxBundleSizeMB)
	}

	return tarFolder(buildDir+"/"+common.MassdriverYamlFilename, buf)
}

func getAllowList(b *Bundle) []string {
	allAllows := []string{}
	allAllows = append(allAllows, common.FileAllows...)

	if b.Steps != nil {
		for _, step := range b.Steps {
			allAllows = append(allAllows, step.Path)
		}
	}

	return common.RemoveDuplicateValues(allAllows)
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

package bundle

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"path/filepath"

	"github.com/massdriver-cloud/massdriver-cli/pkg/compress"
)

func PackageBundle(filePath string, buf io.Writer) error {
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

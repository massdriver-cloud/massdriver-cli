package compress_test

import (
	"archive/tar"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/compress"
	"golang.org/x/mod/sumdb/dirhash"
)

func TestTarDirectory(t *testing.T) {
	type test struct {
		name     string
		dirPath  string
		prefix   string
		wantPath string
	}
	tests := []test{
		{
			name:     "simple",
			dirPath:  "testdata/zipdir/",
			prefix:   "wantdir",
			wantPath: "testdata/wantdir",
		},
		{
			name:     "nested",
			dirPath:  "testdata/zipdir/",
			prefix:   "nested/path",
			wantPath: "testdata/nested/path",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testDir := t.TempDir()

			outFile, err := os.Create(path.Join(testDir, "out.tar"))
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}
			defer outFile.Close()

			tarWriter := tar.NewWriter(outFile)
			err = compress.TarDirectory(tc.dirPath, tc.prefix, tarWriter)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			err = tarWriter.Close()
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			cmd := exec.Command("tar", "-xf", outFile.Name(), "-C", testDir)
			err = cmd.Run()
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			wantMD5, err := dirhash.HashDir(tc.wantPath, "", dirhash.DefaultHash)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			gotMD5, err := dirhash.HashDir(path.Join(testDir, tc.prefix), "", dirhash.DefaultHash)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if gotMD5 != wantMD5 {
				t.Errorf("got %v, want %v", gotMD5, wantMD5)
			}
		})
	}
}

func TestTarFile(t *testing.T) {
	type test struct {
		name     string
		filePath string
		prefix   string
		wantPath string
	}
	tests := []test{
		{
			name:     "simple",
			filePath: "testdata/zipdir/main.txt",
			prefix:   "wantdir",
			wantPath: "testdata/wantdir/main.txt",
		},
		{
			name:     "nested",
			filePath: "testdata/zipdir/main.txt",
			prefix:   "nested/path",
			wantPath: "testdata/nested/path/main.txt",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testDir := t.TempDir()

			outFile, err := os.Create(path.Join(testDir, "out.tar"))
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}
			defer outFile.Close()

			tarWriter := tar.NewWriter(outFile)
			err = compress.TarFile(tc.filePath, tc.prefix, tarWriter)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			err = tarWriter.Close()
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			cmd := exec.Command("tar", "-xf", outFile.Name(), "-C", testDir)
			err = cmd.Run()
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			wantContent, err := os.ReadFile(tc.wantPath)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			gotContent, err := os.ReadFile(path.Join(testDir, tc.prefix, filepath.Base(tc.filePath)))
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if string(wantContent) != string(gotContent) {
				t.Errorf("got %v, want %v", string(wantContent), string(gotContent))
			}
		})
	}
}

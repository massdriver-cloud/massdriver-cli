package bundle_test

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
)

func TestPackageBundle(t *testing.T) {
	type test struct {
		name       string
		bundlePath string
		wantPath   string
	}
	tests := []test{
		{
			name:       "simple",
			bundlePath: "testdata/zipdir/massdriver.yaml",
			wantPath:   "testdata/bundle",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got bytes.Buffer
			err := bundle.PackageBundle(tc.bundlePath, &got)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			// Create a temp dir, write out the archive, then shell out to untar
			testDir := t.TempDir()
			zipOut := path.Join(testDir, "out.tar.gz")
			gotBytes := got.Bytes()
			err = os.WriteFile(zipOut, gotBytes, 0644)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}
			cmd := exec.Command("tar", "-xzf", zipOut, "-C", testDir)
			err = cmd.Run()
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}
			compareDirs(t, tc.wantPath, path.Join(testDir, "bundle"))
		})
	}
}

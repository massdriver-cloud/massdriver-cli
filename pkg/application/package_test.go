package application_test

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/application"
	"golang.org/x/mod/sumdb/dirhash"
)

func TestPackage(t *testing.T) {
	type test struct {
		name            string
		applicationPath string
		wantPath        string
	}
	tests := []test{
		{
			name:            "simple",
			applicationPath: "testdata/zipdir/app.yaml",
			wantPath:        "testdata/bundle",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got bytes.Buffer
			err := application.Package(tc.applicationPath, &got)
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

			wantMD5, err := dirhash.HashDir(tc.wantPath, "", dirhash.DefaultHash)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			gotMD5, err := dirhash.HashDir(path.Join(testDir, "bundle"), "", dirhash.DefaultHash)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if gotMD5 != wantMD5 {
				t.Errorf("got %v, want %v", gotMD5, wantMD5)
			}
		})
	}
}

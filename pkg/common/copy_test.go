package common_test

import (
	"path"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"golang.org/x/mod/sumdb/dirhash"
)

func TestCopyFolder(t *testing.T) {
	type test struct {
		name       string
		bundlePath string
		wantPath   string
		config     *common.CopyConfig
	}
	tests := []test{
		{
			name:       "allowlist",
			bundlePath: "testdata/bundle",
			wantPath:   "testdata/bundle-allowlist",
			config: &common.CopyConfig{
				Allows: []string{
					"massdriver.yaml",
					"src",
				},
				Ignores: []string{},
			},
		},
		{
			name:       "allowlist-steps",
			bundlePath: "testdata/bundle-steps",
			wantPath:   "testdata/bundle-steps-allowlist",
			config: &common.CopyConfig{
				Allows: []string{
					"massdriver.yaml",
					"src",
					"core-services",
				},
				Ignores: []string{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testDir := t.TempDir()

			err := common.CopyFolder(tc.bundlePath, testDir, tc.config)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			wantMD5, err := dirhash.HashDir(tc.wantPath, "", dirhash.DefaultHash)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			gotMD5, err := dirhash.HashDir(path.Join(testDir), "", dirhash.DefaultHash)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if gotMD5 != wantMD5 {
				t.Errorf("got %v, want %v", gotMD5, wantMD5)
			}
		})
	}
}

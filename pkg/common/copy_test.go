package common_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"github.com/sergi/go-diff/diffmatchpatch"
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
			name:       "CopyOnlyAllowed",
			bundlePath: "testdata/bundle",
			wantPath:   "testdata/bundle-tar",
			config: &common.CopyConfig{
				Allows: []string{
					"massdriver.yaml",
					"src",
				},
				Ignores: []string{
					".terraform",
				},
			},
		},
		{
			name:       "CopyMultiStep",
			bundlePath: "testdata/bundle-multi-step",
			wantPath:   "testdata/bundle-multi-step-tar",
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

			_, err := common.CopyFolder(tc.bundlePath, testDir, tc.config)
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
				walkAndCompare(tc.wantPath, tc.bundlePath)
			}
		})
	}
}

func TestCopyFolderFromLocal(t *testing.T) {
	copyFrom := "testdata/copy-from"
	copyTo := t.TempDir()
	wantPath := "testdata/copy-to"

	originalDir, err0 := os.Getwd()
	if err0 != nil {
		t.Fatalf("%d, unexpected error", err0)
	}
	err1 := os.Chdir(copyFrom)
	if err1 != nil {
		t.Fatalf("%d, unexpected error", err1)
	}

	config := &common.CopyConfig{
		Allows: []string{
			"main.tf",
		},
		Ignores: []string{},
	}

	_, err2 := common.CopyFolder(".", copyTo, config)
	if err2 != nil {
		t.Fatalf("%d, unexpected error", err2)
	}
	// change back to the original dir so want paths and got paths make sense
	if err := os.Chdir(originalDir); err != nil {
		panic(err)
	}

	wantMD5, err := dirhash.HashDir(wantPath, "", dirhash.DefaultHash)
	if err != nil {
		t.Fatalf("%d, unexpected erroraaaaaaa", err)
	}

	gotMD5, err := dirhash.HashDir(path.Join(copyTo), "", dirhash.DefaultHash)
	if err != nil {
		t.Fatalf("%d, unexpected errorbbbb", err)
	}

	if gotMD5 != wantMD5 {
		t.Errorf("got %v, want %v", gotMD5, wantMD5)
		walkAndCompare(wantPath, ".")
	}
}

func walkAndCompare(wantDir string, gotDir string) {
	_ = gotDir
	err := filepath.Walk(wantDir,
		func(path string, info os.FileInfo, err error) error {
			isDir, _ := isDirectory(path)

			if isDir {
				return nil
			}

			relativeFilePath := strings.TrimPrefix(path, wantDir)
			gotFilePath := filepath.Join(gotDir, relativeFilePath)

			if err != nil {
				return err
			}

			fmt.Printf("Comparing (want) %s and (got) %s\n", path, gotFilePath)

			dmp := diffmatchpatch.New()
			gotText, _ := readFile(gotFilePath)
			wantText, _ := readFile(path)
			diffs := dmp.DiffMain(wantText, gotText, false)

			fmt.Println(dmp.DiffToDelta(diffs))

			return nil
		})
	if err != nil {
		log.Println(err)
	}
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func readFile(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

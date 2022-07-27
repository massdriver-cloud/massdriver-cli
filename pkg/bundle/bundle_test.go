package bundle_test

import (
	"testing"

	"golang.org/x/mod/sumdb/dirhash"
)

func compareDirs(t *testing.T, wantDir string, gotDir string) {
	wantMD5, err := dirhash.HashDir(wantDir, "", dirhash.DefaultHash)
	if err != nil {
		t.Fatalf("%d, unexpected error", err)
	}

	gotMD5, err := dirhash.HashDir(gotDir, "", dirhash.DefaultHash)
	if err != nil {
		t.Fatalf("%d, unexpected error", err)
	}

	if gotMD5 != wantMD5 {
		t.Errorf("got %v, want %v", gotMD5, wantMD5)
	}
}

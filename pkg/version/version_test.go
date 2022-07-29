package version_test

import (
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/version"
)

func TestCheckForNewerVersionAvailable(t *testing.T) {
	tests := []struct {
		name          string
		isOld         bool
		latestVersion string
		wantErr       bool
	}{
		{
			name:          "fail if versions.go is not up to date",
			isOld:         false,
			latestVersion: version.MassVersion(),
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := version.CheckForNewerVersionAvailable()
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckForNewerVersionAvailable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.isOld {
				t.Errorf("CheckForNewerVersionAvailable() got = %v, want %v", got, tt.isOld)
			}
			if got1 != tt.latestVersion {
				t.Errorf("CheckForNewerVersionAvailable() got1 = %v, want %v", got1, tt.latestVersion)
			}
		})
	}
}

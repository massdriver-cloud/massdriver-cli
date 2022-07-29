package version

import (
	"net/http"
	"strings"

	"golang.org/x/mod/semver"
)

const (
	LatestReleaseURL = "https://github.com/massdriver-cloud/massdriver-cli/releases/latest"
)

// var needs to be used instead of const as ldflags is used to fill this
// information in the release process
var (
	version = "v0.1.0"
	gitSHA  = "unknown" // sha1 from git, output of $(git rev-parse HEAD)
)

// MassVersion returns the current version of the github.com/massdriver-cloud/massdriver-cli.
func MassVersion() string {
	return version
}

func MassGitSHA() string {
	return gitSHA
}

func CheckForNewerVersionAvailable() (bool, string, error) {
	resp, err := http.Get(LatestReleaseURL) //nolint:noctx
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()
	// Github will redirect releases/latest to the appropriate releases/tag/vX.X.X
	redirectURL := resp.Request.URL.String()
	parts := strings.Split(redirectURL, "/")
	latestVersion := parts[len(parts)-1]
	if semver.Compare(version, latestVersion) < 0 {
		return true, latestVersion, nil
	}
	return false, latestVersion, nil
}

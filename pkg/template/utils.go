package template

import "strings"

func TypeToName(artifactType string) string {
	noDashes := strings.ReplaceAll(artifactType, "-", "_")
	noSlashes := strings.ReplaceAll(noDashes, "/", "_")
	return noSlashes
}

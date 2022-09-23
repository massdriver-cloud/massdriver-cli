//
package api_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hasura/go-graphql-client"
	"github.com/massdriver-cloud/massdriver-cli/pkg/api"
)

func TestGetPackage(t *testing.T) {
	pkgName := "ecomm-prod-cache"
	mux := http.NewServeMux()
	mux.HandleFunc(APIURL, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data := fmt.Sprintf(`{"data":{"getPackageByNamingConvention":{"manifest": {"id": "manifest-id"}, "target": {"id": "target-id"}, "namePrefix":"%s-8m8q","paramsSchema":"{\"examples\": [{\"__name\": \"Development\",\"name\": \"John Doe\",\"age\": 25}],\"required\": [\"name\"],\"properties\": {\"name\": {\"type\": \"string\"},\"age\": {\"type\": \"integer\",\"default\": 0}}}"}}}`, pkgName)
		mustWrite(w, data)
	})
	client := graphql.NewClient(APIURL, &http.Client{Transport: localRoundTripper{handler: mux}})

	pkg, err := api.GetPackage(client, "faux-org-id", "ecomm-prod-cache")

	if err != nil {
		t.Fatal(err)
	}

	if got, want := pkg.NamePrefix, "ecomm-prod-cache-8m8q"; got != want {
		t.Errorf("got pkg.NamePrefix: %q, want: %q", got, want)
	}
}

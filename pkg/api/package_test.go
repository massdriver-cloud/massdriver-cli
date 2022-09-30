//
package api_test

import (
	"encoding/json"
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

var expectedParams = map[string]interface{}{
	"md_metadata": map[string]interface{}{
		"name_prefix": "local-dev-testbundle-000",
		"default_tags": map[string]string{
			"md-project":  "local",
			"md-target":   "dev",
			"md-manifest": "testbundle",
			"md-package":  "local-dev-testbundle-000",
		},
		"observability": map[string]interface{}{
			"alarm_webhook_url": "https://placeholder.com",
		},
	},
	"cluster_configuration": map[string]interface{}{
		"enable_binary_authorization": false,
	},
	"cluster_networking": map[string]interface{}{
		"cluster_ipv4_cidr_block":  "/16",
		"master_ipv4_cidr_block":   "172.16.0.0/28",
		"services_ipv4_cidr_block": "/20",
	},
	"core_services": map[string]interface{}{
		"cloud_dns_managed_zones": []interface{}{},
		"enable_ingress":          false,
	},
	"k8s_version": "1.21",
	"node_groups": []interface{}{
		map[string]interface{}{
			"name":         "small-pool",
			"machine_type": "e2-highcpu-2",
			"min_size":     1,
			"max_size":     5,
		},
	},
	"observability": map[string]interface{}{
		"logging": map[string]interface{}{
			"destination": "disabled",
		},
	},
}
var mockPkg = &api.Package{
	ID:         "local-dev-testbundle-000",
	Name:       "local-dev-testbundle-000",
	NamePrefix: "local-dev-testbundle-000",
	ProjectID:  "local",
	ManifestID: "testbundle",
	TargetID:   "dev",
}

func TestReadParamsFromFile(t *testing.T) {
	type args struct {
		path string
		pkg  *api.Package
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "json",
			args: args{
				path: "testdata/params.json",
				pkg:  mockPkg,
			},
			want: expectedParams,
		},
		{
			name: "yaml",
			args: args{
				path: "testdata/params.yaml",
				pkg:  mockPkg,
			},
			want: expectedParams,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := api.ReadParamsFromFile(tt.args.path, tt.args.pkg)
			if err != nil {
				t.Errorf("ReadParamsFromFile() error = %v", err)
				return
			}
			gotJSONStr := mustMarshalJSON(t, got)
			want := mustMarshalJSON(t, tt.want)
			if gotJSONStr != want {
				t.Errorf("ReadParamsFromFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mustMarshalJSON(t *testing.T, got map[string]interface{}) string {
	gotJSONStr, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	return string(gotJSONStr)
}

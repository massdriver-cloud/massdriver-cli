package bundle_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
)

func TestLintSchema(t *testing.T) {
	type test struct {
		name string
		bun  *bundle.Bundle
		err  error
	}
	tests := []test{
		{
			name: "Valid pass",
			bun: &bundle.Bundle{
				Name:        "example",
				Description: "description",
				Access:      "private",
				Schema:      "draft-07",
				Type:        "infrastructure",
				Params:      map[string]interface{}{},
				Connections: map[string]interface{}{},
				Artifacts:   map[string]interface{}{},
				UI:          map[string]interface{}{},
			},
			err: nil,
		},
		{
			name: "Invalid missing schema",
			bun: &bundle.Bundle{
				Name:        "example",
				Description: "description",
				Access:      "private",
				Type:        "infrastructure",
				Params:      map[string]interface{}{},
				Connections: map[string]interface{}{},
				Artifacts:   map[string]interface{}{},
				UI:          map[string]interface{}{},
			},
			err: errors.New(`massdriver.yaml has schema violations:
	- schema: schema must be one of the following: "draft-07"
`),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := bundle.LintSchema(tc.bun)
			if tc.err != nil {
				if err == nil {
					t.Errorf("expected an error, got nil")
				} else if tc.err.Error() != err.Error() {
					t.Errorf("got %v, want %v", err.Error(), tc.err.Error())
				}
			} else if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}
		})
	}
}

func TestLintParamsConnectionsNameCollision(t *testing.T) {
	type test struct {
		name string
		bun  *bundle.Bundle
		err  error
	}
	tests := []test{
		{
			name: "Valid Pass",
			bun: &bundle.Bundle{
				Params: map[string]interface{}{
					"properties": map[string]interface{}{
						"param": "foo",
					},
				},
				Connections: map[string]interface{}{
					"properties": map[string]interface{}{
						"connection": "foo",
					},
				},
			},
			err: nil,
		},
		{
			name: "Invalid Error",
			bun: &bundle.Bundle{
				Params: map[string]interface{}{
					"properties": map[string]interface{}{
						"database": "foo",
					},
				},
				Connections: map[string]interface{}{
					"properties": map[string]interface{}{
						"database": "foo",
					},
				},
			},
			err: fmt.Errorf("a parameter and connection have the same name: database"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := bundle.LintParamsConnectionsNameCollision(tc.bun)
			if tc.err != nil {
				if err == nil {
					t.Errorf("expected an error, got nil")
				} else if tc.err.Error() != err.Error() {
					t.Errorf("got %v, want %v", err.Error(), tc.err.Error())
				}
			} else if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}
		})
	}
}

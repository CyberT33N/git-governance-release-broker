package packaging

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDeploymentWorkflowContracts(t *testing.T) {
	root := filepath.Join("..", "..")
	for _, testCase := range []struct {
		name      string
		path      string
		required  []string
		forbidden []string
	}{
		{
			name: "staging",
			path: filepath.Join(".github", "workflows", "gcp-broker-staging.yml"),
			required: []string{
				"environment: gcp-broker-staging",
				"test \"$GITHUB_REF\" = \"refs/heads/develop\"",
				"BROKER_CREDENTIAL_PROFILE=release-automation",
			},
		},
		{
			name: "production",
			path: filepath.Join(".github", "workflows", "gcp-broker-production.yml"),
			required: []string{
				"environment: gcp-broker-production",
				"test \"$GITHUB_REF\" = \"refs/heads/main\"",
				"@sha256:",
				"BROKER_CREDENTIAL_PROFILE=release-automation",
			},
			forbidden: []string{"docker build"},
		},
		{
			name: "reconciliation publisher",
			path: filepath.Join(".github", "workflows", "gcp-reconciliation-publisher-production.yml"),
			required: []string{
				"environment: gcp-reconciliation-publisher-deployment",
				"test \"$GITHUB_REF\" = \"refs/heads/main\"",
				"@sha256:",
				"BROKER_CREDENTIAL_PROFILE=reconciliation-publisher",
			},
			forbidden: []string{"docker build"},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			contents, err := os.ReadFile(filepath.Join(root, testCase.path))
			if err != nil {
				t.Fatalf("ReadFile(%q) error = %v", testCase.path, err)
			}
			workflow := string(contents)
			for _, value := range testCase.required {
				if !strings.Contains(workflow, value) {
					t.Fatalf("workflow missing %q", value)
				}
			}
			for _, value := range testCase.forbidden {
				if strings.Contains(workflow, value) {
					t.Fatalf("workflow contains forbidden %q", value)
				}
			}
		})
	}

	if _, err := os.Stat(filepath.Join(root, ".github", "workflows", "gcp-deploy.yml")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("retired develop-bound deployment workflow error = %v", err)
	}
}

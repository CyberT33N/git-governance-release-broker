package packaging

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWorkflowContracts(t *testing.T) {
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
		{
			name: "production artifact promotion",
			path: filepath.Join(".github", "workflows", "gcp-broker-production-promotion.yml"),
			required: []string{
				"environment: gcp-broker-production",
				"GCP_PRODUCTION_ARTIFACT_PROMOTION_WIF_PROVIDER",
				"GCP_PRODUCTION_ARTIFACT_PROMOTER_SERVICE_ACCOUNT",
				"GCP_PRODUCTION_SOURCE_ARTIFACT_REPOSITORY",
				"test \"$GITHUB_REF\" = \"refs/heads/main\"",
				"git merge-base --is-ancestor \"$SOURCE_COMMIT\" HEAD",
				"docker pull \"$source_image\"",
				"docker tag \"$source_image\" \"$target_tag\"",
				"docker push \"$target_tag\"",
				"test \"$target_digest\" = \"$source_digest\"",
			},
			forbidden: []string{
				"docker build",
				"gcloud run deploy",
				"latest",
			},
		},
		{
			name: "reconciliation publisher artifact promotion",
			path: filepath.Join(".github", "workflows", "gcp-reconciliation-publisher-promotion.yml"),
			required: []string{
				"environment: gcp-reconciliation-publisher-deployment",
				"GCP_RECONCILIATION_PUBLISHER_ARTIFACT_PROMOTION_WIF_PROVIDER",
				"GCP_RECONCILIATION_PUBLISHER_ARTIFACT_PROMOTER_SERVICE_ACCOUNT",
				"GCP_RECONCILIATION_PUBLISHER_SOURCE_ARTIFACT_REPOSITORY",
				"test \"$GITHUB_REF\" = \"refs/heads/main\"",
				"git merge-base --is-ancestor \"$SOURCE_COMMIT\" HEAD",
				"docker pull \"$source_image\"",
				"docker tag \"$source_image\" \"$target_tag\"",
				"docker push \"$target_tag\"",
				"test \"$target_digest\" = \"$source_digest\"",
			},
			forbidden: []string{
				"docker build",
				"gcloud run deploy",
				"latest",
			},
		},
		{
			name: "protected shared line",
			path: filepath.Join(".github", "workflows", "create-protected-line.yml"),
			required: []string{
				"github.repository == 'CyberT33N/git-governance-release-broker'",
				"github.ref == 'refs/heads/main'",
				"environment: release",
				"request_id:",
				"source=\"origin/develop\"",
				"git push origin \"${SOURCE}:refs/heads/${TARGET}\"",
			},
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

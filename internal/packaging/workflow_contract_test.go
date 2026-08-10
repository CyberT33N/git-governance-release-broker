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
				"attestations: write",
				"GCP_STAGING_EVIDENCE_ARTIFACT_REPOSITORY",
				"anchore/sbom-action@e22c389904149dbc22b58101806040fa8d37a610",
				"actions/attest@f7c74d28b9d84cb8768d0b8ca14a4bac6ef463e6",
				"sigstore/cosign-installer@6f9f17788090df1f26f669e9d70d6ae9567deba6",
				"cosign-release: v3.1.3",
				"cosign sign --yes \"$IMAGE\"",
				"gcloud artifacts generic upload",
			},
			forbidden: []string{
				"COSIGN_EXPERIMENTAL",
				"--registry-referrers-mode",
				"--experimental-oci11",
				"sigstore/cosign-installer@4959ce089c160fddf62f7b42464195ba1a56d382",
			},
		},
		{
			name: "container evidence verifier",
			path: filepath.Join(".github", "actions", "verify-broker-evidence", "action.yml"),
			required: []string{
				"sigstore/cosign-installer@6f9f17788090df1f26f669e9d70d6ae9567deba6",
				"cosign-release: v3.1.3",
				"gcloud artifacts generic download",
				"cosign verify",
				"gh attestation verify",
				"https://spdx.dev/Document/v2.3",
				"--bundle-from-oci",
				"--deny-self-hosted-runners",
				"--source-digest",
				"provenance_bundle_sha256",
				"sbom_bundle_sha256",
			},
			forbidden: []string{
				"cosign sign",
				"docker build",
				"gcloud run deploy",
				"COSIGN_EXPERIMENTAL",
				"--registry-referrers-mode",
				"--experimental-oci11",
				"sigstore/cosign-installer@4959ce089c160fddf62f7b42464195ba1a56d382",
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
				"attestations: read",
				"GCP_PRODUCTION_EVIDENCE_ARTIFACT_REPOSITORY",
				"./.github/actions/verify-broker-evidence",
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
				"GCP_PRODUCTION_SOURCE_EVIDENCE_ARTIFACT_REPOSITORY",
				"GCP_PRODUCTION_EVIDENCE_ARTIFACT_REPOSITORY",
				"attestations: read",
				"git merge-base --is-ancestor \"$SOURCE_COMMIT\" HEAD",
				"./.github/actions/verify-broker-evidence",
				"gcloud artifacts generic upload",
			},
			forbidden: []string{
				"docker build",
				"gcloud run deploy",
				"latest",
			},
		},
		{
			name: "release credential verification",
			path: filepath.Join(".github", "workflows", "gcp-release-credential-verification-production.yml"),
			required: []string{
				"environment: gcp-release-credential-verification-deployment",
				"GCP_RELEASE_CREDENTIAL_VERIFICATION_DEPLOYMENT_WORKLOAD_IDENTITY_PROVIDER",
				"GCP_RELEASE_CREDENTIAL_VERIFICATION_DEPLOYMENT_SERVICE_ACCOUNT",
				"test \"$GITHUB_REF\" = \"refs/heads/main\"",
				"@sha256:",
				"BROKER_CREDENTIAL_PROFILE=release-credential-verification",
				"attestations: read",
				"GCP_RELEASE_CREDENTIAL_VERIFICATION_EVIDENCE_ARTIFACT_REPOSITORY",
				"./.github/actions/verify-broker-evidence",
			},
			forbidden: []string{"docker build"},
		},
		{
			name: "release credential verification artifact promotion",
			path: filepath.Join(".github", "workflows", "gcp-release-credential-verification-promotion.yml"),
			required: []string{
				"environment: gcp-release-credential-verification-deployment",
				"GCP_RELEASE_CREDENTIAL_VERIFICATION_ARTIFACT_PROMOTION_WIF_PROVIDER",
				"GCP_RELEASE_CREDENTIAL_VERIFICATION_ARTIFACT_PROMOTER_SERVICE_ACCOUNT",
				"GCP_RELEASE_CREDENTIAL_VERIFICATION_SOURCE_ARTIFACT_REPOSITORY",
				"GCP_RELEASE_CREDENTIAL_VERIFICATION_SOURCE_EVIDENCE_ARTIFACT_REPOSITORY",
				"GCP_RELEASE_CREDENTIAL_VERIFICATION_EVIDENCE_ARTIFACT_REPOSITORY",
				"attestations: read",
				"test \"$GITHUB_REF\" = \"refs/heads/main\"",
				"git merge-base --is-ancestor \"$SOURCE_COMMIT\" HEAD",
				"./.github/actions/verify-broker-evidence",
				"docker pull \"$SOURCE_IMAGE\"",
				"docker tag \"$SOURCE_IMAGE\" \"$target_tag\"",
				"docker push \"$target_tag\"",
				"test \"$target_digest\" = \"$SOURCE_DIGEST\"",
				"gcloud artifacts generic upload",
			},
			forbidden: []string{
				"docker build",
				"gcloud run deploy",
				"latest",
			},
		},
		{
			name: "hotfix delivery",
			path: filepath.Join(".github", "workflows", "gcp-hotfix-delivery-production.yml"),
			required: []string{
				"environment: gcp-hotfix-delivery-deployment",
				"GCP_HOTFIX_DELIVERY_DEPLOYMENT_WORKLOAD_IDENTITY_PROVIDER",
				"GCP_HOTFIX_DELIVERY_DEPLOYMENT_SERVICE_ACCOUNT",
				"test \"$GITHUB_REF\" = \"refs/heads/main\"",
				"@sha256:",
				"BROKER_CREDENTIAL_PROFILE=hotfix-delivery",
				"attestations: read",
				"GCP_HOTFIX_DELIVERY_EVIDENCE_ARTIFACT_REPOSITORY",
				"./.github/actions/verify-broker-evidence",
			},
			forbidden: []string{"docker build"},
		},
		{
			name: "hotfix delivery artifact promotion",
			path: filepath.Join(".github", "workflows", "gcp-hotfix-delivery-promotion.yml"),
			required: []string{
				"environment: gcp-hotfix-delivery-deployment",
				"GCP_HOTFIX_DELIVERY_ARTIFACT_PROMOTION_WIF_PROVIDER",
				"GCP_HOTFIX_DELIVERY_ARTIFACT_PROMOTER_SERVICE_ACCOUNT",
				"GCP_HOTFIX_DELIVERY_SOURCE_ARTIFACT_REPOSITORY",
				"GCP_HOTFIX_DELIVERY_SOURCE_EVIDENCE_ARTIFACT_REPOSITORY",
				"GCP_HOTFIX_DELIVERY_EVIDENCE_ARTIFACT_REPOSITORY",
				"attestations: read",
				"test \"$GITHUB_REF\" = \"refs/heads/main\"",
				"git merge-base --is-ancestor \"$SOURCE_COMMIT\" HEAD",
				"./.github/actions/verify-broker-evidence",
				"docker pull \"$SOURCE_IMAGE\"",
				"docker tag \"$SOURCE_IMAGE\" \"$target_tag\"",
				"docker push \"$target_tag\"",
				"test \"$target_digest\" = \"$SOURCE_DIGEST\"",
				"gcloud artifacts generic upload",
			},
			forbidden: []string{
				"docker build",
				"gcloud run deploy",
				"latest",
			},
		},
		{
			name: "reconciliation publisher",
			path: filepath.Join(".github", "workflows", "gcp-reconciliation-publisher-production.yml"),
			required: []string{
				"environment: gcp-reconciliation-publisher-deployment",
				"test \"$GITHUB_REF\" = \"refs/heads/main\"",
				"@sha256:",
				"BROKER_CREDENTIAL_PROFILE=reconciliation-publisher",
				"attestations: read",
				"GCP_RECONCILIATION_PUBLISHER_EVIDENCE_ARTIFACT_REPOSITORY",
				"./.github/actions/verify-broker-evidence",
			},
			forbidden: []string{"docker build"},
		},
		{
			name: "reconciliation publisher artifact promotion",
			path: filepath.Join(".github", "workflows", "gcp-reconciliation-publisher-promotion.yml"),
			required: []string{
				"environment: gcp-reconciliation-publisher-deployment",
				"GCP_RECONCILIATION_PUBLISHER_ARTIFACT_PROMOTION_WIF_PROVIDER",
				"GCP_RECONCILIATION_PUBLISHER_ARTIFACT_PROMOTER_SERVICE_ACCOUNT",
				"GCP_RECONCILIATION_PUBLISHER_SOURCE_ARTIFACT_REPOSITORY",
				"GCP_RECONCILIATION_PUBLISHER_SOURCE_EVIDENCE_ARTIFACT_REPOSITORY",
				"GCP_RECONCILIATION_PUBLISHER_EVIDENCE_ARTIFACT_REPOSITORY",
				"attestations: read",
				"git merge-base --is-ancestor \"$SOURCE_COMMIT\" HEAD",
				"./.github/actions/verify-broker-evidence",
				"gcloud artifacts generic upload",
			},
			forbidden: []string{
				"docker build",
				"gcloud run deploy",
				"latest",
			},
		},
		{
			name: "hotfix propagation publisher",
			path: filepath.Join(".github", "workflows", "gcp-hotfix-propagation-publisher-production.yml"),
			required: []string{
				"environment: gcp-hotfix-propagation-publisher-deployment",
				"test \"$GITHUB_REF\" = \"refs/heads/main\"",
				"@sha256:",
				"BROKER_CREDENTIAL_PROFILE=hotfix-propagation-publisher",
				"attestations: read",
				"GCP_HOTFIX_PROPAGATION_PUBLISHER_EVIDENCE_ARTIFACT_REPOSITORY",
				"./.github/actions/verify-broker-evidence",
			},
			forbidden: []string{"docker build"},
		},
		{
			name: "hotfix propagation publisher artifact promotion",
			path: filepath.Join(".github", "workflows", "gcp-hotfix-propagation-publisher-promotion.yml"),
			required: []string{
				"environment: gcp-hotfix-propagation-publisher-deployment",
				"GCP_HOTFIX_PROPAGATION_PUBLISHER_ARTIFACT_PROMOTION_WIF_PROVIDER",
				"GCP_HOTFIX_PROPAGATION_PUBLISHER_ARTIFACT_PROMOTER_SERVICE_ACCOUNT",
				"GCP_HOTFIX_PROPAGATION_PUBLISHER_SOURCE_ARTIFACT_REPOSITORY",
				"GCP_HOTFIX_PROPAGATION_PUBLISHER_SOURCE_EVIDENCE_ARTIFACT_REPOSITORY",
				"GCP_HOTFIX_PROPAGATION_PUBLISHER_EVIDENCE_ARTIFACT_REPOSITORY",
				"attestations: read",
				"test \"$GITHUB_REF\" = \"refs/heads/main\"",
				"git merge-base --is-ancestor \"$SOURCE_COMMIT\" HEAD",
				"./.github/actions/verify-broker-evidence",
				"docker pull \"$SOURCE_IMAGE\"",
				"docker tag \"$SOURCE_IMAGE\" \"$target_tag\"",
				"docker push \"$target_tag\"",
				"test \"$target_digest\" = \"$SOURCE_DIGEST\"",
				"gcloud artifacts generic upload",
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

func TestStagingWorkflowPreparesEvidenceWorkspaceBeforeSBOM(t *testing.T) {
	workflowPath := filepath.Join("..", "..", ".github", "workflows", "gcp-broker-staging.yml")
	contents, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", workflowPath, err)
	}

	workflow := string(contents)
	workspaceIndex := strings.Index(workflow, "- name: Prepare staging evidence workspace")
	if workspaceIndex < 0 {
		t.Fatal("staging workflow is missing the evidence workspace step")
	}
	if !strings.Contains(workflow[workspaceIndex:], "mkdir -p .build/evidence") {
		t.Fatal("staging evidence workspace step does not create .build/evidence")
	}

	sbomIndex := strings.Index(workflow, "- name: Generate staging image SBOM")
	if sbomIndex < 0 {
		t.Fatal("staging workflow is missing the SBOM generation step")
	}
	if workspaceIndex > sbomIndex {
		t.Fatal("staging evidence workspace is prepared after SBOM generation")
	}
}

func TestCosignV3UsesBundleDefaults(t *testing.T) {
	stagingPath := filepath.Join("..", "..", ".github", "workflows", "gcp-broker-staging.yml")
	stagingContents, err := os.ReadFile(stagingPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", stagingPath, err)
	}

	stagingWorkflow := string(stagingContents)
	signatureStepStart := strings.Index(stagingWorkflow, "- name: Sign and verify immutable staging image")
	if signatureStepStart < 0 {
		t.Fatal("staging workflow is missing the Cosign signature step")
	}
	signatureStepEnd := strings.Index(stagingWorkflow[signatureStepStart:], "- id: provenance")
	if signatureStepEnd < 0 {
		t.Fatal("staging workflow is missing the provenance step after Cosign")
	}
	signatureStep := stagingWorkflow[signatureStepStart : signatureStepStart+signatureStepEnd]
	if !strings.Contains(stagingWorkflow, "sigstore/cosign-installer@6f9f17788090df1f26f669e9d70d6ae9567deba6") {
		t.Fatal("staging workflow does not pin the Cosign v3-compatible installer")
	}
	if !strings.Contains(stagingWorkflow, "cosign-release: v3.1.3") {
		t.Fatal("staging workflow does not pin Cosign v3.1.3")
	}
	if !strings.Contains(signatureStep, "cosign sign --yes \"$IMAGE\"") {
		t.Fatal("staging workflow does not use the Cosign v3 bundle signing default")
	}
	if strings.Contains(stagingWorkflow, "COSIGN_EXPERIMENTAL") {
		t.Fatal("staging workflow retains experimental Cosign mode")
	}
	if strings.Contains(stagingWorkflow, "--registry-referrers-mode") {
		t.Fatal("staging workflow retains deprecated OCI referrer mode")
	}
	stagingVerifyStart := strings.Index(signatureStep, "cosign verify \\")
	if stagingVerifyStart < 0 {
		t.Fatal("staging Cosign signature step is missing verification")
	}
	stagingVerifyEnd := strings.Index(signatureStep[stagingVerifyStart:], "$IMAGE\" >/dev/null")
	if stagingVerifyEnd < 0 {
		t.Fatal("staging Cosign verification is missing its image target")
	}
	stagingVerify := signatureStep[stagingVerifyStart : stagingVerifyStart+stagingVerifyEnd]
	if strings.Contains(stagingVerify, "--registry-referrers-mode") {
		t.Fatal("staging Cosign verification uses unsupported OCI referrers mode")
	}
	if strings.Contains(stagingVerify, "--experimental-oci11") {
		t.Fatal("staging Cosign verification retains removed OCI 1.1 discovery mode")
	}

	verifierPath := filepath.Join("..", "..", ".github", "actions", "verify-broker-evidence", "action.yml")
	verifierContents, err := os.ReadFile(verifierPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", verifierPath, err)
	}

	verifier := string(verifierContents)
	verificationStepStart := strings.Index(verifier, "- id: verify")
	if verificationStepStart < 0 {
		t.Fatal("evidence verifier is missing the Cosign verification step")
	}
	if !strings.Contains(verifier, "sigstore/cosign-installer@6f9f17788090df1f26f669e9d70d6ae9567deba6") {
		t.Fatal("evidence verifier does not pin the Cosign v3-compatible installer")
	}
	if !strings.Contains(verifier, "cosign-release: v3.1.3") {
		t.Fatal("evidence verifier does not pin Cosign v3.1.3")
	}
	if strings.Contains(verifier, "COSIGN_EXPERIMENTAL") {
		t.Fatal("evidence verifier retains experimental Cosign mode")
	}
	if strings.Contains(verifier, "--registry-referrers-mode") {
		t.Fatal("evidence verifier retains deprecated OCI referrer mode")
	}
	verifierVerifyStart := strings.Index(verifier[verificationStepStart:], "cosign verify \\")
	if verifierVerifyStart < 0 {
		t.Fatal("evidence verifier is missing Cosign verification")
	}
	verifierVerifyStart += verificationStepStart
	verifierVerifyEnd := strings.Index(verifier[verifierVerifyStart:], "$source_image\" >/dev/null")
	if verifierVerifyEnd < 0 {
		t.Fatal("evidence verifier is missing its image target")
	}
	verifierVerify := verifier[verifierVerifyStart : verifierVerifyStart+verifierVerifyEnd]
	if strings.Contains(verifierVerify, "--registry-referrers-mode") {
		t.Fatal("evidence verifier uses unsupported OCI referrers mode")
	}
	if strings.Contains(verifierVerify, "--experimental-oci11") {
		t.Fatal("evidence verifier retains removed OCI 1.1 discovery mode")
	}
}

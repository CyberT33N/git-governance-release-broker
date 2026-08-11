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
				"token_format: access_token",
				"docker/login-action@dbcb813823bdd20940b903addbd779551569679f",
				"username: oauth2accesstoken",
				"password: ${{ steps.auth.outputs.access_token }}",
				"logout: true",
				"gcloud artifacts generic upload",
			},
			forbidden: []string{
				"COSIGN_EXPERIMENTAL",
				"--registry-referrers-mode",
				"--experimental-oci11",
				"sigstore/cosign-installer@4959ce089c160fddf62f7b42464195ba1a56d382",
				"gcloud auth configure-docker",
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
				"evidence_type == \"provenance\"",
				"evidence_type == \"attestation\"",
				"evidence_type == \"signature\"",
				"require-promotion:",
				"gcloud auth print-access-token",
				"docker login \"$registry\" --username oauth2accesstoken --password-stdin",
				"docker logout \"$registry\"",
			},
			forbidden: []string{
				"cosign sign",
				"docker build",
				"gcloud run deploy",
				"COSIGN_EXPERIMENTAL",
				"--registry-referrers-mode",
				"--experimental-oci11",
				"sigstore/cosign-installer@4959ce089c160fddf62f7b42464195ba1a56d382",
				"gcloud auth configure-docker",
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
				"./.github/actions/record-broker-promotion-evidence",
			},
			forbidden: []string{
				"docker build",
				"gcloud run deploy",
				"latest",
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
				"./.github/actions/record-broker-promotion-evidence",
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
				"./.github/actions/record-broker-promotion-evidence",
			},
			forbidden: []string{
				"docker build",
				"gcloud run deploy",
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
				"./.github/actions/record-broker-promotion-evidence",
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
				"./.github/actions/record-broker-promotion-evidence",
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

func TestAttestationRegistryAuthUsesShortLivedStaticDockerCredentials(t *testing.T) {
	stagingPath := filepath.Join("..", "..", ".github", "workflows", "gcp-broker-staging.yml")
	stagingContents, err := os.ReadFile(stagingPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", stagingPath, err)
	}

	stagingWorkflow := strings.ReplaceAll(string(stagingContents), "\r\n", "\n")
	bootstrapIndex := strings.Index(stagingWorkflow, "\n  bootstrap:\n")
	if bootstrapIndex < 0 {
		t.Fatal("staging workflow is missing the bootstrap job")
	}
	deployIndex := strings.Index(stagingWorkflow, "\n  deploy:\n")
	if deployIndex < 0 {
		t.Fatal("staging workflow is missing the deploy job")
	}
	if bootstrapIndex > deployIndex {
		t.Fatal("staging workflow declares the deploy job before bootstrap")
	}

	bootstrapJob := stagingWorkflow[bootstrapIndex:deployIndex]
	if strings.Contains(bootstrapJob, "- id: auth") {
		t.Fatal("bootstrap job exposes a deploy-only authentication output")
	}
	if strings.Contains(bootstrapJob, "token_format: access_token") {
		t.Fatal("bootstrap job requests a deploy-only access token")
	}

	deployJob := stagingWorkflow[deployIndex:]
	authIndex := strings.Index(deployJob, "- id: auth")
	if authIndex < 0 {
		t.Fatal("deploy job is missing the access-token authentication step")
	}
	loginIndex := strings.Index(deployJob, "- name: Authenticate Docker client to Artifact Registry")
	if loginIndex < 0 {
		t.Fatal("deploy job is missing static Docker registry authentication")
	}
	buildIndex := strings.Index(deployJob, "- name: Build immutable staging image")
	if buildIndex < 0 {
		t.Fatal("deploy job is missing the image build step")
	}
	if authIndex > loginIndex || loginIndex > buildIndex {
		t.Fatal("deploy registry authentication is not ordered before the image build")
	}
	for _, required := range []string{
		"token_format: access_token",
		"docker/login-action@dbcb813823bdd20940b903addbd779551569679f",
		"username: oauth2accesstoken",
		"password: ${{ steps.auth.outputs.access_token }}",
		"logout: true",
	} {
		if !strings.Contains(deployJob, required) {
			t.Fatalf("deploy job is missing %q", required)
		}
	}
	if strings.Contains(stagingWorkflow, "gcloud auth configure-docker") {
		t.Fatal("staging workflow retains an unsupported Docker credential helper")
	}

	verifierPath := filepath.Join("..", "..", ".github", "actions", "verify-broker-evidence", "action.yml")
	verifierContents, err := os.ReadFile(verifierPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", verifierPath, err)
	}

	verifier := string(verifierContents)
	for _, required := range []string{
		"docker_config=\"$(mktemp -d)\"",
		"export DOCKER_CONFIG=\"$docker_config\"",
		"gcloud auth print-access-token",
		"docker login \"$registry\" --username oauth2accesstoken --password-stdin",
		"unset access_token",
		"docker logout \"$registry\"",
		"rm -rf \"$docker_config\" \"$evidence_directory\"",
	} {
		if !strings.Contains(verifier, required) {
			t.Fatalf("evidence verifier is missing %q", required)
		}
	}
	if strings.Contains(verifier, "gcloud auth configure-docker") {
		t.Fatal("evidence verifier retains an unsupported Docker credential helper")
	}
}

func TestEvidenceVerifierCallersCreateGoogleCredentialsFiles(t *testing.T) {
	root := filepath.Join("..", "..")
	for _, path := range []string{
		filepath.Join(".github", "workflows", "gcp-broker-production.yml"),
		filepath.Join(".github", "workflows", "gcp-broker-production-promotion.yml"),
		filepath.Join(".github", "workflows", "gcp-reconciliation-publisher-production.yml"),
		filepath.Join(".github", "workflows", "gcp-reconciliation-publisher-promotion.yml"),
		filepath.Join(".github", "workflows", "gcp-release-credential-verification-production.yml"),
		filepath.Join(".github", "workflows", "gcp-release-credential-verification-promotion.yml"),
		filepath.Join(".github", "workflows", "gcp-hotfix-delivery-production.yml"),
		filepath.Join(".github", "workflows", "gcp-hotfix-delivery-promotion.yml"),
		filepath.Join(".github", "workflows", "gcp-hotfix-propagation-publisher-production.yml"),
		filepath.Join(".github", "workflows", "gcp-hotfix-propagation-publisher-promotion.yml"),
	} {
		contents, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			t.Fatalf("ReadFile(%q) error = %v", path, err)
		}

		workflow := string(contents)
		if !strings.Contains(workflow, "uses: ./.github/actions/verify-broker-evidence") {
			t.Fatalf("workflow %q does not call the evidence verifier", path)
		}
		if !strings.Contains(workflow, "create_credentials_file: true") {
			t.Fatalf("workflow %q does not provide Google credentials to the evidence verifier", path)
		}
	}
}

func TestStagingWorkflowMaterializesEGP1Subjects(t *testing.T) {
	stagingPath := filepath.Join("..", "..", ".github", "workflows", "gcp-broker-staging.yml")
	contents, err := os.ReadFile(stagingPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", stagingPath, err)
	}

	workflow := string(contents)
	for _, required := range []string{
		"Materialize local builder and isolated dependency evidence",
		"--network=none",
		"docker pull --platform=linux/amd64 \"$builder_image\"",
		"local_builder_image=\"broker-local-builder:${builder_digest#sha256:}\"",
		"docker tag \"$builder_image\" \"$local_builder_image\"",
		"docker image inspect \"$local_builder_image\"",
		"--pull=false --network=none",
		"dependency-resolution.json",
		"test-result.json",
		"cosign download signature \"$IMAGE\"",
		"signature.json",
		"evidence-graph/v1",
		"source.subject.payload.json",
		"dependency-resolution.subject.payload.json",
		"build.subject.payload.json",
		"artifact.subject.payload.json",
		"for subject_type in promotion operation",
		"${subject_type}.not-recorded.subject.payload.json",
		"Seal immutable staging evidence subjects",
		"./.github/actions/seal-evidence-subject",
		"Record immutable staging deployment subject",
		"./.github/actions/record-broker-deployment-evidence",
		"lane: staging",
		"workflow: .github/workflows/gcp-broker-staging.yml",
		"source-ref: refs/heads/develop",
		"environment: gcp-broker-staging",
		"source_tree=\"$(git rev-parse \"${GITHUB_SHA}^{tree}\")\"",
		"go_mod_sha256",
		"dockerfile_sha256",
		"go_toolchain",
		"dependency_resolution:",
		"admission: {",
		"phase_status_at_issuance:",
	} {
		if !strings.Contains(workflow, required) {
			t.Fatalf("staging EGP-1 workflow is missing %q", required)
		}
	}
	for _, forbidden := range []string{
		"schema_version: 2",
		"manifest.json",
		"external-prerequisite-unverified",
		"not-evaluated",
	} {
		if strings.Contains(workflow, forbidden) {
			t.Fatalf("staging EGP-1 workflow retains obsolete %q", forbidden)
		}
	}

	dependencyIndex := strings.Index(workflow, "- name: Materialize local builder and isolated dependency evidence")
	buildIndex := strings.Index(workflow, "- name: Build immutable staging image")
	pushIndex := strings.Index(workflow, "- name: Push immutable staging image")
	sbomIndex := strings.Index(workflow, "- name: Generate staging image SBOM")
	if dependencyIndex < 0 || buildIndex < 0 || pushIndex < 0 || sbomIndex < 0 ||
		dependencyIndex > buildIndex || buildIndex > pushIndex || pushIndex > sbomIndex {
		t.Fatal("staging builder and dependency evidence is not materialized before the isolated local build and image publication")
	}
	deployIndex := strings.Index(workflow, "gcloud run deploy")
	recordIndex := strings.Index(workflow, "uses: ./.github/actions/record-broker-deployment-evidence")
	if deployIndex < 0 || recordIndex < 0 || deployIndex > recordIndex {
		t.Fatal("staging deployment evidence is not recorded after the Cloud Run mutation")
	}
}

func TestStagingWorkflowUsesMaterializedLocalBuilderForOfflineConsumer(t *testing.T) {
	stagingPath := filepath.Join("..", "..", ".github", "workflows", "gcp-broker-staging.yml")
	contents, err := os.ReadFile(stagingPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", stagingPath, err)
	}

	workflow := strings.ReplaceAll(string(contents), "\r\n", "\n")
	materializationStart := strings.Index(workflow, "- name: Materialize local builder and isolated dependency evidence")
	buildStart := strings.Index(workflow, "- name: Build immutable staging image")
	if materializationStart < 0 || buildStart < 0 || materializationStart > buildStart {
		t.Fatal("staging workflow does not materialize the local builder before the final image build")
	}

	materialization := workflow[materializationStart:buildStart]
	pullIndex := strings.Index(materialization, "docker pull --platform=linux/amd64 \"$builder_image\"")
	tagIndex := strings.Index(materialization, "docker tag \"$builder_image\" \"$local_builder_image\"")
	firstOfflineRun := strings.Index(materialization, "docker run --rm")
	if pullIndex < 0 || tagIndex < 0 || firstOfflineRun < 0 || pullIndex > tagIndex || tagIndex > firstOfflineRun {
		t.Fatal("staging workflow does not materialize and tag the builder before offline consumption")
	}

	for _, required := range []string{
		"builder_digest=\"${builder_image##*@}\"",
		"[[ \"$builder_digest\" =~ ^sha256:[0-9a-f]{64}$ ]]",
		"local_builder_digests=\"$(docker image inspect \"$local_builder_image\" --format '{{range .RepoDigests}}{{println .}}{{end}}')\"",
		"awk -F@ -v expected=\"$builder_digest\" '$NF == expected { found = 1 } END { exit !found }'",
		"test \"$(docker image inspect \"$local_builder_image\" --format '{{.Os}}/{{.Architecture}}')\" = \"linux/amd64\"",
		"\"$local_builder_image\" list -m -json all",
		"\"$local_builder_image\" mod verify",
		"\"$local_builder_image\" test -mod=readonly ./...",
		"\"$local_builder_image\" env GOVERSION",
		"--pull=never",
		"--network=none",
	} {
		if !strings.Contains(materialization, required) {
			t.Fatalf("local builder materialization is missing %q", required)
		}
	}
	for _, forbidden := range []string{
		"\"$builder_image\" list -m -json all",
		"\"$builder_image\" mod verify",
		"\"$builder_image\" test -mod=readonly ./...",
		"\"$builder_image\" env GOVERSION",
	} {
		if strings.Contains(materialization, forbidden) {
			t.Fatalf("offline consumer retains direct builder reference %q", forbidden)
		}
	}

	buildEnd := strings.Index(workflow[buildStart:], "- name: Push immutable staging image")
	if buildEnd < 0 {
		t.Fatal("staging workflow is missing image publication after the final image build")
	}
	buildStep := workflow[buildStart : buildStart+buildEnd]
	if !strings.Contains(buildStep, "docker build --platform=linux/amd64 --pull=false --network=none --tag \"$image\" .") {
		t.Fatal("final staging image build is not explicitly local and network-isolated")
	}
	if strings.Contains(buildStep, " --pull ") {
		t.Fatal("final staging image build retains an implicit builder refresh")
	}
}

func TestEvidenceVerifierRequiresEGP1SubjectsAndVerifiedPromotion(t *testing.T) {
	verifierPath := filepath.Join("..", "..", ".github", "actions", "verify-broker-evidence", "action.yml")
	contents, err := os.ReadFile(verifierPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", verifierPath, err)
	}

	verifier := string(contents)
	for _, required := range []string{
		"require-promotion:",
		"promotion-lane:",
		"promotion-workflow:",
		"require-promotion must be true or false",
		"verify_subject()",
		"source.subject.json",
		"dependency-resolution.subject.json",
		"build.subject.json",
		"artifact.subject.json",
		"promotion.subject.json",
		"promotion-approval.json",
		"evidence-graph/v1",
		"canonical_payload_digest",
		"utf8-json-sorted-keys-v1",
		"cosign verify-blob",
		".dependency_resolution.admission.status == \"verified\"",
		".lifecycle.status == \"verified\"",
		".lifecycle.phase_references.source_artifact_canonical_payload_digest",
		"promotion_certificate_identity",
		"@refs/heads/main",
	} {
		if !strings.Contains(verifier, required) {
			t.Fatalf("evidence verifier is missing %q", required)
		}
	}
	for _, forbidden := range []string{
		"schema_version == 2",
		"manifest.json",
		".subjects.",
		"promotion.json",
	} {
		if strings.Contains(verifier, forbidden) {
			t.Fatalf("evidence verifier retains obsolete %q", forbidden)
		}
	}
}

func TestProductionDeploymentsRequireLaneBoundPromotionEvidence(t *testing.T) {
	root := filepath.Join("..", "..")
	for _, testCase := range []struct {
		path               string
		lane               string
		promotionWorkflow  string
		deploymentWorkflow string
		environment        string
	}{
		{
			path:               filepath.Join(".github", "workflows", "gcp-broker-production.yml"),
			lane:               "release-automation",
			promotionWorkflow:  ".github/workflows/gcp-broker-production-promotion.yml",
			deploymentWorkflow: ".github/workflows/gcp-broker-production.yml",
			environment:        "gcp-broker-production",
		},
		{
			path:               filepath.Join(".github", "workflows", "gcp-reconciliation-publisher-production.yml"),
			lane:               "reconciliation-publisher",
			promotionWorkflow:  ".github/workflows/gcp-reconciliation-publisher-promotion.yml",
			deploymentWorkflow: ".github/workflows/gcp-reconciliation-publisher-production.yml",
			environment:        "gcp-reconciliation-publisher-deployment",
		},
		{
			path:               filepath.Join(".github", "workflows", "gcp-release-credential-verification-production.yml"),
			lane:               "release-credential-verification",
			promotionWorkflow:  ".github/workflows/gcp-release-credential-verification-promotion.yml",
			deploymentWorkflow: ".github/workflows/gcp-release-credential-verification-production.yml",
			environment:        "gcp-release-credential-verification-deployment",
		},
		{
			path:               filepath.Join(".github", "workflows", "gcp-hotfix-delivery-production.yml"),
			lane:               "hotfix-delivery",
			promotionWorkflow:  ".github/workflows/gcp-hotfix-delivery-promotion.yml",
			deploymentWorkflow: ".github/workflows/gcp-hotfix-delivery-production.yml",
			environment:        "gcp-hotfix-delivery-deployment",
		},
		{
			path:               filepath.Join(".github", "workflows", "gcp-hotfix-propagation-publisher-production.yml"),
			lane:               "hotfix-propagation-publisher",
			promotionWorkflow:  ".github/workflows/gcp-hotfix-propagation-publisher-promotion.yml",
			deploymentWorkflow: ".github/workflows/gcp-hotfix-propagation-publisher-production.yml",
			environment:        "gcp-hotfix-propagation-publisher-deployment",
		},
	} {
		t.Run(testCase.lane, func(t *testing.T) {
			contents, err := os.ReadFile(filepath.Join(root, testCase.path))
			if err != nil {
				t.Fatalf("ReadFile(%q) error = %v", testCase.path, err)
			}

			workflow := string(contents)
			for _, required := range []string{
				"uses: ./.github/actions/verify-broker-evidence",
				"require-promotion: true",
				"promotion-lane: " + testCase.lane,
				"promotion-workflow: " + testCase.promotionWorkflow,
				"uses: ./.github/actions/record-broker-deployment-evidence",
				"lane: " + testCase.lane,
				"workflow: " + testCase.deploymentWorkflow,
				"source-ref: refs/heads/main",
				"environment: " + testCase.environment,
			} {
				if !strings.Contains(workflow, required) {
					t.Fatalf("production workflow is missing %q", required)
				}
			}
			deployIndex := strings.Index(workflow, "gcloud run deploy")
			recordIndex := strings.Index(workflow, "uses: ./.github/actions/record-broker-deployment-evidence")
			if deployIndex < 0 || recordIndex < 0 || deployIndex > recordIndex {
				t.Fatal("production deployment evidence is not recorded after the Cloud Run mutation")
			}
		})
	}
}

func TestPromotionWorkflowsRecordLaneBoundEGP1Subjects(t *testing.T) {
	root := filepath.Join("..", "..")
	for _, testCase := range []struct {
		path              string
		lane              string
		promotionWorkflow string
		environment       string
	}{
		{
			path:              filepath.Join(".github", "workflows", "gcp-broker-production-promotion.yml"),
			lane:              "release-automation",
			promotionWorkflow: ".github/workflows/gcp-broker-production-promotion.yml",
			environment:       "gcp-broker-production",
		},
		{
			path:              filepath.Join(".github", "workflows", "gcp-reconciliation-publisher-promotion.yml"),
			lane:              "reconciliation-publisher",
			promotionWorkflow: ".github/workflows/gcp-reconciliation-publisher-promotion.yml",
			environment:       "gcp-reconciliation-publisher-deployment",
		},
		{
			path:              filepath.Join(".github", "workflows", "gcp-release-credential-verification-promotion.yml"),
			lane:              "release-credential-verification",
			promotionWorkflow: ".github/workflows/gcp-release-credential-verification-promotion.yml",
			environment:       "gcp-release-credential-verification-deployment",
		},
		{
			path:              filepath.Join(".github", "workflows", "gcp-hotfix-delivery-promotion.yml"),
			lane:              "hotfix-delivery",
			promotionWorkflow: ".github/workflows/gcp-hotfix-delivery-promotion.yml",
			environment:       "gcp-hotfix-delivery-deployment",
		},
		{
			path:              filepath.Join(".github", "workflows", "gcp-hotfix-propagation-publisher-promotion.yml"),
			lane:              "hotfix-propagation-publisher",
			promotionWorkflow: ".github/workflows/gcp-hotfix-propagation-publisher-promotion.yml",
			environment:       "gcp-hotfix-propagation-publisher-deployment",
		},
	} {
		t.Run(testCase.lane, func(t *testing.T) {
			contents, err := os.ReadFile(filepath.Join(root, testCase.path))
			if err != nil {
				t.Fatalf("ReadFile(%q) error = %v", testCase.path, err)
			}

			workflow := string(contents)
			for _, required := range []string{
				"uses: ./.github/actions/record-broker-promotion-evidence",
				"target-image: ${{ steps.promotion.outputs.image }}",
				"lane: " + testCase.lane,
				"workflow: " + testCase.promotionWorkflow,
				"source-ref: refs/heads/main",
				"environment: " + testCase.environment,
				"promoter-service-account: ${{ env.GCP_PROMOTER_SERVICE_ACCOUNT }}",
			} {
				if !strings.Contains(workflow, required) {
					t.Fatalf("promotion workflow is missing %q", required)
				}
			}
			for _, forbidden := range []string{
				"promotion.json",
				"manifest_sha256",
				"schema_version: 1",
			} {
				if strings.Contains(workflow, forbidden) {
					t.Fatalf("promotion workflow retains obsolete %q", forbidden)
				}
			}
		})
	}
}

func TestEvidenceGraphSubjectSealerContract(t *testing.T) {
	sealerPath := filepath.Join("..", "..", ".github", "actions", "seal-evidence-subject", "action.yml")
	contents, err := os.ReadFile(sealerPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", sealerPath, err)
	}

	sealer := string(contents)
	for _, required := range []string{
		"sigstore/cosign-installer@6f9f17788090df1f26f669e9d70d6ae9567deba6",
		"cosign-release: v3.1.3",
		"*.subject.payload.json",
		"jq -cS . \"$payload\"",
		"cosign sign-blob --yes --bundle \"$bundle\" \"$canonical\"",
		"canonical_payload_digest",
		"utf8-json-sorted-keys-v1",
		"immutable_reference:",
		"issued_at: $signed_at",
		"credential-like value",
		"rm -f \"$canonical\" \"$payload\"",
	} {
		if !strings.Contains(sealer, required) {
			t.Fatalf("evidence-graph subject sealer is missing %q", required)
		}
	}
	for _, forbidden := range []string{
		"gcloud artifacts generic upload",
		"gcloud run deploy",
		"docker build",
	} {
		if strings.Contains(sealer, forbidden) {
			t.Fatalf("evidence-graph subject sealer contains forbidden %q", forbidden)
		}
	}
}

func TestPromotionAndDeploymentEvidenceActionsRemainAppendOnly(t *testing.T) {
	root := filepath.Join("..", "..")
	for _, testCase := range []struct {
		name      string
		path      string
		required  []string
		forbidden []string
	}{
		{
			name: "promotion",
			path: filepath.Join(".github", "actions", "record-broker-promotion-evidence", "action.yml"),
			required: []string{
				"artifact.subject.json",
				".lifecycle.status == \"verified\"",
				"promotion-approval.json",
				"promotion.subject.payload.json",
				"relation_type: \"promotes\"",
				"uses: ./.github/actions/seal-evidence-subject",
				"gcloud artifacts generic upload",
				"--skip-existing",
			},
			forbidden: []string{
				"docker build",
				"gcloud run deploy",
				"latest",
			},
		},
		{
			name: "deployment",
			path: filepath.Join(".github", "actions", "record-broker-deployment-evidence", "action.yml"),
			required: []string{
				"gcloud run services describe",
				"gcloud run revisions describe",
				"deployment-health.json",
				"deployment-approval.json",
				"deployment.subject.payload.json",
				"relation_type: \"deploys\"",
				"deployed_digest=\"${deployed_image##*@}\"",
				"test \"$deployed_digest\" = \"$digest\" || fail deployed-digest-mismatch",
				"uses: ./.github/actions/seal-evidence-subject",
				"gcloud artifacts generic upload",
				"--skip-existing",
			},
			forbidden: []string{
				"docker build",
				"gcloud run deploy",
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			contents, err := os.ReadFile(filepath.Join(root, testCase.path))
			if err != nil {
				t.Fatalf("ReadFile(%q) error = %v", testCase.path, err)
			}

			action := string(contents)
			for _, required := range testCase.required {
				if !strings.Contains(action, required) {
					t.Fatalf("%s evidence action is missing %q", testCase.name, required)
				}
			}
			for _, forbidden := range testCase.forbidden {
				if strings.Contains(action, forbidden) {
					t.Fatalf("%s evidence action contains forbidden %q", testCase.name, forbidden)
				}
			}
		})
	}
}

func TestDeploymentEvidenceUsesDigestIdentityAndBoundedFailureCodes(t *testing.T) {
	deploymentPath := filepath.Join("..", "..", ".github", "actions", "record-broker-deployment-evidence", "action.yml")
	contents, err := os.ReadFile(deploymentPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", deploymentPath, err)
	}

	action := string(contents)
	for _, required := range []string{
		"fail() {",
		"deployment evidence validation failed: $1",
		"artifact-subject-missing",
		"artifact-subject-invalid",
		"artifact-integrity-missing",
		"runtime-revision-missing",
		"runtime-revision-document-missing",
		"runtime-revision-condition-missing",
		"runtime-revision-not-ready",
		"revision_document=\"$(gcloud run revisions describe \"$revision\"",
		"--format=json)",
		"printf '%s' \"$revision_document\" | jq -er",
		".status.conditions[]?",
		".conditions[]?",
		"map(select(.type == \"Ready\"))",
		"expected exactly one Ready condition",
		".status // .state // empty",
		"True|CONDITION_SUCCEEDED)",
		".spec.containers[]?",
		".containers[]?",
		"expected exactly one container",
		"ready_condition_status: $ready_status",
		"deployed-image-missing",
		"deployed-image-not-digest-pinned",
		"deployed-digest-mismatch",
		"deployed_digest=\"${deployed_image##*@}\"",
		"[[ \"$deployed_digest\" =~ ^sha256:[0-9a-f]{64}$ ]]",
		"test \"$deployed_digest\" = \"$digest\" || fail deployed-digest-mismatch",
	} {
		if !strings.Contains(action, required) {
			t.Fatalf("deployment evidence action is missing %q", required)
		}
	}
	for _, forbidden := range []string{
		"test \"$deployed_image\" = \"$IMAGE\"",
		"deployment evidence validation failed: $IMAGE",
		"deployment evidence validation failed: $deployed_image",
		"status.conditions[?type=Ready].status",
		"status.conditions[?type=Ready].state",
		"conditions[?type=Ready].state",
		"status.conditions[?type=\"Ready\"].status",
		"spec.containers.image",
		"test \"$ready_state\" = \"CONDITION_SUCCEEDED\"",
		"spec.containers[0].image",
		"containers[0].image",
	} {
		if strings.Contains(action, forbidden) {
			t.Fatalf("deployment evidence action retains forbidden %q", forbidden)
		}
	}
	if count := strings.Count(action, "gcloud run revisions describe"); count != 1 {
		t.Fatalf("deployment evidence action invokes gcloud run revisions describe %d times, want 1", count)
	}
}

#!/usr/bin/env bash
set -euo pipefail

fail() {
  echo "builder-evidence-validation-failed:$1" >&2
  exit 1
}

for variable in \
  BUILDER_IMAGE \
  GCP_PROJECT_ID \
  GCP_REGION \
  BUILDER_ARTIFACT_REPOSITORY \
  BUILDER_EVIDENCE_ARTIFACT_REPOSITORY \
  BUILDER_SIGNER_REPOSITORY \
  BUILDER_SIGNER_WORKFLOW \
  BUILDER_SOURCE_REF; do
  test -n "${!variable:-}" || fail "missing-${variable,,}"
done

builder_digest="sha256:${BUILDER_IMAGE##*@sha256:}"
[[ "$builder_digest" =~ ^sha256:[0-9a-f]{64}$ ]] || fail builder-digest-invalid

expected_builder_image="${GCP_REGION}-docker.pkg.dev/${GCP_PROJECT_ID}/${BUILDER_ARTIFACT_REPOSITORY}/builder@${builder_digest}"
test "$BUILDER_IMAGE" = "$expected_builder_image" || fail builder-image-invalid

evidence_version="${builder_digest#sha256:}"
evidence_package="broker-builder-evidence"
evidence_directory="$(mktemp -d)"
cleanup() {
  rm -rf "$evidence_directory"
}
trap cleanup EXIT

gcloud artifacts generic download \
  --project="$GCP_PROJECT_ID" \
  --location="$GCP_REGION" \
  --repository="$BUILDER_EVIDENCE_ARTIFACT_REPOSITORY" \
  --package="$evidence_package" \
  --version="$evidence_version" \
  --destination="$evidence_directory"

for required_file in \
  builder.subject.json \
  builder.subject.integrity.sigstore.json \
  builder.spdx.json \
  builder.signature.json \
  builder.provenance.intoto.jsonl \
  builder.sbom.intoto.jsonl \
  builder.policy.json \
  builder.approval.json \
  builder.revocation.json; do
  test -f "$evidence_directory/$required_file" || fail "missing-${required_file}"
done

builder_subject="$evidence_directory/builder.subject.json"
builder_subject_id="$(jq -er '.subject.id' "$builder_subject")" || fail builder-subject-id-invalid
builder_source_commit="$(jq -er '.source.commit | select(test("^[0-9a-f]{40}$"))' "$builder_subject")" || fail builder-source-commit-invalid
signature_file="$(jq -er '.integrity.signature.immutable_reference.file' "$builder_subject")" || fail builder-integrity-reference-invalid
signature_bundle="$evidence_directory/$signature_file"
test -f "$signature_bundle" || fail builder-integrity-bundle-missing

canonical_payload="$(mktemp)"
trap 'rm -f "$canonical_payload"; cleanup' EXIT
jq 'del(.integrity)' "$builder_subject" | jq -cS . > "$canonical_payload"
canonical_payload_digest="sha256:$(sha256sum "$canonical_payload" | awk '{print $1}')"
signature_bundle_digest="sha256:$(sha256sum "$signature_bundle" | awk '{print $1}')"
certificate_identity="https://github.com/${BUILDER_SIGNER_REPOSITORY}/${BUILDER_SIGNER_WORKFLOW}@${BUILDER_SOURCE_REF}"

jq -e \
  --arg builder_digest "$builder_digest" \
  --arg builder_image "$BUILDER_IMAGE" \
  --arg subject_id "$builder_subject_id" \
  --arg canonical_payload_digest "$canonical_payload_digest" \
  --arg signature_bundle_digest "$signature_bundle_digest" \
  --arg signature_file "$signature_file" \
  --arg certificate_identity "$certificate_identity" \
  --arg evidence_repository "$BUILDER_EVIDENCE_ARTIFACT_REPOSITORY" \
  --arg evidence_package "$evidence_package" \
  --arg evidence_version "$evidence_version" \
  '.schema == "evidence-graph/v1" and
   .document.type == "subject" and
   .document.issuer == $certificate_identity and
   .subject.type == "artifact" and
   .subject.id == $subject_id and
   .subject.primary_digest == $builder_digest and
   .artifact.format == "oci" and
   .artifact.kind == "builder" and
   .artifact.reference == $builder_image and
   .artifact.digest == $builder_digest and
   .policy.decision == "verified" and
   .policy.approval.status == "verified" and
   .lifecycle.status == "verified" and
   .lifecycle.target_lane == "builder" and
   .integrity.canonicalization == "utf8-json-sorted-keys-v1" and
   .integrity.canonical_payload_digest == $canonical_payload_digest and
   .integrity.signature.evidence_type == "signature" and
   .integrity.signature.subject_id == $subject_id and
   .integrity.signature.digest == $signature_bundle_digest and
   .integrity.signature.issuer == $certificate_identity and
   .integrity.signature.immutable_reference.repository == $evidence_repository and
   .integrity.signature.immutable_reference.package == $evidence_package and
   .integrity.signature.immutable_reference.version == $evidence_version and
   .integrity.signature.immutable_reference.file == $signature_file' \
  "$builder_subject" >/dev/null || fail builder-subject-invalid

verify_evidence() {
  local evidence_type="$1"
  local evidence_file="$2"
  local evidence_path="$evidence_directory/$evidence_file"
  local evidence_digest

  evidence_digest="sha256:$(sha256sum "$evidence_path" | awk '{print $1}')"
  jq -e \
    --arg evidence_type "$evidence_type" \
    --arg subject_id "$builder_subject_id" \
    --arg evidence_repository "$BUILDER_EVIDENCE_ARTIFACT_REPOSITORY" \
    --arg evidence_package "$evidence_package" \
    --arg evidence_version "$evidence_version" \
    --arg evidence_file "$evidence_file" \
    --arg evidence_digest "$evidence_digest" \
    'any(.evidence[];
      .evidence_type == $evidence_type and
      .subject_id == $subject_id and
      .immutable_reference.repository == $evidence_repository and
      .immutable_reference.package == $evidence_package and
      .immutable_reference.version == $evidence_version and
      .immutable_reference.file == $evidence_file and
      .digest == $evidence_digest and
      (.issuer | type == "string" and length > 0) and
      (.issued_at | type == "string" and length > 0)
    )' \
    "$builder_subject" >/dev/null || fail "builder-${evidence_type}-evidence-invalid"
}

verify_evidence sbom builder.spdx.json
verify_evidence signature builder.signature.json
verify_evidence provenance builder.provenance.intoto.jsonl
verify_evidence attestation builder.sbom.intoto.jsonl
verify_evidence policy builder.policy.json
verify_evidence approval builder.approval.json
verify_evidence revocation builder.revocation.json

cosign verify-blob "$canonical_payload" \
  --bundle "$signature_bundle" \
  --certificate-identity="$certificate_identity" \
  --certificate-oidc-issuer="https://token.actions.githubusercontent.com" >/dev/null

cosign verify \
  --certificate-identity="$certificate_identity" \
  --certificate-oidc-issuer="https://token.actions.githubusercontent.com" \
  "$BUILDER_IMAGE" >/dev/null

gh attestation verify \
  "oci://${BUILDER_IMAGE}" \
  --repo "$BUILDER_SIGNER_REPOSITORY" \
  --signer-workflow "${BUILDER_SIGNER_REPOSITORY}/${BUILDER_SIGNER_WORKFLOW}" \
  --source-ref "$BUILDER_SOURCE_REF" \
  --source-digest "$builder_source_commit" \
  --deny-self-hosted-runners \
  --bundle-from-oci >/dev/null

gh attestation verify \
  "oci://${BUILDER_IMAGE}" \
  --repo "$BUILDER_SIGNER_REPOSITORY" \
  --signer-workflow "${BUILDER_SIGNER_REPOSITORY}/${BUILDER_SIGNER_WORKFLOW}" \
  --source-ref "$BUILDER_SOURCE_REF" \
  --source-digest "$builder_source_commit" \
  --deny-self-hosted-runners \
  --predicate-type "https://spdx.dev/Document/v2.3" \
  --bundle-from-oci >/dev/null

echo "image=$BUILDER_IMAGE" >> "$GITHUB_OUTPUT"
echo "subject-id=$builder_subject_id" >> "$GITHUB_OUTPUT"
echo "canonical-payload-digest=$canonical_payload_digest" >> "$GITHUB_OUTPUT"

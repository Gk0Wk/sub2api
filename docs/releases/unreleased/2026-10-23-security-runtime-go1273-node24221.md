---
type: fix
scope: runtime
audience: internal
summary: Pin Sub2API build paths to patched Go 1.27.3 and Node 24.22.1
breaking: false
demo_ready: false
tests: "docker buildx build --check --platform linux/amd64 -f Dockerfile .; docker buildx build --check --platform linux/amd64 -f deploy/Dockerfile .; docker buildx build --check --platform linux/amd64 -f backend/Dockerfile backend; go test -tags=unit ./... -run '^$' -count=1; go vet ./..."
artifacts: Dockerfile; deploy/Dockerfile; backend/Dockerfile; backend/go.mod; .github/workflows/backend-ci.yml; .github/workflows/release.yml; .github/workflows/security-scan.yml; DEV_GUIDE.md
---

## What changed

The root, deployment, and backend build paths now use immutable official Go
1.27.3 and Node 24.22.1 image indexes. `backend/go.mod` and the CI/release
workflows enforce the same versions.

## Why it matters

The patched Go toolchain carries the `crypto/x509` trust-list fix associated
with CVE-2026-111943. The Node runtime pin keeps the build and release paths
on the patched 24.x line associated with CVE-2026-111497.

## Demo posture / limitations

This is a source-only candidate from current `origin/main`. It does not push,
publish, restart, or replace a running image. Full vulnerability scanning and
the complete Go unit suite remain release gates; the existing candidate run
passed Dockerfile checks, version probes, vet, and focused Linux tests while
retaining two shared-state service test failures for owner review.

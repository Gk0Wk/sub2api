---
type: fix
scope: runtime
audience: internal
summary: Pin the live Sub2API lineage to patched Go and Node runtime images
breaking: false
demo_ready: false
tests:
  - "docker buildx build --check --platform linux/amd64 -f Dockerfile ."
  - "docker buildx build --check --platform linux/amd64 -f deploy/Dockerfile ."
  - "docker buildx build --check --platform linux/amd64 -f backend/Dockerfile backend"
  - "docker buildx build --check --platform linux/arm64 -f deploy/Dockerfile ."
  - "docker run --rm node:24.22.1-alpine3.24@sha256:06d7adf3e0ef61a13f1d28682726366f97a15aa1888c934d3ecec86c33918b53 node --version"
  - "go test -tags=unit ./... -run '^$' -count=1"
  - "go vet ./..."
  - "go test -tags=unit ./internal/repository -run TestPgDumper -count=1 (pinned Go Linux container)"
  - "go test -tags=unit ./... -count=1 (partial: known unrelated service package failures)"
artifacts:
  - Dockerfile
  - deploy/Dockerfile
  - backend/Dockerfile
  - backend/go.mod
  - .github/workflows/backend-ci.yml
  - .github/workflows/release.yml
  - .github/workflows/security-scan.yml
  - DEV_GUIDE.md
  - TODO.md
---

## What changed

The source baseline `cd8555dd6fea76b778bef8e0a62ade2b2ff799dd` now uses
Go 1.27.3 and Node 24.22.1 in the root, deployment, and backend build paths.

The Go builder is pinned to the official multi-platform image index
`sha256:d530cf60723abf1ffae0cd16bb546b379ec0550d2f3815e5b16ed69dbc00fba4`.
The Node builder is pinned to the official multi-platform image index
`sha256:06d7adf3e0ef61a13f1d28682726366f97a15aa1888c934d3ecec86c33918b53`.
Both indexes cover Linux amd64 and arm64. `backend/go.mod` declares Go 1.27.3
so local and container toolchains enforce the same minimum. CI's Go version
assertions and frontend setup use the same Go and Node versions.

## Why it matters

Go 1.27.3 contains the fixed `crypto/x509` trust-list behavior for
GO-2026-6809 / CVE-2026-111943. Node 24.22.1 contains the TLS, DNS, c-ares,
NSS CA, and OpenSSL fixes associated with CVE-2026-111497. The live
codex-relay image was built from the `internal-pr-fixes-20261021` lineage and
still runs the older build inputs until a new image is promoted.

## Demo posture / limitations

This change is an isolated source candidate. It does not push, publish an
image, replace the running codex-relay image, restart a service, or authorize a
production rollout. Dockerfile validation and Node version execution passed;
the three `pg_dump` tests that fail on the Windows host pass in the pinned Go
Linux container. The full Go unit command still reports two `internal/service`
failures when many package tests share process state; both failing cases pass
when run alone. A CI rerun must clear that existing package-test instability
before release.

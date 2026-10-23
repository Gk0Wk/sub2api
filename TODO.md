# TODO

## Active: OAuth refresh deadline boundary (2026-10-23)

- [x] Confirm source ownership and isolate branch from `c9b820271`;
      production remains HOLD, no push/deploy.
- [x] Reproduce late credentials with an explicit deadline and an open Done
      channel using a controlled clock; verify both persistence paths.
- [x] Check cancellation and deadlines before first persistence, preserving
      durable-success cleanup and provider error classification.
- [x] Run focused RED/GREEN, old/new dependency controls, then unit/integration
      checks under 2 CPU / 6 GiB in CTYunGPU1 `gpu/gpu-dev`; source, logs and
      caches are centralized at `/data/tmp/ispark-security-202610/sub2api-deadline-b783`.
      The first 3 GiB compilation was OOM-killed before any test ran.
- [ ] Build and scan an exact-source candidate, obtain independent review,
      hand evidence to the original production coordinator, and clean caches.

Focused old/new dependency RED/GREEN and independent source review are complete.
Captured full unit: 57 packages passed, 51 had no tests, schema failed on a
readonly `.entc` fixture. Its byte-identical writable recheck passed, as did the
complete service integration-tag suite. The first full unit output was rotated
by Docker and remains incomplete evidence. Keep original failures and HOLD.

The approved scope is private implementation and deterministic regressions;
public contracts, schema and production are outside this unit. Stop the
affected work on write conflicts, repeated verification failure, or a need to
change persistence/rotation semantics. Detailed acceptance and results belong
in the release entry; this branch does not supersede the production window.

## Active: Sub2API security runtime baseline

- [x] Create an isolated worktree from the live lineage baseline
      `cd8555dd6fea76b778bef8e0a62ade2b2ff799dd`.
- [x] Pin Go 1.27.3 and Node 24.22.1 by immutable multi-platform image index
      in the root, deployment, and backend build paths.
- [x] Verify Dockerfile amd64/arm64 resolution and run the pinned Node image.
- [x] Pinned Go 1.27.3 Linux full unit suite and vet pass; internal/service
      passes in the full suite. Runtime image built, pushed and promoted to
      internal codex-relay; see the
      [promotion record](docs/releases/unreleased/2026-10-22-outbound-runtime-promotion.md).
- [ ] Run a separate full vulnerability scan. The optional SBOM scanner image
      pull failed; the runtime version checks do not represent a full scan.
- [ ] Rebase the remaining local Sub2API branches, then remove the superseded
      security worktree only after the accepted branch and evidence are merged.

### Local branch inventory

| Branch | Current runtime-input gap |
| --- | --- |
| `main` | Go 1.27.2 and unpinned Node in all build paths |
| `integration/internal-pr-fixes-20261021` | Go 1.27.2 and unpinned Node in all build paths |
| `fix/model-catalog-transient-availability-20261021` | Go 1.27.2 and unpinned Node in all build paths |
| `feat/openai-fast-capability-fallback` | Go 1.27.2 and unpinned Node in all build paths |
| `integration/internal-native-model-20261021` | Go 1.27.2 and unpinned Node in all build paths |
| `fix/model-catalog-temporary-unavailability-20261022` | root Go is 1.27.3, but deploy/backend Go and Node remain old |
| `fix/model-upstream-balance-diagnostics-20261022` | root Go is 1.27.3, but deploy/backend Go and Node remain old |
| `Sttot/security-runtime-20261022` | preserves an earlier dirty WIP; reconcile against `Sttot/security-runtime-clean-20261022` |

Remote-only branches were not changed or audited in this unit.

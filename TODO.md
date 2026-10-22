# TODO

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

# TODO

## Active: Sub2API security runtime baseline

- [x] Create an isolated worktree from the live lineage baseline
      `cd8555dd6fea76b778bef8e0a62ade2b2ff799dd`.
- [x] Pin Go 1.27.3 and Node 24.22.1 by immutable multi-platform image index
      in the root, deployment, and backend build paths.
- [x] Verify Dockerfile amd64/arm64 resolution and run the pinned Node image.
- [ ] Re-run the full Go test suite in the project CI or a Linux test
      environment. The three Windows `sh` failures pass in the pinned Linux
      container; two service tests pass alone but remain unstable in the full
      package run and need a CI retest or separate pre-existing fix.
- [ ] Build and scan a new immutable codex-relay image, then follow the
      deployment owner's rollout and rollback procedure.
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

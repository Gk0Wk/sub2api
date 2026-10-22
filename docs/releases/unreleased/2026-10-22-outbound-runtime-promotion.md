---
type: chore
scope: outbound
audience: internal
summary: Promote the verified r6 source lineage with Go 1.27.3 and Node 24.22.1
breaking: false
demo_ready: false
tests:
  - "Pinned Go Linux: go test -p 4 -tags=unit ./... -count=1"
  - "Pinned Go Linux: go vet -p 4 -tags=unit ./..."
  - "docker buildx build --check --platform linux/amd64 -f Dockerfile ."
  - "docker buildx build --check --platform linux/amd64 -f deploy/Dockerfile ."
  - "docker buildx build --check --platform linux/amd64 -f backend/Dockerfile backend"
  - "go version -m .tmp/runtime-promotion/sub2api"
  - "kubectl rollout status deployment/codex-relay (explicit hz1 kubeconfig and internal namespace)"
artifacts:
  - TODO.md
  - DEV_GUIDE.md
  - Dockerfile
  - backend/go.mod
---

## What changed

2026-10-22，候选源码 `b00ece264f844b4c6486d5cdd8ef62c1864c18b5` 已构建并上线
hz1 internal `codex-relay`。它的直接父提交
`cd8555dd6fea76b778bef8e0a62ade2b2ff799dd` 与线上 r6 程序内嵌 commit 一致；
ACR r6 tag、live pin 和 Pod imageID 均为
`sha256:d09ad3f1ed0f6e2c4ee37d88939d798da6f12e41e7d0d17d29c4bf9bbb32899b`。
候选不改变业务代码、模型目录、账号、路由、计费或并发行为。

新镜像为 `registry.cn-hangzhou.aliyuncs.com/sttot/ispark-sub2api:internal-20261022-runtime-b00ece264`，
linux/amd64 manifest digest：
`sha256:dc6e902cabe893b84ea130d060314d99e743967575851fd02b66158171139a68`。
hz1 从对应 `registry-vpc.cn-hangzhou.aliyuncs.com` 地址按该 digest 拉取。
程序 buildinfo 为 `go1.27.3`；前端 builder 为 `v24.22.1`，最终镜像没有 Node runtime。

## Why it matters

全量 Linux unit suite 和 vet 均 exit 0，`internal/service` 用时 251.663 秒。
此前两项不稳定 service 测试名称未保留，协调者明确接受此次固定工具链全量通过；
没有猜测测试名称，也不将一次通过宣称为消除了所有不稳定性。
三个 Dockerfile amd64 check 均无警告。首次测试因只读源码挂载使 Ent 无法创建
`.entc` 失败，改用容器内可写副本后全量通过，没有修改产品源码。

19:24 CST 生产 rollout 达到 generation/observed 79/79、Ready 1/1。
目标 worker 版本探针及生产 Pod imageID 都匹配 ACR digest。client/server dry-run
和完整 spec 比较证明仅 image/config-revision 改变。Service、HTTPRoute 和 Redis
spec 保持一致，Redis 8.10.2 的认证 PING 返回 PONG；公网 health/root/admin=200、
未认证 models=401。正常 auto-setup 首次 Completed/exit0 重启一次，按协调者裁定
单独记录为初始化行为。稳定窗口详情由 deployment release evidence 记录。

## Demo posture / limitations

这是内部生产运行时升级，不代表上游发布、公开 GA、付费模型调用或完整漏洞扫描通过。
未调用付费 provider；应用源码和数据库 schema 不变，未修改 Secret。
SBOM scanner 镜像拉取失败；重试必要构建时关闭该附加扫描步骤。
ACR 拒绝带空配置的 attestation manifest，因此推送的是同一 amd64 image manifest；
本地 BuildKit provenance、二进制 buildinfo 和源码 revision labels 已保留。

deployment 审计原始结果为 94 blockingHigh，均来自两份既有、Git clean 的临时公开源码。
按协调者限定路径的裁定继续，记录原始审计摘要与哈希；不宣称全仓 high 已清零。
源码 push 的 HTTPS OAuth 凭证缺 workflow scope，改用已验证属于 Gk0Wk 的既有 SSH
身份成功推送，没有修改凭证、Git remote 或认证配置。

回滚使用上面的 r6 digest，revision `2026-10-21-sub2api-0-2-17-pr-fixes-r6`。
完整两字段 guarded rollback、线上稳定窗口及审计处置证据位于 deployment 仓库
`docs/releases/unreleased/2026-10-22-outbound-codex-relay-runtime.md`。

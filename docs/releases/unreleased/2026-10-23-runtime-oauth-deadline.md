---
type: fix
scope: runtime
audience: internal
summary: 在 OAuth 首次凭据写入前检查显式截止时间，防止取消通知延迟时写入迟到结果
breaking: false
demo_ready: false
tests:
  - "Focused RED with original logic, original and patched dependencies: 24 fail / 7 pass events in each run"
  - "Focused GREEN with fixed logic, original and patched dependencies: 31 pass events in each run, exit 0, no OOM"
  - "Independent source review and git diff --check: passed"
  - "Full unit captured run: 57 packages passed, 51 without tests, schema package failed on readonly .entc; byte-identical writable schema recheck passed"
  - "go test -mod=readonly -tags=integration -p 1 -parallel 1 -count=1 ./internal/service: passed, 8513 pass and 3 skip test events"
artifacts:
  - backend/internal/service/oauth_refresh_deadline.go
  - backend/internal/service/oauth_refresh_api.go
  - backend/internal/service/token_refresh_service.go
  - backend/internal/service/token_refresh_deadline_test.go
  - backend/internal/service/oauth_refresh_api_test.go
  - backend/internal/service/token_refresh_pool_health_test.go
  - backend/internal/service/token_refresh_service_test.go
  - backend/internal/service/README.md
---

## What changed

基于后端安全依赖提交 `c9b820271800d8b686d8d6b57f7650315f6f931d`（完整基线以
Git parent 为准），在后台刷新兼容路径首次写入前、统一 OAuth API 首次普通写入
或 Grok CAS 前检查 `ctx.Err()` 和显式 deadline。取消错误优先，恰好到达截止时间
也返回超时；生产使用 `time.Now`，测试使用同一私有可控时钟，无公共配置变动。

父级取消或到期不触发账号封禁或冷却。单次尝试超时、provider containment /
configuration 错误保持原分类；CAS 结果不确定不能转成可重试超时。已确认持久化
成功后仅清理阶段超过 attempt deadline，仍同步状态并返回成功；原 shortCircuit
的 skip 语义保持。此变更不提供跨数据库事务的原子截止时间保证。

两个平台、兼容与统一两条路径覆盖 deadline 前、恰好到达、之后、父级到期、
父级取消及迟到 invalid_grant。用例在 `Err()==nil`、`Done()` 尚未关闭时推进
同一时钟，替代原先依赖短 Sleep 的调度测试。另验证 CAS 不确定和持久化后清理
越期，不新增重试、不清掉已有冷却状态。

## Why it matters

[Go Context 文档](https://pkg.go.dev/context)允许取消通知异步关闭；只看 `Err()`
不能证明显式截止时间尚未到达。本次修复拦住已过预算的结果进入首次写入，且不把
已经成功的轮换误当作失败重试。没有调整依赖、schema、token 轮换协议或接口。

旧逻辑在原/新依赖均产生 24 fail / 7 pass，修复后均为 31 pass；计数包含父子
测试事件，不是 24 个漏洞。首次 3 GiB 编译 OOM 时尚无测试执行，不计作 RED。
独立 reviewer 核验七个实现/测试文件与 GREEN 快照一致，未发现必须修改项。

## Demo posture / limitations

这次不代表已经发布、生产健康或完整业务验收。源码写入归本安全主线，生产窗口
仍归原协管；本单元禁止 push/deploy。完整 service integration-tag 测试已通过，
8513 pass / 3 skip测试事件、exit0、无OOM。本轮没有重复全部仓库的集成测试，
也没有把service测试视为新的真实数据库或业务验收。精确源码镜像和扫描仍待完成。

运行环境为 CTYunGPU1（SSH 别名 `CTYunGPT1`）的 `gpu/gpu-dev`，每容器
2 CPU / 6 GiB，Go `-p 1 -parallel 1`。早期 `red-old` 与 `green` 容器运行窗口
曾重叠约10.80秒，不能声称全程严格串行；两次均未OOM。后续执行脚本增加活跃
任务容器检查，拒绝重叠启动，后续完整unit、schema复验及service integration
的运行窗口不重叠。统一材料目录为
`/data/tmp/ispark-security-202610/sub2api-deadline-b783`，实际位于大容量数据卷。
首轮完整单元退出 1，Docker 日志轮转导致前半段缺失；复验改为直接写集中目录。
只读挂载阻挡 schema 测试创建 `.entc` 的失败单独保留，不通过修改产品代码绕过。
完整捕获的 unit 结果为57包pass、51包无测试、1包fail；唯一失败为上述schema
夹具。42个文件哈希一致的可写副本复验该包通过，没有改变源代码。该结果不应
写成“首次全量unit通过”。

| 证据 | 结果 | 日志 SHA256 |
| --- | --- | --- |
| RED / 新依赖 | 24 fail / 7 pass | `bfb9c9053418cb46b3f73353a9a30c8b78b2e67ed9f40e46d62c8624f9697b38` |
| RED / 原依赖 | 24 fail / 7 pass | `ad03f9e9c4f95990ddf670dcf7f5343efb2f434b0e54bdd84f94609bc964b60d` |
| GREEN / 新依赖 | 31 pass | `72f4a8d8ce433a1a0eb03dc76a32c05d185fff0cd124082215be90d24ab41520` |
| GREEN / 原依赖 | 31 pass | `c69b896f353402a8b9fabef9b17077348e2ce2d1dde324c6c528d3b835580c15` |
| 完整 unit 捕获 | 57 pass / 51 无测试 / 1 fail 包 | `606d51e1c33b9d71f062f66263d23b0f23fe06502064534a70ff894242b84fbe` |
| schema 可写复验 | 1 pass 包 | `376747d62b0258d1c391fe9d054c93d97e07343be5f6afd9a5ab977995a8680d` |
| service integration | 1 pass 包 | `20ce73acc65b7a74543272e8b71205023c842ace7b96637bcc49ecdfbb3749ce` |

19:51 CST 本机 Docker 查询恢复，但本任务旧容器、卷和 c9 候选镜像在当前 daemon
均不可见。本线未重启 Docker/WSL 或执行全局 prune；来源未确定。源码、归档和
历史报告仍在，新的候选需保留远端镜像归档，不依赖本机旧镜像可继续发布。

---
type: fix
scope: runtime
audience: developer
summary: 原生远程压缩按账号能力与实际上游协议调度，避免静默进入普通聊天转换
breaking: false
demo_ready: false
tests:
  - "PASS: go test -tags unit ./internal/service ./internal/handler -run 'Test(SchedulerRemoteCompactionCapability|SupportsRemoteCompactionV2|NativeCompaction|ForwardRejectsUnsupportedNativeCompaction|OpenAIResponsesRequiredCapability|OpenAIResponsesCompactionRoutingFlags|NormalizeOpenAIResponsesCompactRequest)' -count=1"
  - "PASS: go test -tags unit ./internal/service ./internal/handler -run '^TestNativeCompaction(ColdCatalogRefusalCanSelectAnotherAccount|FailoverExhaustedReportsCapability)$' -count=1"
  - "PASS: go test -tags unit ./internal/service ./internal/handler -count=1 (service 196.307s; handler 40.303s; after separately documented fixture fixes)"
  - "PASS: git diff --check"
artifacts:
  - backend/internal/service/provider_profile.go
  - backend/internal/service/account.go
  - backend/internal/service/openai_gateway_scheduling.go
  - backend/internal/service/openai_account_scheduler.go
  - backend/internal/service/model_protocol_catalog.go
  - backend/internal/service/openai_gateway_forward.go
  - backend/internal/service/openai_compact_body_signal.go
  - backend/internal/service/gateway_service.go
  - backend/internal/handler/openai_gateway_handler.go
  - backend/internal/service/openai_remote_compaction_capability_test.go
  - backend/internal/service/openai_remote_compaction_catalog_test.go
  - docs/REMOTE_COMPACTION.md
---

## What changed

原生 `stream: true` + `compaction_trigger` 使用独立的请求能力，复用现有 OpenAI Compact override/probe、provider 默认端点声明、账号协议、模型映射和模型目录。DeepSeek 默认端点不参与原生压缩；自定义端点保持 unknown。普通 Responses 和 legacy Compact 保留既有行为。原生 v2 的已确认支持优先排序与 legacy 专属模型映射分离。

候选过滤只读有效目录缓存，不触发候选数量倍增的网络等待；过期负向记录不永久阻止重新观测。Forward 以一次实际协议判断决定转发，不将原生压缩转为普通 Chat Completions/Messages。无法处理的账号进入既有的有上限 failover，不污染账号健康分数；耗尽后返回明确的压缩能力错误，并覆盖已输出 keepalive 的 SSE 情形。

## Why it matters

普通协议兼容性不等于原生压缩能力。入口、选号和转发使用同一能力边界，保留分组/模型的归属，不用普通摘要伪造压缩输出，也不将 DeepSeek 模型静默切换为 GPT。设计与官方能力依据见 [原生远程压缩](../../REMOTE_COMPACTION.md)。

已完成的本轮工作：恢复重启前 WIP；核对 provider/账号/分组/池；实现能力与协议门控；覆盖三类调度、粘性选择、配置限制、模型目录、failover/SSE；完成独立只读审查并修正其发现的排序、目录等待、双重协议解析和错误归因问题。

## Demo posture / limitations

仅完成本地实现和审查，未提交、未推送、未发布、未改线上资源。没有新增数据库字段、认证策略或前端 Compact UI。这次不代表 DeepSeek 获得了原生 Compact，也不代表 unknown 上游已验证可用。跨 provider 管理/探测以及 agent 本地摘要降级不在本轮范围内。账号更新仍按原有快照/outbox 传播。

验证结果（2026-10-18）：完整 service/handler 单元测试最终通过（196.307s / 40.303s）。验证中发现的三处夹具问题已单独修正并审查，见 [测试修正条目](2026-10-18-test-service-fixtures.md)，没有增加 Compact 生产补丁。

以下保留首次验证及基线对照证据，未将此前失败改写为通过：

- 定向能力、协议、三类调度、粘性、冷/过期目录、转发和 SSE 回归通过；两项补充集成测试也通过。
- 首次完整 handler 包通过（40.298s）。完整 service 包运行 207.191s，有两项失败，当时不能标为整包全绿。
- `TestOllamaProbeCallback_StaleLongDoesNotOverrideNewShort`：当前代码单跑仍失败；通过 Go overlay 将全部本轮修改的 Go 文件还原为 `HEAD 83d12d1a` 并屏蔽新增测试后，同一断言仍失败，确认是现有基线问题。
- `TestInflightEstimate_AccountMappingNoDBAndBoundedMemory`：首次整包时测得 8,399,184 字节，超过 8,388,608 字节阈值；当前版本和 HEAD overlay 单跑均通过。随后 profile 定位到测试夹具未初始化倍率解析器造成的重复 janitor/ticker 分配，修正夹具后保留原阈值并通过整包；计费生产代码未改。

基线对照命令：`go test -overlay=<temporary HEAD overlay.json> -tags unit ./internal/service -run '^(TestInflightEstimate_AccountMappingNoDBAndBoundedMemory|TestOllamaProbeCallback_StaleLongDoesNotOverrideNewShort)$' -count=1`。overlay 仅影响测试编译，没有回写源码或更改工作树。

独立审查创建参数为 `gpt-6-astra/max`；工具没有暴露实际运行时型号字段。最终生产 diff 与后续三处测试修正均未发现剩余阻断问题。本地验证已收口，不据此宣称生产验收通过。

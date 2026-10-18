---
type: test
scope: runtime
audience: developer
summary: 修正估价测试夹具初始化及两个依赖时钟精度的测试
breaking: false
demo_ready: false
tests:
  - "PASS: go test -tags unit ./internal/service -run '^(TestInflightEstimate_|TestOllamaProbeCallback_|TestHandle429_Ollama|TestCommandCodeGatewayPassesThroughCatalogProtocols)' -count=5"
  - "PASS: go test -tags unit ./internal/service ./internal/handler -count=1 (service 196.307s; handler 40.303s)"
  - "PASS: git diff --check"
artifacts:
  - backend/internal/service/billing_inflight_reservation_test.go
  - backend/internal/service/ratelimit_service_ollama_429_test.go
  - backend/internal/service/model_protocol_catalog_test.go
---

## What changed

估价测试夹具按生产构造器的方式初始化并复用倍率解析器。Ollama 过期回调用例从已记录的触发时间生成明确不同的新冷却时间。模型目录直通用例用局部序号隔离缓存键，结束后清理自己的条目。

## Why it matters

Compact 验证暴露了原有测试不稳定点。估价夹具此前每次请求都会走未初始化解析器的兜底，创建新的 go-cache 清理协程；整包诊断的堆 profile 指向这些协程、ticker 和 janitor，而非按模型名的业务缓存。生产构造器已有解析器初始化，本次不修改计费实现，原有 8 MiB 阈值、两万次调用及无数据库查询断言全部保留。

Windows 上连续 `time.Now()` 可返回同一值。诊断确认 Ollama 用例声称“新一代”的 reset 与旧 reset 完全相等，不能期待 CAS 将其区分。目录测试也观测到重复时间戳；其失败缓存写入按设计保留旧目录，因此同键不能模拟独立的“目录不可用”场景。改用明确的状态差异和独立缓存键，保留原有失败路径断言。

## Demo posture / limitations

仅修改测试夹具；不改变 Ollama CAS、计费或目录缓存的生产语义。这次不代表生产问题已通过上线验证。该测试修正与 Compact 功能改动分开记录，可独立审阅和提交。定向回归连续五轮通过；完整 service/handler 单元测试通过（196.307s / 40.303s）。

原 Astra reviewer 已完成这三处变更的独立只读复核，未发现阻断问题。限制保留：`HeapAlloc` 是进程级统计；Ollama 用例验证不同 reset 的新旧代际，不证明时间字段完全相同的两次更新可被区分。生产既有 CAS 身份语义未扩大。

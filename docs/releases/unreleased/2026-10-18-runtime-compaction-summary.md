---
type: feature
scope: runtime
audience: developer
summary: 复用既有账号与转发链路提供显式启用的摘要压缩、持久载荷还原和按实际目标的原生能力探测
breaking: false
demo_ready: false
tests:
  - "PASS on backend snapshot b8lprjgu: go test -tags=unit ./internal/service ./internal/handler ./internal/repository -count=1 (300.037s / 55.297s / 18.136s)"
  - "PASS: frontend scoped Vitest 197 cases, vue-tsc --noEmit, i18n completeness 3 cases, targeted ESLint"
  - "PASS: isolated fake-upstream generation lab 15 cases; latest default scheduling lab 9 cases"
  - "PASS: final advanced scheduling error-class regression, provider model-picker correction, and wildcard/error-redaction boundaries"
  - "PASS: immutable candidate build with embedded UI and isolated fake-upstream acceptance"
  - "INCOMPLETE: production release gate and real external client acceptance"
artifacts:
  - backend/internal/service/openai_compaction_summary_codec.go
  - backend/internal/service/openai_compaction_summary_store.go
  - backend/internal/service/openai_compaction_summary_request.go
  - backend/internal/service/openai_compaction_summary_response.go
  - backend/internal/service/openai_compaction_probe_target.go
  - backend/internal/repository/security_secret_store.go
  - backend/internal/repository/account_repo.go
  - frontend/src/components/account/CompactCapabilitySettings.vue
  - docs/REMOTE_COMPACTION.md
---

## What changed

原生压缩判断绑定实际模型、协议、端点及配置身份，v2 与 legacy 分别记录证据。复用账号 extra、repository/outbox、默认及高级调度和现有模型目录；WebUI 在创建、编辑、批量和测试入口区分原生策略与摘要兼容。

摘要兼容默认关闭。显式开启后，通过选中账号的普通生成链路执行一次摘要，封装为可回传的自有压缩项。服务端使用独立持久密钥验证用户/组归属，并在后续协议转换前还原。原生上游密文保留；不再删除唯一 compaction 历史以换取一次表面成功的重试。

## Why it matters

普通 Responses 兼容不等于原生压缩能力。目标证据避免把同账号一个模型的结果扩展到全部模型；网关摘要让不具备原生能力的账号在显式启用后完成压缩和续接，不要求客户端实现额外摘要协议。复用既有路由、鉴权和用量链路，不增加第二套 provider 调度系统或数据库 schema。

## Demo posture / limitations

当前仍未上线。摘要有信息损失与一次推理成本；仅支持文本历史，最终模型须有已知上下文窗口并预留输出与安全余量。不能恢复第三方不可解密密文或未提供的隐藏历史。假上游验收不证明真实模型质量、全部 OAuth/厂商路径或真实客户端兼容。最终冻结候选已完成后端/前端构建、服务与 admin 回归、三协议隔离验收、高级调度验收和真实浏览器 UI/probe 验收；这些证据不替代真实外部客户端或生产验收。详细输入和密钥恢复限制见 [功能说明](../../REMOTE_COMPACTION.md)。

发布须先具备读取能力并保持签发关闭，再启用指定账号。启用后不能回滚到不认识自有载荷的 compactcap1；须保留读取能力和摘要密钥备份恢复链。当前不支持在线换钥，禁止覆盖密钥记录后滚动重启。生产资源、共享审计与最终客户端验收继续由发布工作单约束，不因本条目写入而放行。

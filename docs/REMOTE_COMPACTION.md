# 原生远程压缩能力

`/responses` 可用并不代表它支持原生 remote compaction v2。原生 v2 请求在普通 Responses 路径上携带 `stream: true` 和 `compaction_trigger`，客户端需要上游返回真正的压缩项。普通文本摘要不能直接替代这个契约。

## 现有配置的作用

| 配置/事实 | 作用 | 不代表什么 |
| --- | --- | --- |
| 分组、组合分组路由、账号模型映射 | 确定 provider、账号池和上游模型 | 不自动把 DeepSeek 的压缩请求改送 GPT |
| `api_protocol`、`api_base_urls`、`protocol_rules`、上游模型目录 | 选择实际上游协议和地址 | 普通 Responses 支持不等于压缩支持 |
| OpenAI 账号的 Compact 模式和探测 | 控制账号的压缩资格；显式模式优先于探测 | Force On 不会绕过普通账号能力限制或改变已配置的协议 |
| `compact_model_mapping` | 用于 legacy `/responses/compact` | 不应用于原生 v2 的初始选号或映射 |
| provider profile 的默认能力 | 描述该接入模式的默认 Responses 端点 | 不自动套用到自定义中转地址 |

当前 DeepSeek 管理界面没有 Compact 开关。不要要求操作人员在该页面设置不存在的 Force Off。这里通过后端已有 provider profile 的默认端点声明能力，不增加账号字段、数据库迁移或新管理界面。

## 判断边界

- 原生 v2 使用独立请求能力，普通 Responses、图片意图和 legacy Compact 保留各自判断。
- OpenAI 继续读取已有 Compact 模式和账号级探测，不新增另一套互相冲突的设置。
- OpenAI 原生 v2 保留“已探测支持优先于 unknown”的排序；该排序与 legacy 专属模型映射分别控制。
- DeepSeek 默认 Responses 端点标为不支持。依据 2026-10-18 核验的 [官方接口文档](https://api-docs.deepseek.com/api/create-response/) 和 [协议兼容表](https://api-docs.deepseek.com/guides/responses_api/)：`compaction_trigger` 不在支持的输入类型中，其他类型被忽略。该事实需要随上游能力更新复核。
- 自定义端点以及其他尚无证据的平台保持 unknown。unknown 保留尝试资格，不能展示为已确认支持。
- 账号仍须满足原有模型归属、权限、可调度性等条件。原生 v2 不能通过 Chat Completions 或 Messages 转换路径发送；普通请求仍可转换。
- 选号只读取未过期的模型目录缓存，不为各个候选发起网络探测。目录缺失或过期时保留 unknown；选中账号后，转发按实际解析的同一个协议决定是否可发送。
- Grok 的专用转发路径和 legacy 摘要适配保留。其原生 v2 路径保持既有行为，不在本次宣称已验证支持。

## 验证与限制

新增回归覆盖默认/自定义端点、adaptive 的分协议地址、OpenAI 手动配置、账号能力限制、两种调度器及 load-batch、粘性账号、模型映射、模型目录和转发前的拒绝。所有上游交互使用本地测试替身，不调用生产模型。

没有可用候选时，响应会说明正在请求原生压缩，并提示核对能力、协议和账号可用性；不会把限流或容量不足一概宣称为“不支持”。转发时若发现所选账号实际不兼容，会在转换或发送推理请求前进入既有的有上限账号切换；不重试同账号，不计为账号健康故障。切换耗尽时返回 `compact_not_supported`。模型目录读取可发生在转发阶段。

本次没有提供跨 provider 的持久化能力探测、管理 UI、按模型存储的 probe，或 agent 本地摘要降级。现有 OpenAI probe 仍是账号级，不能证明该账号所有模型都支持相同能力。unknown 上游仍可能忽略触发器；其返回结果不能被当成已成功压缩。本地通过不等于生产验证或已经上线。

账号设置更新按既有快照同步/outbox 传播；正在执行的选号不保证读后写一致，本次不新增强一致性机制。

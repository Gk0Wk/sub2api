# 压缩能力补全：复用现有链路的设计审阅稿

状态：2026-10-18 用户在确认原生密文隔离、认证不变和兼容风险后回复“对的对的，那我们就按照最佳的实现方式做”，批准按下述设计实施。当前进入实现与验证阶段；尚未部署。

本轮唯一源码写入者：本机原生 Chat UI `codex_cc3ffbdc_4e54_4c8e_b5fb_e9c9d9bef063`，provider `01a14db8-01b6-7c50-b491-e910381c1c74`，实际 `deepseek-flash / max` 已核对并读取接手回报。其直属上级是本 Sub2API leader `codex_83df68ba_e273_4c24_8b38_432152b751ff`，leader 向用户报告；deployment 主协调者保持平级协作和共享发布事项收口。leader 维护本 brief/TODO、审查与验收；执行者维护相关源码/测试/功能文档，避免同文件并行写入。原生团队 Astra reviewer 只读审查。检查点为功能单元完成或最长十分钟，回报在原 Chat UI，由 leader 主动读取。

## 目标与范围

客户端继续按已有 Responses 压缩协议请求并回传压缩项；Sub2API 对支持的实际目标使用上游原生压缩，对缺少原生能力的目标生成摘要、封装压缩项，并在后续请求中还原。客户端本地摘要仅是当前恢复措施，不能作为本功能完成的验收条件。管理员在现有账号/测试/分组 WebUI 中配置并理解压缩能力，使调度、真实转发与测试使用同一组账号、模型、协议和端点事实。已有分组、账号类型、聚合平台、模型映射、限流、粘性会话、计费和 compat 继续拥有各自职责。不新建压缩服务、路由系统、独立账号池或跨 provider 自动迁移。

验收以“配置保存读回 → 测试 → 调度 → 转发 → 下一轮续接”的完整链条为准，不能只以 HTTP 200、一个 compaction 标签或 WebUI 支持徽章为准。无上游推理调用的确定性测试先行；真实上游探测、部署和 Secret 操作保持各自授权范围。

本轮关键词：native remote compaction v2、legacy responses/compact、ProviderProfile、account extra、composite group、channel mapping、resolveUpstreamProtocol、apicompat、WebUI account test。

## 已核实的现状

源码事实源：当前 `fix/native-compaction-capability` 分支，HEAD `9c5daabd`，前置 `83d12d1a` 为 encrypted input failover 隔离，`b8fe61f7` 为独立测试夹具调整。上线后客户端本地压缩恢复成功不证明上述提交应整体进入同一 PR。

| 已有机制 | 事实源 | 本次复用与缺口 |
| --- | --- | --- |
| 分组与组合路由 | `docs/COMPOSITE_GROUPS.md`、`backend/internal/handler/composite_platform.go`、现有 composite route service | 先决定所属 provider 和 public/upstream model，压缩资格只在授权池内筛选；现有 responses 路由不区分普通推理与 v2 compact 意图 |
| 渠道 | `backend/internal/service/channel_service.go` | 继续负责映射、定价、映射链展示；不作为新的压缩状态 owner |
| 账号类型与接入模式 | `account.go`、`provider_profile.go` | OAuth/API key/影子账号及按量/套餐模式保持原认证逻辑；官方默认端点能力不能覆盖自定义中转 |
| 协议与聚合平台 | `upstream_protocol_routing.go`、`model_protocol_catalog.go` | 已有显式协议、protocol_rules、模型目录、provider 默认顺序；OpenCode 同账号不同模型可以选择不同协议 |
| 原生 Compact 配置 | `account.go: GetOpenAICompactMode/OpenAICompactSupportKnown` | auto/force_on/force_off 及 probe 仅在 OpenAI 读取；其他 provider 没有同等操作入口 |
| 账号测试 | `account_test_service.go`、`openai_compact_probe.go` | 已探测 v2；adaptive 与 by_model 测试入口会走独立路径，mode 未全部传到 Compact 探测；探针自己构建 URL/认证，与真实协议选择存在漂移风险 |
| probe 持久化 | `accounts.extra` 的 `openai_compact_*` | 当前是账号级布尔值，没有保存被测模型/最终端点/协议适用范围；一个模型通过不能代表整条聚合链全部通过 |
| compat | `backend/internal/pkg/apicompat/` | 复用普通请求/结果协议转换；不能从缺失内容中制造可被上游识别的原生加密状态 |
| Grok 摘要兼容 | `openai_gateway_grok_compact.go`、`openai_gateway_grok.go` | legacy compact 已有摘要请求和回放转换，但要求上游返回 reasoning.encrypted_content；不等于通用纯文本摘要封装 |
| WebUI | `CreateAccountModal.vue`、`EditAccountModal.vue`、`BulkEditAccountModal.vue`、`AccountTestModal.vue`、`AccountsView.vue` | Compact 入口按 OpenAI 显示；说明仍偏 legacy 路径，实际 probe 已变 v2；profile catalog 尚未暴露后端默认压缩能力 |

当前 probe 接受 SSE added/done、仅 completed.output 以及整体 JSON 中存在 compaction 类型的结果。该事实不充分证明真实客户端可以完成 v2 流式消费；需要基于现有 SSE 处理与实际客户端契约做一致性验证，不能再加独立宽松解析器。

## 推荐的最小实现顺序

### 1. 统一能力判定的事实来源

继续在现有 account/provider/protocol 层判定，不新增独立调度器。输入使用当前请求的最终账号、账号类型/模式、映射后模型、协议和最终端点。

优先级分两层：

- 不可覆盖约束：分组/模型权限、实际协议、账号认证与健康/并发/配额，以及加密状态归属；Force On 不跳过这些条件。
- 能力声明：管理员显式开关 → 与实际目标匹配的有效 probe → 当前默认端点的 provider 声明 → unknown。Force On 表示管理员允许尝试，不等于探测已验证。

auto 下 unknown 延续可尝试的兼容语义；已确认支持优先。失败分类必须区分“确定不支持”与“没有空闲/有效账号”，避免把限流永久记为能力不支持。

### 2. 扩展现有账号配置和 probe，避免账号级误判

复用现有模式字段及 admin account update/test 接口，保留旧 OpenAI 数据的读取与原语义；不做数据库迁移或批量改存量账号。

需要审阅的增量数据契约：在现有 extra 中增加有界的 Compact probe 目标记录，最少包含协议种类（v2/legacy 分开）、映射后模型、目标配置指纹、supported/unsupported/unknown、检查时间及错误类别。账号身份、凭据来源、端点/协议/映射等影响目标的变化使旧证据失效。指纹由现有实际目标构造信息计算，不保存或展示凭证；条目容量和失效规则需在测试中固定。

旧账号级布尔值只按其原有兼容范围读取，不复制成所有模型的探测结果。更新使用已有 repository extra 更新及快照/outbox；避免并发 probe 以整份 extra 回写覆盖其他设置。无法安全复用的字段先列为待裁定项，不并行维护新旧两套可写策略。

测试入口必须复用现有协议/端点/模型映射解析，仍走原账号类型的认证、代理、headers override。只提取被真实请求和 probe 同时需要的目标解析，不为测试复制整个 forwarding/billing 链路。

### 3. WebUI 接通既有能力

把现有 Compact 区域提取为共享账号组件，在已有创建、编辑、批量和测试页使用；显示是否能原生调用 Responses、当前策略、适用模型、探测时间、来源和失败原因。Chat/Messages-only 账号显示明确限制，不提供无效的原生探测入口。

后端 profile catalog 提供默认事实，WebUI 不维护另一份厂商能力白名单。DeepSeek 自定义 Responses 端点可测试；官方默认端点显示默认不支持且允许管理员显式复核。聚合账号测试必须选定模型，不能把某一个模型的成功显示为整个账号所有模型支持。

Compact 专用模型映射需要准确解释：当前首次选择专用模型仅在 legacy 路径；`openai_compact_fallback.go` 的模型/上下文错误重试也覆盖 v2，不能标成“只对 legacy 生效”。本次保留已有行为，对新摘要路径继续使用原账号和已经授权的模型映射，不因 Compact 把 DeepSeek 请求改投 GPT。账号/组合路由预览优先扩展现有接口展示最终路径及不满足原因；如必须增加请求意图字段，单独审阅其兼容契约。

### 4. 真实流式契约与失败处理

复用现有 SSE parser/response handling 验证完整终态及实际可消费的压缩项。覆盖纯文本 200、空 encrypted content、多项、重复 added/done、只在 completed 出现、截断和中途 error；不要因为类型字符串出现就报告成功。

在响应开始前可按既有有界 failover 选择同授权池内其他账号；开始透传后不能承诺透明重试或二次计费安全。native encrypted state 的归属继续遵守现有 failover 隔离，不能为了压缩成功跨账号/厂商乱转。

## 已确认的最终目标与服务端闭环方案

用户原话：“最终肯定是希望，sub2api自己就能做到，而不是需要客户端配合”。因此交付范围包括服务端摘要及后续还原，不能停在能力检测和报错改善。

### 请求与返回

1. 请求经过现有鉴权、分组授权和大小限制。识别 Sub2API 自有压缩项并验证归属，在协议转换或原生 encrypted-content 清理前还原；保留原生上游载荷的原有归属规则。
2. 继续使用组合分组、渠道映射、账号池和 `resolveUpstreamProtocol` 得到实际目标。在同一调度体系中选择执行方式（native 或 summary），不能先把所有非原生账号过滤掉再在外围补救。
3. 原生能力满足时走原链路。摘要模式将明确的压缩请求转换为普通的无工具摘要推理，去掉触发压缩的标记，避免递归；沿用该账号认证、代理、协议转换、并发、配额和 usage 归集。不能 HTTP 自调用网关，不能复制一套 provider 转发实现。
4. 摘要推理必须正常完成且得到非空、未截断的结果；输出长度有界。失败不能封装成成功压缩。服务端封装一个可还原的 compaction 项；v2 通过既有 SSE 输出机制返回 item 和完整终态，legacy 保持对应 JSON 契约。
5. 下一轮还原出的摘要以普通会话上下文进入已有 compat 转换，不升级为 system/developer 指令；继续保留本轮用户输入和原有高优先级指令。重新压缩必须包含前次摘要，不能逐次丢掉早期任务状态。

这是一种由网关实现的摘要压缩，不能称为获得了上游私有原生压缩能力，也不承诺无损。普通客户端只需保留并回传已有 compaction 项；实际客户端能否完整保留该项必须通过端到端验证，不能只看 JSON 类型合法。

### 已批准的数据与配置边界

- 复用 `accounts.extra` 和现有账号更新接口。保留 `openai_compact_mode` 的原生许可语义；新增一个默认关闭的摘要兼容开关（建议 `compact_summary_enabled`），在现有 Compact 区域显示。原生 Force Off 不应静默等于“摘要已获准”；显式开启摘要后才有 native/summary 两种可选执行方式。WebUI 分别显示原生支持与摘要兼容，不能把后者写回 `openai_compact_supported=true`。
- 加入前述有界、按实际目标绑定的 probe 证据；无数据库 schema 迁移。原生模式 unknown 可以受控尝试；401、429、超时、5xx 不触发盲目摘要重试，不永久标成不支持。仅明确的不支持且尚无不可撤销输出时允许切换；每次实际发生的上游用量均须归集。
- 新载荷仍放在现有 `compaction.encrypted_content` 字段，采用独立的版本前缀与 key ID。载荷为经过认证加密的摘要、版本、签发时间和 owner 范围；不存完整会话，不建立服务端会话数据库。版本、key ID 和用途参与认证，拒绝篡改、超限、未知版本和错误归属。
- 建议绑定当前用户和授权分组，允许同范围内 API key 轮换和既有合法账号 failover；摘要载荷不自行授予模型或组权限。跨用户、跨组不能回放；改变所属组需明确迁移处理，不能静默放开。
- 复用 `security_secrets` 已有持久化及并发首次创建机制，增加独立用途的 32 字节随机密钥条目 `compaction_summary_key_v1`。使用 Go 标准库 AES-GCM；不复用 JWT/TOTP 密钥，不增加外部依赖或新密钥服务。不在日志中记录密钥、摘要或完整载荷。仅在服务端启用该能力时初始化；失败只阻止该能力启用，不能让普通请求或整个服务启动失败。
- 压缩载荷不依赖 Redis TTL 或单实例内存，以免会话隔夜或实例重启后失效；已有备份恢复链必须覆盖此密钥。保留 key ID 供后续轮换读取，不在本轮顺带开发管理平台；旧密钥不能在仍需读取相应会话时删除。

代码证据：`openai_live_attestation.go` 已有网关加密封装，但从 JWT 派生且用途不同，不能直接拿来保护长期会话；`security_secret_bootstrap.go` 已有跨实例一致的持久密钥初始化。复用的是这些实现模式和存储设施，不扩大现有认证密钥用途。

### 大上下文与无法恢复的历史

本轮采用单次摘要，输入须适合最终目标的上下文窗口，包含摘要提示词和输出余量；超窗明确失败并保留历史。这与原生 compact 仍要求输入能装入窗口的边界一致。原设计中的自动分段是 leader 主动扩展的策略，2026-10-18 审查发现现有单次请求计费按 request ID 去重，逐段计费或粗合并会破坏原定价不变量，因此收敛为单次调用，不为本功能新增多调用计费编排。复用已有模型上下文信息；信息不足时不能声称已保证预算，也不另建厂商常量表。不得为了成功静默截断历史，不能把不支持的图片/文件当成已读文本。仅有 previous_response_id 的隐藏历史，除非既有 replay 确实恢复全文，否则明确拒绝摘要。

Sub2API 无法从其他上游的不可解密原生状态恢复已经丢弃的原文。此类旧会话需要原来源继续读取，或者用户持有的完整历史；新摘要能力不能宣称自动修复所有旧密文。自有载荷验证失败时不得删除该项后继续回答。

### 独立审查发现的必要正修

`83d12d1a` 的入口/换账号 sanitize 与旧 invalid-encrypted/lineage 清理将 compaction 当作 reasoning 一并删除；缺少原文时会丢失唯一历史。本次必须直接收敛此行为：原生 compaction 保持原文，不能删除后重试成功；自有载荷先验证还原；reasoning 的已有兼容逻辑按其可重放边界保留。复用可信的 response/account 绑定，缺乏归属证据时拒绝不安全重放，不新增原生密文 owner 数据库。

恢复需覆盖 OpenAI 和 Anthropic 的 Responses handler、WS 首帧和后续帧、HTTP bridge、passthrough；单改一个入口不算完成。Forward 正文直接写 c.Writer，内部摘要须在现有转发/规范化响应中形成有界结果接缝，明确关闭 native gate 与自动工具注入，不能使用伪造 Context/ResponseRecorder 或递归 HTTP 自调用掩盖耦合。

### 发布与撤回

先发布只读支持并保持摘要签发关闭，验证各实例能读取相同载荷；再在指定账号启用签发。此后撤回目标必须保留读取能力，关闭签发即可止住新行为。不能把当前不识别新载荷的 `compactcap1` 镜像直接当成启用后的安全回滚版本。

只读阶段发布前须只读盘点实际账号的 `compact_summary_enabled`：若已有开启值，不能假定新版本默认关闭就能保证不签发；应先明确该账号的启用顺序和授权，再执行发布。该盘点不修改账号配置、不引入第二个全局签发开关。

验收使用无本地压缩覆盖的隔离客户端，覆盖首次压缩、普通续接、连续多次压缩、实例重启、多副本、合法账号切换和关闭签发后旧载荷续接。确认通过后才移除目前用于恢复会话的客户端覆盖。

### 实施授权边界

用户已批准本方案的独立用途密钥、owner 绑定、可持久回传载荷和新增 extra 字段。执行者可在源码和隔离测试中实现；实际生产启用继续按既有部署门禁、指定资源范围和回滚读取能力核验，不因方案批准跳过发布验收。若需扩大本设计的认证/安全边界，再报告具体差异。

## 验证计划与停止条件

- 单元/确定性仿真：官方/自定义端点、同名模型跨上游、同账号不同模型、映射变更、协议覆盖、probe 失效、并发更新、429/401/5xx、unsupported、unknown、Force On/Off。
- 调度链：分组/组合路由/渠道映射、粘性账号、两种调度器、load-batch、影子账号、failover、鉴权/配额与计费一次性不变量。
- 协议回归：普通 Responses/Chat/Messages 不受影响；native v2 与 legacy 不混用；压缩后下一轮可消费。
- 前端：Vitest + vue-tsc；真实浏览器检查创建/编辑保存重开、批量编辑、切换 provider/mode/protocol、测试事件和状态解释。静态代码 scout 不算 UI 验收。
- 整体：既有 service/handler 测试、隔离假上游集成；源码审查通过后再准备不可变镜像和独立发布验收。
- 停止：需要 schema/公共配置/认证/Secret 边界变更而未审阅；unknown 被静默标 supported；需跨 provider 兜底；重试产生重复推理/计费；同一问题连续 2–3 次验证失败；改动演变为第二套路由/状态系统。

## PR 收敛

先对上游基线审阅整个最终 diff。现有 `9c5daabd` 中被统一判定替代的分支直接收敛，不继续叠新 wrapper；`83d12d1a` 的 encrypted input 隔离单独审查与提交；`b8fe61f7` 的夹具修正独立。移除与本项无关的索引/忽略规则改动需先核对 owner，不能为清理 PR 覆盖其他 WIP。尚未授权改变已发布历史或强推。

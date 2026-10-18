# 待办

原生 Compact 能力门控的实现、完整 service/handler 单元测试与独立审查已完成，记录归档至 [功能条目](docs/releases/unreleased/2026-10-18-runtime-native-compaction-capability.md) 和 [测试夹具修正](docs/releases/unreleased/2026-10-18-test-service-fixtures.md)。

- [ ] 整理 PR：功能基线 `9c5daabd`、测试夹具修正 `b8fe61f7` 与本轮压缩候选仍待拆分后提交；当前仅完成最终 diff、候选构建和隔离验收，尚未 push/创建 PR。

## 压缩能力补全

客户端恢复已完成并归档到 deployment 工作单；三个原会话的本地 compact 和普通续接通过。用户要求不改 Orca，继续补全 Sub2API，并复用现有账号、分组、compat、聚合和 WebUI。设计与验证计划见 [brief](working-delta/2026-10-18-compaction-completion/brief.md)。

- [x] 核对原源码 checkout、WIP 和当前上线补丁；完成配置/测试/调度/compat/UI scout。
- [x] 用户确认最终由 Sub2API 完成摘要压缩和续接，无需客户端配置；已补充闭环设计与分阶段发布约束。
- [x] 用户批准按已审阅设计实施独立用途密钥、归属绑定及可持久回传载荷、账号 extra 增量契约。
- [x] 第一单元本地收口：codec/持久密钥/HTTP 与 WS 还原、摘要生成、原生密文保护；Astra 本轮已报问题已关闭，完整 service/handler/repository 测试及 15 项隔离 HTTP lab 通过。真实客户端与最终发布验收留在后续单元。
- [x] 实际目标、按目标的 probe 证据与调度统一：目标记录、DB fresh-unsupported 诊断、wildcard 映射边界和错误脱敏已收口；执行证据见 Orca `run_c7f9206b3a94` / `task_abf979a9fdf2` / `ctx_87315e0d2854`。
- [x] 接通现有 WebUI 创建/编辑/批量/测试入口；隔离真实浏览器已验证 DeepSeek 模型下拉、原生 v2 probe、固定 Chat 禁用提示、OpenAI 映射保留和无 console/HTTP 错误。
- [x] 在既有转发/compat 中实现摘要生成、封装和下一轮还原，覆盖超限输入、连续压缩、三协议、legacy、跨实例/重启读取、关闭签发继续读取及原生密文保留。
- [x] 完整协议、调度、计费及压缩后续接回归；Astra 独立审查接受最终五文件冻结版。
- [ ] 按最终范围准备发布，不改当前恢复中的客户端策略或扩大 live 资源范围。

独立隔离验收入口已准备：`C:/Users/nmg_w/AppData/Local/Temp/sub2api-summary-compat-lab.py --image <candidate>`，只用假上游和新建 tmpfs 数据库，已运行后端诊断快照，尚非完整发布候选。覆盖原生保持、摘要默认关闭、Responses/Chat/Messages、首次与重复压缩、legacy、归属/篡改拒绝、一次计费、key 轮换、关闭签发继续读、第二实例与重启。部署共享 audit 仍未解除；先完成可审查产物与恢复验收，再按具体发布范围处理。

generation 诊断检查点：固定后端快照 `C:/Users/nmg_w/AppData/Local/Temp/sub2api-summary-generation-snapshot-gz0g0rh8` 构建成功，诊断镜像 `sub2api-summary-generation:lab1` 尚非发布候选。原生压缩/续接与摘要关闭无出站通过；摘要 Responses 往返虽产生了可还原载荷，但被旧 compact fallback 错改为 gpt-5.5，legacy 仍请求上游 compact 端点返回 400；强制 Chat/Messages 账号返回 503、无出站。证据 `C:/Users/nmg_w/AppData/Local/Temp/sub2api-summary-lab-npza0az1`、`sub2api-summary-lab-fasu_r5g`。503 原始错误已在 `sub2api-summary-lab-jymr7a72/app.log` 核实为 pool=1/capability_mismatch；组绑定正确，根因是 `repository/scheduler_cache.go::filterSchedulerExtra` 投影遗漏摘要开关与原生 Compact 选择字段。Astra 同轮还确认多协议 SSE 回程、跨 attempt 状态残留、Antigravity 失败用量丢失和上游提前 EOF 缺少终态处理。暂停新增 UI/probe，先由原执行者修复真实生成链路并补回归，再继续后续单元；此结果不代表候选通过或已上线。

2026-10-18 generation 第二检查点：固定快照 `sub2api-summary-generation-snapshot-qpzjojp9` 的诊断镜像 `sha256:8851543d36c3225bde5a21d18ddd7e25661edb93031e90aa60bf562cbafa1048` 已通过 14 项隔离端到端验收，证据 `C:/Users/nmg_w/AppData/Local/Temp/sub2api-summary-lab-5euam89_`。覆盖三协议摘要/续接/再压缩/legacy、归属与篡改拒绝、key轮换、授权账号切换、关闭签发继续读、失败一次推理一次记账、第二实例与重启。仅证明假上游 API key 路径，不覆盖 OAuth 透传、Vertex/Antigravity 最终模型及真实客户端。Astra 同快照仍发现 3 个 P1：预算早于最终模型解析；legacy 透传仍套用 compact 模型映射和 OAuth compact 参数变换；OAuth 透传因上游 stream=true 绕过摘要封装。已交原唯一 writer 在原单元修复及真实链路回归；暂不进入 UI/probe/docs 或部署。

2026-10-18 generation 第三检查点：固定快照 `sub2api-summary-generation-snapshot-ccdjvtpa`、诊断镜像 `sha256:7a807c6c0a59f8ffbcd57540c404c7654bcc643687af3e4e6aa871ac2d338449` 的 15 项隔离验收全部通过，新增透传 v2/legacy、普通及 compact 映射存在时的实际目标不变断言；证据 `C:/Users/nmg_w/AppData/Local/Temp/sub2api-summary-lab-cjp2e5ps`。Astra 对同快照复核确认上述 3 个 P1 closed。执行者窄测试 service 17.705s、handler 0.432s 通过；leader 完整 service/handler/repository suite 正在运行，尚无完成结论。当前前端 3 个文件共 129 tests 和 vue-tsc --noEmit 通过；本机 pnpm11 自动安装产生的锁文件及 workspace 改动已撤回，依赖以 CI 使用的 pnpm9.15.9 按 frozen lockfile 恢复，未升级依赖。

2026-10-18 第一单元验收补充：首轮 full service 发现两个 WS 遗留项，原 writer 已修复且 Astra 静态复核 closed；原生 compaction 存在时不得清 reasoning/previous ID 后猜测重连，duplicate type fixture 改为严格拒绝及零出站断言。最终快照 `sub2api-summary-generation-snapshot-vvpenzu_` 构建通过，`go test -tags=unit ./internal/service -count=1` 全套 210.412s 通过，日志同快照 `full-service-unit-test.log`。相较 ccdjvtpa 只变更三个 WS service 文件；原全套 handler 41.128s 通过，repository 补齐 Windows Git sh PATH 后全套 9.738s 通过（未改备份代码）。HTTP lab 15 项证据继续适用于未改变的 HTTP 路径；不代表 OAuth/真实 provider/真实客户端/最终 UI 或生产已经验收。

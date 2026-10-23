---
type: fix
scope: runtime
audience: internal
summary: 吸收官方安全修复并移除未使用的图标依赖链，保持构建产物一致
breaking: false
demo_ready: false
tests:
  - "pnpm 9.15.9 --dir frontend install --frozen-lockfile: passed"
  - "pnpm 9.15.9 --dir frontend audit --prod --json: 0 critical, 2 high, 0 moderate, 0 low"
  - "python tools/check_pnpm_audit_exceptions.py --audit .tmp/security-20261023/frontend-audit-final.json --exceptions .github/audit-exceptions.yml: passed"
  - "pnpm 9.15.9 --dir frontend run lint:check: passed"
  - "Vitest after official patches, before unused-dependency removal: 371 files / 2868 tests passed, two workers"
  - "pnpm 9.15.9 --dir frontend exec vue-tsc -b: passed after dependency removal"
  - "pnpm 9.15.9 --dir frontend exec vite build: passed after dependency removal"
  - "Build SHA256 comparison before/after unused-dependency removal: all 178 files identical"
  - "Lockfile package-block comparison: removed 532, added 0, retained package resolutions unchanged"
artifacts:
  - frontend/package.json
  - frontend/pnpm-lock.yaml
  - DEV_GUIDE.md
---

## What changed

从 Sub2API 官方上游提交 `cacc80b8620e11e9655ab307471f71eb9ee31a0e` 和
`8286268047cbbca4e7f5fbb63ff0890c1ea37cf5` 精确吸收两个依赖文件：DOMPurify
`3.4.14 -> 3.4.17`、follow-redirects `1.16.0 -> 1.16.1`。已有 pnpm override
同时覆盖 Mermaid 和 axios 使用的传递副本。没有合并其他上游业务改动。

随后移除未使用的 `@lobehub/icons`：已跟踪源码仅有两个内置 SVG 组件的来源注释，
没有导入该包。锁文件从 1,025 个 package entry 减到 493 个；没有新增包，保留包的
完整 package block 均未变化。未修改应用源码或 SVG。

## Why it matters

初始生产依赖审计有 1 Critical、5 High；更新后本轮 DOMPurify 与
follow-redirects 告警消失。移除未使用依赖链后，其 13 Moderate / 2 Low 告警也消失，
最终生产审计为 0 Critical / 2 High / 0 Moderate / 0 Low。
厂商依据为 [DOMPurify 3.4.17](https://github.com/cure53/DOMPurify/releases/tag/3.4.17)
及 [follow-redirects 修复公告](https://github.com/follow-redirects/follow-redirects/security/advisories/GHSA-8r9p-f939-6h3c)。

剩余两条 High 是 `xlsx` 的 GHSA-4r6h-8v6p-xvw6、GHSA-5pgg-2g8v-p4x9。
现有例外到期日为 `2027-01-06`；本轮重新检查到 `UsageView.vue` 只有动态导入和
`XLSX.write` 导出，无 `XLSX.read` 解析路径。没有新增、延期或放宽审计例外。
审计通过依据是已存在的限定例外，不表示生产依赖零漏洞。

## Demo posture / limitations

官方补丁后的全量前端测试通过（371 文件 / 2,868 测试）；删除未使用依赖后重新完成
frozen/offline install、lint、类型检查、生产构建和生产审计，178 个构建文件的路径及
SHA256 与删除前完全一致。执行使用 pnpm `9.15.9`（与 CI 的 major 相同）、两个
Vitest worker、Node heap 上限 2048 MiB；没有重复启动全量测试。

全依赖审计另外揭示开发/测试工具链尚有严重问题：合计 3 Critical / 36 High /
22 Moderate / 2 Low（含生产 xlsx 两项），涉及 Vite、Vitest、tinypool 等，须独立
采用官方兼容修复并验证。生产审计 Critical 为零不代表整个仓库 Critical 为零。
新 head 的远端 CI、内部生产 fork 和 registry/codex-relay 发布分别验收；此记录
仅证明本地源码候选及上述验证，不代表生产已采用前端补丁。

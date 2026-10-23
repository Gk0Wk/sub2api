---
type: fix
scope: ci
audience: internal
summary: 更新前端开发工具安全版本，并阻止审计服务错误被当成审计通过
breaking: false
demo_ready: false
tests:
  - "pnpm 9.15.9 frozen/offline install: passed"
  - "pnpm 9.15.9 lint:check and vue-tsc -b: passed"
  - "Vite 6.4.4 production build: passed"
  - "UsageView.spec.ts with Vitest 4.1.12: 11 tests passed after constructor mock migration"
  - "Full Vitest 4.1.12 suite after migration: 371 files / 2868 tests passed, one worker"
  - "python -m unittest discover -s tools -p test_check_pnpm_audit_exceptions.py -v: 5 passed"
  - "Official npm production audit: 0 critical / 2 high / 0 moderate / 0 low; existing exceptions validated"
  - "Full dependency audit: 0 critical / 3 high / 1 moderate / 0 low"
  - "Production dependency graph comparison with 458b77862: 116 snapshot entries, all unchanged"
artifacts:
  - frontend/package.json
  - frontend/pnpm-lock.yaml
  - frontend/src/views/user/__tests__/UsageView.spec.ts
  - tools/check_pnpm_audit_exceptions.py
  - tools/test_check_pnpm_audit_exceptions.py
  - .github/workflows/security-scan.yml
  - DEV_GUIDE.md
---

## What changed

采用官方 [Vite 6.4.4](https://github.com/vitejs/vite/releases/tag/v6.4.4)、
[Vitest 4.1.12](https://github.com/vitest-dev/vitest/releases/tag/v4.1.12)，同步
coverage-v8，并更新兼容范围内的开发传递依赖。editorconfig 精确锁住的旧
minimatch 9.0.1 通过限定在受影响 9.x 范围的 override 使用 9.0.9；没有跨主版本
强制覆盖剩余依赖。Vitest 4 使用自己的 pool，移除了旧 tinypool 依赖。

完整前端测试首次发现 UsageView 的四项导出失败；根因为三处 Blob mock 使用箭头
函数，而 Vitest 4 按真实构造器调用 mock。按[官方迁移说明](https://v4.vitest.dev/guide/migration.html#spyon-and-fn-support-constructors)
将这三处改为普通函数。该文件11项回归通过，应用代码、CSV 安全断言和 coverage
门槛均未修改。

复验还发现镜像 registry 返回 `ERR_PNPM_AUDIT_ENDPOINT_NOT_EXISTS` 时，旧
审计检查器会打印成功。新增先红后绿的回归测试，检查器现在拒绝错误响应、非对象或
缺失合法漏洞集合的报告；仍支持两种既有报告格式和期限例外。CI 的审计请求显式指向
npm 官方 registry，并运行该检查器的回归测试。

## Why it matters

相对运行依赖安全 PR 的开发依赖基线，全依赖审计从 3 Critical / 36 High /
22 Moderate / 2 Low 降到 0 Critical / 3 High / 1 Moderate / 0 Low。
这些是审计记录计数，不是已证实的可利用路径数量。所有生产直接依赖及其116个
可达 snapshot 均与 `458b77862` 一致；工具链更新没有漂移生产依赖图。

生产审计仍为两个 XLSX High，沿用现有仅导出使用的例外，到期日未改。开发依赖
另剩 braces 3.0.3 High（审计无修复版本）与 postcss-selector-parser 6.1.5
Moderate（修复要求7.1.6，尚未完成跨主版本兼容升级）。未添加或延长任何例外。

## Demo posture / limitations

这次不代表内部 codex-relay 已部署工具链候选，也不表示全仓或最终镜像零漏洞。
现网两文件依赖补丁在独立分支构建，本分支不混入该发布窗口。

本机使用官方哈希校验的独立 Node 24.22.1 与 pnpm 9.15.9；不修改全局 Node。
测试一个 worker、1GiB Node heap，类型检查/构建2GiB heap，重构建串行。
首次1GiB类型检查的内存不足与首次四项 mock 失败均保留证据；最终全量复验通过
371文件 / 2,868测试，用时679秒。
Vite 保留既有大 chunk 提示，没有降低阈值。覆盖率工具已同步版本，但本轮未运行
coverage 模式，不宣称覆盖率达标。

浏览器 API 返回空的 apps/browsers 清单，本轮未完成真实页面、console、network
或认证流程验证；组件测试和构建结果不能代替这部分验收。

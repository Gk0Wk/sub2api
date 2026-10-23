---
type: fix
scope: runtime
audience: internal
summary: 在现网 Sub2API 源码上形成仅两文件的前端依赖安全候选
breaking: false
demo_ready: false
tests:
  - "git -c core.whitespace=cr-at-eol diff --check 67e32e1f38cfa0f84036935ab95a6dace03def37 43c7637750ec23d33d3a6f272952202dd119aa02"
  - "pnpm frozen/offline install, lint, vue-tsc -b and production build: passed"
  - "Vitest 2 with one worker: 372 files / 2886 tests passed, zero failures"
  - "Official npm production audit: 0 Critical / 2 High / 0 Moderate / 0 Low; existing exceptions validated"
  - "docker buildx build: linux/amd64 image completed from exact source archive under 2 CPU / 4 GiB"
  - "Offline readonly UID1000 container: sub2api -version, pg_dump --version and psql --version passed"
  - "go version -m: go1.27.3, linux/amd64, CGO_ENABLED=0, embed build tag"
  - "Trivy 0.76.1 with preserved 2026-10-21 database: candidate and live baseline findings identical"
artifacts:
  - frontend/package.json
  - frontend/pnpm-lock.yaml
---

## What changed

从已核对的现网源码 `67e32e1f38cfa0f84036935ab95a6dace03def37` 创建独立候选，
仅修改 `frontend/package.json` 与 `frontend/pnpm-lock.yaml`：采用官方
DOMPurify 3.4.17、follow-redirects 1.16.1，移除应用没有导入的 `@lobehub/icons`。
两个文件与已审运行依赖 PR #1 的 head `458b77862074651f7cb1873955430ea882a9d08d`
逐字节一致；未合入开发工具链 PR #2 或其他业务差异。

镜像源码 commit 为 `43c7637750ec23d33d3a6f272952202dd119aa02`，tree 为
`577eeaf9920c0a0af3cf6faeafc20b2240a7d164`。本记录属于后续文档变更，不改变
镜像源码身份。源码归档 SHA256 为
`5664bc59820abd4219d16ed89017e41989d69658f1521e9b2be84ba7c99249a6`。

## Why it matters

官方运行依赖修复可以进入内部现网谱系，同时保留已有业务修复和部署边界。全量
前端测试、类型检查与178文件生产构建通过；有效的官方npm审计仅剩既有两项XLSX
High。两项例外到期日仍为2027-01-06，使用范围为导出，不新增或延期例外。

旧审计检查器曾把registry错误JSON判为通过，修复已在独立PR #2合入；本候选使用
该修复后的检查器在树外核验真实报告，未把检查器变更混进这两文件源码提交。

## Demo posture / limitations

linux/amd64镜像已构建，manifest为
`sha256:9403848fa11c757b186ccc2fced5c2d55cbecf816443a2220ada59524d790102`，
config为`sha256:81d7e872cb3a5aadd309d78a615e803887e6605fc3b66a00b157825670c7e126`。
尚未推送或部署，真实页面、认证和生产请求未验收。
构建使用精确源码的git archive，未带入宿主缓存、node_modules或认证配置。
前端Node 24.22.1和后端Go 1.27.3均沿用已固定的官方镜像。

第一次构建因Go官方层下载超过8分钟无进展而主动停止。第二次采用官方支持的
[OCI named context](https://docs.docker.com/reference/cli/docker/buildx/build/#use-an-oci-layout-directory-as-build-context)
从本机读取相同镜像：12个blob全部通过SHA256核验，官方index为
`sha256:d530cf60723abf1ffae0cd16bb546b379ec0550d2f3815e5b16ed69dbc00fba4`，
amd64 manifest为
`sha256:bff16a94c0b4b6d9642581e2a3aa81beb3edb940b2f5a3aac1d22c6a68cb7747`。
Dockerfile和编译版本未改；builder仍限制2CPU / 4GiB，未使用X1O资源。

最终二进制回读为Sub2API 0.2.17 / commit43c763775，Go buildinfo为1.27.3、
linux/amd64、CGO_ENABLED=0、embed tag。无网络、只读、UID1000、cap-drop ALL的
受限容器版本检查exit0，pg_dump与psql均为18.6；没有启动应用健康服务。
二进制SHA256为`896ba33e0ddf19d973c6a5cec765a209bac8fa85b02d37139ba6e86635b63747`。
验证容器已删除，builder已停止。

最新Trivy数据库经两个官方分发源下载均失败，已停止重试。本轮扫描仅使用
哈希核验后的2026-10-21T19:18:22Z基线库：Alpine32包零告警，Go143包为
0 Critical / 3 High / 3 Medium / 1 Unknown。与原live本地镜像`afa9cfea6ea4…`
同库重扫的finding集合完全一致，没有新增。它不代表最新漏洞覆盖或产品零风险。

三项High来自compress1.18.2及x/image0.41.0；官方修复下限分别为
[compress1.18.7](https://github.com/klauspost/compress/security/advisories/GHSA-259r-337f-4rfw)
及[x/image0.45.0](https://pkg.go.dev/vuln/GO-2026-6222)，TIFF另一High的
[下限为0.43.0](https://pkg.go.dev/vuln/GO-2026-5062)。Go官方二进制符号扫描亦报出
这两模块的已知条目，但不能据此推断外部请求可达性。
[GO-2026-5932](https://pkg.go.dev/vuln/GO-2026-5932)没有修复版本，保持Unknown。
后端残余进入单独的go.mod/go.sum修复单元，不覆盖本镜像身份。
原始验证、构建与数据库限制证据位于本工作树`.tmp/security-20261023/`。
生产协管负责独立复核和新的唯一执行窗口，候选完成不等于已上线。

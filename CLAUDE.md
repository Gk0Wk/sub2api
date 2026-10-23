# Sub2API

Go 后端、Vue 前端的 API 网关。入口与安装见 [README_CN.md](README_CN.md)。

- `backend/`：服务、存储及生成代码；`frontend/`：控制台；`deploy/`：部署模板。
- 后端测试：`cd backend && make test-unit`；集成测试：`make test-integration`。
- 版本跟随 `backend/go.mod`、前端 lockfile；Go 检查规则见 `backend/.golangci.yml`。
- OAuth 刷新边界与定向验证见 [service README](backend/internal/service/README.md)。
- 本地未决事项见 [TODO.md](TODO.md)，发布记录位于 `docs/releases/unreleased/`。

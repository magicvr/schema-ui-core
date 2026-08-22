---
id: E-001
doc: execution-entry
goal: GOAL-002-rotation-contract-freeze
status: recorded
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# E-001 · R1 配置面落地

## 事实（2026-08-22）

按 GOAL-002 D-001 实施，全部改动限于 `apps/api`：

| 改动 | 位置 |
|------|------|
| `Config.AuthJWTSecretPrevious` 字段（含合同注释） | `internal/config/config.go`（Config struct） |
| YAML 键 `auth.jwt_secret_previous` 解析（`strPtrOr`） | `internal/config/config.go`（yamlFile Auth 块 + Load 应用行） |
| env `AUTH_JWT_SECRET_PREVIOUS`（`envOr`，env 覆盖 YAML） | `internal/config/config.go`（Load env 段） |
| `ValidateProd`：previous 非空时 ≥32 字符 + 含字母与数字 + 与 current 不同值；错误只点名键名 | `internal/config/config.go`（非开发环境块） |
| 样例注记（默认注释、轮换操作顺序） | `internal/config/config.default.yaml`、`configs/config.yaml` |
| 测试 `TestJWTSecretPreviousConfig`（8 子用例） | `internal/config/config_test.go` |

## 验证证据

- `go build ./...`：通过（go1.26.0 windows/amd64）。
- `go test ./internal/config/`：`ok ... 1.221s`。
- `go test ./internal/config/ -run TestJWTSecretPreviousConfig -v`：8/8 PASS（合规通过 / 短 previous 拒绝 / 全数字 previous 拒绝 / 同值拒绝 / dev 低门槛 / 缺省单密钥 / YAML 解析 / env 覆盖）。
- `go test ./internal/config/ -run TestValidateProd -v`：既有 9 子用例全部 PASS（单密钥行为零变化）。
- 全仓 `go vet ./...` + `go test ./...`：见 E-002（关门审计前补记）。

## 对照检查点

- 检查点 1（字段+双通道+校验）：完成。
- 检查点 2（单测矩阵）：完成。
- 检查点 3（样例注记）：完成。
- 检查点 4（self 审计 + goal-tree）：进行中。

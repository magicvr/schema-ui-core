---
id: E-001
doc: execution-entry
goal: GOAL-003-dual-key-jwt
status: recorded
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# E-001 · R2 双密钥实现落地

## 事实（2026-08-22）

按 GOAL-003 D-001 实施：

| 改动 | 位置 |
|------|------|
| `Authenticator.previousSecret` 字段 + 合同注释 | `internal/auth/auth.go`（struct） |
| 新构造器 `NewWithRepositoryAndPrevious(current, previous, ...)`；既有 `New`/`NewWithRepository` 签名不变 | `internal/auth/auth.go` |
| `verifyAccess(raw)`：current 先验，失败且配置 previous 再试 previous；两次都强制过期/方法检查 | `internal/auth/auth.go`（Middleware 改调此方法） |
| composition 接线：`newAuthenticator` 读 `cfg.AuthJWTSecretPrevious` | `internal/composition/composition.go`（`NewApp` 对外签名不变） |
| 测试 `TestDualKeyRotationOverlapWindow`（4 子用例） | `internal/auth/auth_test.go` |

## 验证证据

- `go build ./...`：通过。
- `go test ./internal/auth/ -run TestDualKeyRotationOverlapWindow -v`：4/4 PASS——
  1. 重叠窗内旧 key token 通过中间件；
  2. 移除 previous 后同一 token 401；
  3. 双密钥下 Login 新 access 只被 current 验证、不被 old 验证；
  4. previous 签发的过期 token 在双密钥下仍 401（回退不延长寿命）。
- 全仓 vet + 全套件结果 → E-002。

## 对照检查点

- 检查点 1（I-003 决策）：done（D-001）。
- 检查点 2（实现 + 接线 + 单测）：done（本条）。
- 检查点 3（全仓验证 + E 记录）：E-002。
- 检查点 4（self → independent 审计 → 响应 → goal-tree）：进行中。

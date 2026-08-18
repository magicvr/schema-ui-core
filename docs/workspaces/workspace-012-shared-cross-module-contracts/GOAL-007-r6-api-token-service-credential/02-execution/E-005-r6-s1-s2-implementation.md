---
id: E-005-r6-s1-s2-implementation
goal: GOAL-007-r6-api-token-service-credential
status: recorded
created: 2026-08-19
updated: 2026-08-19
parent: GOAL-007-r6-api-token-service-credential
version: 0.1.0
---

# E-005 · R6 S1/S2 实施与全量验证

## 实施事实

1. `aa00f33` 增加 0044 `service_credentials`、0045 operation-log 三事件与 correlation-safe rebuild、hash-only repository、同事务 operation recorder，以及 raw 不落库、大小写唯一、审计回滚和迁移相关性测试。
2. `ce8d952` 增加 `sui_sc_` 256-bit prefixed secret、service principal 与 Bearer 分流；service prefix 在 JWT/dev fallback 前判定，unknown/expired/revoked 统一 401，`last_used_at` 与 use audit 为 best-effort。
3. `1864f49` 增加 human-only list/detail/create/revoke API、中心 `service-credentials.read/write` 权限、creator permission ceiling、reserved scope deny、一次性 secret、事务 create/revoke audit，以及 user-only `UserIdentityFrom` 隔离。
4. `2a2d0dd` 补强 concurrent duplicate/revoke 单次转换和 R5 maintenance/degraded/read-only 对 credential mutation 的 server-level 黑盒门禁。

## 验证事实

- `apps/api`: `go test ./...` 全量通过；handler 259.365s，auth/composition/authsession/operationlog/store/docscheck 等全部 exit 0。
- 定向：credential repository/API/auth、0045 migration、R5 operational gate、Profile/kernel/Manifest 与 system-data reconciliation 均通过；并发 duplicate/revoke 以 `-count=5` 复跑通过。
- `apps/web`: `npm run build` 通过（tsc + Vite，1859 modules）；只有既有 chunk-size warning。build 生成的三个 conformance claim 文件核对后恢复，协议资产无交付 diff。
- `apps/api`: `go test ./internal/docscheck -count=1` 通过；`git diff --check` 和 untracked whitespace 检查通过；验证结束后工作树干净。

## 当前门禁

S1/S2 实施事实完成，A-006 self verdict=pass，required=0。因 R6 属 security/data/migration/cross-module 边界，S3 仍须 grok-build independent 关门审计；未在本条提前关闭目标。

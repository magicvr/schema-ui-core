---
id: A-008-r6-a007-response
goal: GOAL-007-r6-api-token-service-credential
doc: audit-entry
record_id: A-008
source: self
scope: response to A-007 F-001 through F-005; R6 S3 remediation evidence
verdict: pass
status: recorded
parent: GOAL-007-r6-api-token-service-credential
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
responds_to: A-007
---

# A-008 · R6 A-007 response

## 审计头

| 项 | 值 |
|----|----|
| source | self |
| scope | A-007 F-001～F-005；提交 `b6ebfec`；R6 S3 finding remediation |
| verdict | pass |
| required findings | 0（实现侧判断；S3 仍待 independent 复核） |

## 逐项响应

| Finding | 原建议 | 响应 | 证据 |
|---------|--------|------|------|
| F-001 | required / med | **fixed**：create 201 一次性字段改为 D-002 冻结的 `secret`；未采用 residual 或 overrule | `service_credentials.go`；`service_credentials_test.go`；提交 `b6ebfec` |
| F-002 | recommended / low | **fixed**：重名 `fieldErrors.reason` 改为 `name already exists` | 同上 |
| F-003 | recommended / low | **fixed**：拒绝空集及超过 64 个 scopes | 同上；65-scope 黑盒用例 |
| F-004 | recommended / low | **fixed**：六类 user-only HTTP 表面逐项 401；补 expired credential 401 | `service_credentials_test.go`；`auth_test.go` |
| F-005 | recommended / low | **fixed**：use audit detail 增加 `scopeCount` | `auth.go`；`composition.go`；相关测试 |

## 验证

- 定向 auth、handler、composition 用例在整改后通过。
- `apps/api` 的 `go test ./...` 全量通过；新增管理 API、机器认证、用户态隔离、迁移、审计和 R5 gate 均包含在本轮包级回归中。
- `git diff --check` 在整改提交前通过，未发现未跟踪文件。

## 结论

A-007 唯一 required finding F-001 以及四条 recommended findings 已由可核对代码和测试修正。按 cross 审计路径，本条仅给出编排器 self 响应；GOAL-007 保持 `active`、progress=`75`，待 independent 对 A-007 finding closure 作出意见后再决定 S3 与关门。

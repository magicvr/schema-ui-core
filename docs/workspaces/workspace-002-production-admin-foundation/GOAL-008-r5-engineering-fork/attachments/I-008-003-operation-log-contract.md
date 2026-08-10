---
title: I-008-003 · 最小操作日志（operation_log）契约
status: active
doc_type: information-contract
created: 2026-08-03
updated: 2026-08-03
parent: GOAL-008-r5-engineering-fork
version: 1.0.0
related_info: I-008-003
related_decisions: D-009
---

# I-008-003 · 最小操作日志（operation_log）契约

## 1. 冻结范围

回答 `I-008-003` 的信息问题：R5 可选加分 S6「最小操作日志」的事件与存储契约是什么。依据 Root D-013 方案甲（SQLite `operation_log` 迁移 + repository，覆盖 records 写 + auth 关键事件）。

- **冻结内容**：事件类型全集、表结构与迁移（0004）、repository 写入/查询边界、接线点、失败语义、验证动作。
- **不在本协议内**：日志查询 API（HTTP 端点）、清理/轮转策略（本波次不交付，见 §7）、任何 UI 展示。本协议不构成 S6 已实施的证据。

## 2. 事件类型

| 事件 | 触发点 | 语义 |
|------|--------|------|
| `records.create` | POST /api/records 成功（201） | 新建记录 |
| `records.update` | PATCH /api/records/{id} 成功（200） | 更新记录字段 |
| `records.delete` | DELETE /api/records/{id} 成功（204） | 删除记录 |
| `auth.login` | POST /api/auth/login 成功（200） | 登录成功 |
| `auth.logout` | POST /api/auth/logout 成功（204） | 登出（撤销 refresh token） |
| `auth.refresh` | POST /api/auth/refresh 成功（200） | 刷新/轮换 token 对 |

仅成功操作记日志；被 400/401/403/404 拒绝或失败的操作不记录（最小范围，失败可由服务日志覆盖）。

## 3. 表结构与迁移

- 迁移版本：**0004**，名称 `operation_log`。
- DDL（在迁移事务内执行，进入 `schema_migrations` 台账）：

```sql
CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at INTEGER NOT NULL
);
CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC);
```

- `created_at` 为 Unix 毫秒（与 records `updated_at` 口径一致，D-004）。
- `record_id`：records 事件填写记录 id；auth 事件为 NULL。
- `detail`：可选 JSON 文本；records 事件存 `{"name":"<记录名>"}` 摘要，auth 事件存 `{"username":"<用户名>"}`；不含 token/password/secret。
- 既有数据库升级：在 0003 之后顺序应用 0004（`pre-v0004` 快照逻辑由既有迁移运行器自动承担）。

## 4. repository 边界

- `RecordOperation(Operation)`：追加一行（幂等由调用方生成唯一 id；失败返回 error）。
- `ListOperations(limit int)`：按 `created_at DESC, id DESC` 返回最近 N 条，供测试与运维核对；`limit <= 0` 返回空。
- **不提供**删除、清理、轮转 API（§7 非目标）。

## 5. 接线点（handler 层）

- `handler/records.go`：create/update/delete 成功写响应前记录（actor 取自 `requirePermission` 返回的 `account.User`）。
- `handler/auth.go`：login/refresh/logout 成功写响应前记录；logout 需 `auth.Authenticator.Logout` 返回被撤销 token 的 userID。
- 记录失败**不阻断主操作**（best-effort）：`RecordOperation` 返回 error 时以服务日志（slog.Error）留痕，HTTP 响应保持成功语义——审计日志不能引入新的失败面；这是加分项，非核心验收门禁。

## 6. 验证动作

- store 层：0004 迁移应用 + 台账；`RecordOperation`/`ListOperations` 往返与排序；CHECK 约束（非法事件拒绝）。
- handler 层：records create/update/delete 后日志含对应事件与 actor；auth login/logout/refresh 后日志含对应事件。
- 回归：`go test ./...`（apps/api）全绿。

## 7. 非目标（明确不交付）

- 日志 HTTP 查询/导出端点、管理 UI。
- 日志清理/轮转/保留策略（无界追加；运维侧容量管理不在本波次）。
- 失败操作（401/403/400/404）的日志。
- Web 端任何展示。

## 8. 修订记录

| 版本 | 日期 | 变更 |
|------|------|------|
| v1.0.0 | 2026-08-03 | 初始冻结（D-009）。 |

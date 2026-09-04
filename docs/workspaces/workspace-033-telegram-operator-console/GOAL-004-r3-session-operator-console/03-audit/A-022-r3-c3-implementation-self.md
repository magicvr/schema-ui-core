---
doc_type: goal-audit
id: A-022-r3-c3-implementation-self
parent: GOAL-004-r3-session-operator-console
date: 2026-09-05
source: self
auditor: Codex
audit_type: implementation
scope: C3 v69 outbound、operator API、RBAC/runtime 接线、会话/成绩单、幂等重试，以及 A-018 F-004～F-007 非阻断项
verdict: pass
open_required: 0
version: 0.1.0
---

# A-022 · R3 C3 实现自审（2026-09-05）

## 审视结论

在当前提交 `7ddc97e1` 上，C3 合同对应的生产实现与专项测试已落盘；A-018 的
required 实现门禁未发现新的开放项，`open_required: 0`。本条是 self opinion，
不替代用户要求的 Grok independent implementation audit，也不修改 C3/R3 的
status 或 progress；独立实现审计完成前不关闭 C3。

## 实现证据

- `apps/api/modules/channel/telegram/migration/migration.go` 增加 v69
  `telegram_outbound_messages` 及 SQLite/PostgreSQL DDL；`retry_root` 的 pending
  partial unique index、`(bot_id, request_id)` 主键、状态约束、索引和迁移
  checksum 已同步到 `apps/api/internal/store/identity.go` 及测试 fixture。
- `apps/api/modules/channel/telegram/store/repository.go` 提供会话/统一成绩单
  投影、pending/failed/sent 状态机和显式 retry。插入使用无 target 的
  `ON CONFLICT DO NOTHING`、`RowsAffected` 与同一事务内读取；request id 同时受
  bot/chat 范围和 `[A-Za-z0-9._-]{1,128}` 约束，避免跨 bot/chat 串读或 mux 歧义。
- `apps/api/internal/handler/telegram_operator.go` 提供四个 operator route，
  在写入 pending 后才调用现有 sender；终态重放不重复外发，sender 失败持久化
  `failed`，发送后的状态更新失败保持 pending 并 fail closed。未知 chat/request、
  in-progress、conflict、retry 禁止和 runtime unavailable 均使用稳定 catalog code；
  下游诊断只保留固定的非敏感类别。
- `apps/api/modules/channel/telegram/provider.go`、
  `apps/api/internal/composition/composition.go` 与 `apps/api/kernel/profile.go`
  已同步 operator routes/permissions、模块注册和 production composition；
  operator handler 在 composition 中经 `a.Middleware` 包装，runtime gate 检查
  running、bot id、webhook/polling receiver 及 `HasBusinessHandlers()`。
- `apps/api/internal/channel/telegram/lease_handler.go` 保留既有 lease session
  行为，并按 D-010 接受 `settings.read OR telegram.operator.read`；operator API
  不隐式启动 lease。`apps/api/internal/errorcatalog/errorcatalog.go` 与
  `error_contract_test.go` 已同步七个 Telegram operator catalog code。

## A-018 非阻断项处理

- F-004：Provider descriptor、kernel/profile contribution、migration/identity
  projection 已同步，routes 与 permissions 的数量和归属有测试覆盖。
- F-005：未知 chat/request 使用稳定 404 code；发送失败、并发冲突、重试状态和
  runtime 不可用均有固定 code，错误诊断不回显 token 或原始下游响应。
- F-006：request id 使用 mux-safe 字符集与长度上限，并有非法 id 测试。
- F-007：post-send 状态更新失败保持 pending；后续 replay 返回 in-progress，
  不会因内部状态不确定而再次外发。

## 验证事实

- 通过：`go test ./modules/channel/telegram/store ./internal/handler -run 'TestRepositoryOperatorProjectionAndOutboundStateSQLite|TestRepositoryOutboundRejectsUnsafeRequestIDSQLite|TestTelegramOperatorHandler' -count=1`。
- 通过：`go test ./internal/store ./internal/handler ./internal/composition ./internal/channel/telegram ./modules/channel/telegram/... ./kernel -count=1`。
- 通过：C3 handler 与 repository 的隔离 `-race` 专项测试；更大范围 `-race`
  仅在全量 handler 包的既有 wallet/SQLite 并发测试出现 database locked/internal
  failure，C3 专项 race 未复现，不能把该环境/基线波动写成 C3 通过证据。
- PostgreSQL 冲突/重试测试已登记为 gated test；当前无 PostgreSQL 环境时为 skip，
  不把 skip 计作 PostgreSQL 已验证通过。
- 提交前 `git diff --check` 与 owned/untracked 文件 trailing-whitespace 检查通过；
  实现提交为 `7ddc97e1 feat(telegram): add C3 operator messaging surface`。

## 后续门禁

将调用本地 `grok-4.6 · reasoning high` 对当前实现、迁移、测试和 F-004～F-007
进行 independent audit。若发现 required finding，必须按 P-003/P-004 修复或等待
用户书面裁决；在此之前 C3 仍为实现阶段，不能关闭。

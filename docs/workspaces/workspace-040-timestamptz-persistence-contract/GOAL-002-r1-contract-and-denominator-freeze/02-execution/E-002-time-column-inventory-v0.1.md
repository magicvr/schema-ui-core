---
id: E-002-time-column-inventory-v0.1
doc: execution-entry
status: recorded
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-002 · 时间列与单位 inventory v0.1

## 已发生事实

- 已对 apps/api compiled migration DDL 与代表性 runtime repositories 执行只读盘点。
- 已确认当前物理形状主要是 SQLite `INTEGER` / PostgreSQL `BIGINT`，不存在已发现的时间列 TEXT；运行时单位存在 seconds 与 milliseconds 并存。
- 已确认核心毫秒面：`core.jobs`、`core.operationlog`、`core.persistence.mail_outbox` 与 `mail_config`；其余主要模块目前按 Unix seconds 写读。
- 已识别特殊 sentinel/nullable 面：`locked_until`、`last_login_failure_at`、`mail_config.updated_at DEFAULT 0`、`telegram_config.updated_at DEFAULT 0`、nullable `read_at`/`restored_at`/`finished_at`/`expires_at` 等。
- 已写入 inventory 附件；本条不宣称 C1 完成，完整 catalog/运行时逐列核对仍在进行。

## 证据

- `attachments/r1-time-column-inventory-v0.1.md`
- `apps/api/modules/*/migration/migration.go`
- `apps/api/internal/jobs/model.go`
- `apps/api/internal/mail/outbox.go`
- `apps/api/internal/mail/runtime.go`
- `apps/api/modules/operationlog/repository.go`
- `apps/api/modules/operationlog/retention.go`

## 下一步

补齐 48-migration catalog 的机械核对、表/列级 SQLite/PG 当前 schema introspection 与每列 codec/转换映射；完成后进入 C2/C3 方案审视。

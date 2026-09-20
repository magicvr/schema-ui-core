---
id: A-008-r1-self-guardrails-readiness
doc_type: goal-audit-entry
source: self
auditor: /govern
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · C2/C3 guardrails draft / D-004 Backup SPI / D-005 wire / R2 entry
verdict: conditional
open_required: 7
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-008 · C2/C3 guardrails self-readiness

## 结论

`conditional`。C2/C3 guardrails 草案已把 A-006 的 required 门禁映射到 codec、NULL/default、追加-only catalog、CHECK/index/predicate、Backup SPI、wire inventory 与 R2 entry gate；但草案尚未冻结，且 Backup SPI 的最终 API owner/surface 仍需用户裁决。

## 已形成的草案证据

- `attachments/r1-c2-c3-guardrails-v0.1.md`：codec boundary、seconds/milliseconds conversion、NULL/sentinel、v73+ append-only、predicate/check/index、Backup SPI/Service、wire/R3 interface 与 R2 gate。
- `D-004-r2-migration-ownership`：按模块追加 migration，历史 v1–v72 immutable。
- `D-006-r1-backup-spi-user-decision`：SQLite/PG native provider + metadata/verification/restore-to-new-db，rollback 优先，完整 backup product 排除。
- `D-003/D-005`：6 位输出与 RFC3339 输入兼容。

## Required findings

### F-R1-002 · codec/DDL/精度/排序仍未冻结

Guardrail 已列出目标，但没有逐列 codec、PG `USING`、SQLite rebuild、microsecond truncation/rounding 与 canonical ordering 的 accepted contract。

### F-R1-003 · NULL/zero/default 与 predicate mapping 仍未冻结

Guardrail 已点名 `login_failures`、`task_runs`、config D0、nullable columns，但尚无逐列 old→new→read/write mapping 与 CHECK/WHERE 改写。

### F-R1-004 · 新合同 backup/rollback 仍未验证

D-006 方向已选，Guardrail 已提出 SPI/metadata/restore contract，但尚无 concrete interface、provider evidence、转换前后 backup 点或 restore verification。

### F-I-005 · checksum/append-only 仍未成为可执行门禁

D-004 已选按模块 v73+ append-only，但 C2 尚未列出 version allocation、72-entry frozen test append 形态、PG type assertion 更新与 conversion descriptor owner。

### F-I-006 · CHECK/index/predicate 列表仍未冻结

至少 jobs state CHECK、recycle partial index、digital entitlement CHECK、login failures comparisons、task-run NULL sentinel 仍需逐项写入设计。

### F-I-010 · wire inventory 仍需独立审与完整覆盖

wire inventory 已落盘，但仍须处理 A-006 指出的漏项（`account_self.go:391`、service credential audit detail、tests、filelibrary/configpkg 例外裁决）并由 independent 复审；当前 formatter 尚未修改。

### F-I-012 · Backup SPI/Service 的 API surface 尚未由用户裁决

用户已选择“统一 Backup SPI/Service”，但尚未决定它是：

- 仅 `internal/backup` 的实现层 service，不进入 `kernel`/模块公共契约；或
- 新增 kernel-level port，供 composition/运维调用。

本选择会影响架构边界、公共 API、审计范围与未来 VP-009/运维交接，不能静默替代。

## 放行

上述 required findings 未闭合前，不得冻结 C2/C3、修改 migration DDL、改公共 formatter、或启动 R2。先向用户询问 F-I-012，再补齐 guardrails，随后调用本地 grok build 做 independent 审计。

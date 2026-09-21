---
id: A-001-r1-self-readiness
doc_type: goal-audit-entry
source: self
auditor: /govern
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · C1/C2/C3 readiness
verdict: conditional
open_required: 4
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-001 · R1 self readiness audit

## 结论

`conditional`。用户已冻结 R1 的关键方案方向，但当前只能证明 inventory v0.1 的初步分组与若干运行时单位；尚不足以冻结逐列合同、实施两方言原地转换或放行 R2。按 P-003/P-005，以下 required findings 保持开放。

## 已核对事实

- VP-013 R1 v1.4 已明确历史形态：SQLite `INTEGER`、PostgreSQL `BIGINT`，单位按列保留秒/毫秒。
- 当前代码至少存在两类运行时单位：`Unix()` 秒（authsession、wallet、dictionary、MFA、notifications、telegram、digitaloffer 等）与 `UnixMilli()` 毫秒（jobs、operationlog、mail outbox/config）。
- 用户已选择 VP-040 新形态：PG `timestamptz(6)`、SQLite fixed-6 UTC RFC3339 `TEXT`、所有绝对时刻列进入分母、两方言各自原地转换、sentinel 0 → NULL。
- inventory v0.1 已落盘，但仍明确标记 collecting。

## Required findings

### F-R1-001 · 逐列 inventory 尚未闭合

- **level**: required
- **影响门禁**: C1/C2、R2
- **主张**: 需要对 48-migration compiled catalog、历史重建/ALTER 路径、当前 runtime repository/reader/writer 做逐表逐列核对。
- **当前证据**: `attachments/r1-time-column-inventory-v0.1.md` 已覆盖主要模块与特殊 sentinel，但还不是机械可核对的完整分母。
- **关闭要求**: 完整清单；每列 old unit、nullable/default/sentinel、reader/writer、目标 codec、migration owner 与证据行号齐全。

### F-R1-002 · 两方言转换 DDL/codec 尚未冻结

- **level**: required
- **影响门禁**: C2/C3、R2
- **主张**: `INTEGER/BIGINT` 秒/毫秒到 `TEXT/timestamptz(6)` 的转换、时区规范化、微秒截断/舍入、索引与排序语义仍未落盘。
- **关闭要求**: 逐列 SQL/Go codec 设计、不可逆点、失败回滚策略、PG `USING` 与 SQLite table-rebuild 形态，并由独立审计核对。

### F-R1-003 · NULL/zero/default 逐列规则尚未闭合

- **level**: required
- **影响门禁**: C2/C3、R2
- **主张**: 用户裁决 sentinel 0 → NULL，但 `locked_until`、`last_login_failure_at`、`mail_config.updated_at DEFAULT 0`、`telegram_config.updated_at DEFAULT 0` 等列的 nullable 化、读写映射与 backfill 尚未逐列冻结。
- **关闭要求**: 每个 nullable/sentinel 列给出 old→new→read/write mapping、约束/default 调整与数据异常处理。

### F-R1-004 · 备份/恢复与失败回滚验证尚未闭合

- **level**: required
- **影响门禁**: C3、R2/R3
- **主张**: VP-013 `pg_dump/pg_restore` 合同与 SQLite snapshot 已有历史证据，但 VP-040 物理形状转换后的验证脚本、转换失败恢复、版本/客户端兼容边界尚未形成本 R1 可执行方案。
- **关闭要求**: 明确复用/扩展 VP-013 证据、转换前后备份点、恢复到新库/副本的校验面与失败回滚界限；不得把旧 catalog dump 证据直接当作新合同已验证。

## 放行结论

在上述 required findings 合法闭合前，不得进入 R2 的 schema/codec 实施，不得将 R1 标记 completed，不得关闭 GOAL-002 或 Root R1 检查点。下一步是补齐 C2/C3 方案证据，然后调用本地 grok build 做 independent 审计。

## 编排响应（2026-09-20）

`F-R1-001` → **fixed**：inventory v0.2 已补齐 90 个 live 绝对时刻列、1 个历史 retired 来源、单位/nullable/sentinel/运行时证据与未决数据抽样清单；见 `attachments/r1-time-column-inventory-v0.2.md` 与 E-003。原 verdict 与原 finding 原文保留。

当前仍开放 required：`F-R1-002`、`F-R1-003`、`F-R1-004`（3 条）。

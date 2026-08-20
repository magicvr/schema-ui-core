---
id: A-005
doc: audit-entry
goal: GOAL-005-r4-repository-surface
source: self
scope: 响应 independent A-004（F-001/F-002 关闭）
verdict: pass
status: recorded
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# A-005 · 响应 independent A-004（instr 改写 + I-002 闭合）

## 范围与区间

- auditor: 本会话编排器（/govern，self 响应）
- type: `response`
- 被响应: `A-004-independent-r4-execution-closeout`（grok-4.6 · reasoning high · `conditional`；F-001 required/high + F-002 required/med）

## 关闭证据表

| Finding | 状态 | 证据路径 |
|---------|------|----------|
| A-004 F-001（9 处运行时 `instr(...)` 在 PG 上失败，required/high） | **fixed** | 7 个文件 9 处全部改写为 `LOWER(col) LIKE '%' || CAST(? AS TEXT) || '%'`（含 wallet/datadictionary/authsession users·notifications·roles/scheduledtasks/operationlog）；`grep instr(` 0 残留；commit `b76d794`；sqlite 全量回归 + live PG 复跑 0 FAIL（`TestFullCatalogPostgresBootstrapIntegration` 含可移植检索断言） |
| A-004 F-002（I-002 到期 required，登记/闭合） | **fixed** | GOAL-005 `00-meta` + `01-decision`：I-002 → **verified**，逐项落盘：INSERT OR IGNORE → `ON CONFLICT DO NOTHING`；`LIKE` → `LOWER(col) LIKE LOWER(?)`（Q 检索）/ `LOWER(col) LIKE '%'||CAST(? AS TEXT)||'%'`（contains 检索）；`COLLATE NOCASE` → `LOWER(col)`；旧 sqlite `instr` → 归入 contains 形态；无 `RETURNING` 需求（插入取 id 用带 PK 的 upsert/回车，见 R3） |

## 仍开放项

- 无 open required（本目标）。R5（升级策略 I-001 / 备份合同 I-004）在新子目标承接。

## 与既有意见的异同

- A-004 原 `conditional` 保留；A-005 fixed 关闭其全部 findings。对 A-003「0 open required」的反对已更正——instr 债补齐。

## 结论与下一步

R4 全部 required 闭合、I-002 verified、sqlite 全量回归 0 FAIL、live PG（bootstrap + 完整启动 + 可移植检索）全绿。**编排器判定：GOAL-005 具备关门条件** → status `done`、progress 6/6；Root R4 ✅（Root 3/5 → 4/5）；创建 **R5 子目标**（双路径验收：SQLite→PG 升级策略 I-001 + PG 备份合同 I-004 + 全链验收）。

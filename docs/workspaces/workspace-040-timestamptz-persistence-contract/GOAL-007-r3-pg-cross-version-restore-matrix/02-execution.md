---
id: GOAL-007-r3-pg-cross-version-restore-matrix
doc: execution
status: done
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-21
updated: 2026-09-21
version: 0.3.0
---

# 执行记录 · GOAL-007-r3-pg-cross-version-restore-matrix

## 执行索引

| E-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| E-001 | 2026-09-21 | R3-C 立项（用户预确认 slug；范围取自 `D-018` §2 第 4 项） | recorded | `02-execution/E-001-goal-created.md` |
| E-002 | 2026-09-21 | 检查点 A：工具兼容行为前置实测 + 矩阵定义与判定口径冻结（`D-001`，关闭 `I-041-009`） | recorded | 见下方事实边界与 `01-decision/D-001` |
| E-003 | 2026-09-21 | 检查点 B：真实迁移链的 9+54 格矩阵实测、形状校验、规则更正（`D-002`）与 `I-041-004` 收口 | recorded | `02-execution/E-003-r3c-matrix-measured.md` |
| E-004 | 2026-09-21 | 检查点 C：响应 `A-002` 三条 recommended（`D-003`）+ 台账同步 + 修正后复跑 + 静默关门 | recorded | `03-audit/A-003-response-to-a002-and-checkpoint-c-closure.md` |

## 事实边界

> 本目标承接 **R3-C**：PG 15/16/17 跨版本 `pg_dump`/`pg_restore` 组合矩阵（逐组合 supported/unsupported 落盘）+ 判据 4「升级后恢复」有界核对，并据此关闭或 residual 化 `I-041-004`。**已关门**（2026-09-21，`done · 3/3`）。
>
> **检查点 A（2026-09-21）**：实测版本（15.19 / 16.15 / 17.11；常驻 15.4）；两条工具规则与两个失败模式；判定口径与组合边界先于矩阵本体冻结（`D-001`）；`I-041-009` → verified。证据 `attachments/r3c-pg-tool-compatibility-probe-v0.1.md`。
>
> **检查点 B（2026-09-21）**：驱动 `internal/backup/pg_cross_version_matrix_test.go`（`VP040_PG_MATRIX=1` 门控，默认 SKIP）在真实 VP-040 迁移链上测出 dump **9 格**（6 supported）与 restore **54 格**（18 supported 且形状校验全通过 / 24 toolgate / 12 serverguc / 0 unexpected）；附带证明迁移链在 **PG 16.15 与 17.11** 上至 v87 可完整应用。`D-001` §3 的外推规则经反例更正为 `client ≥ server`、`client ≥ dumper_client`、`client == 17 ⇒ target == 17`（`D-002`）。`I-041-004` → verified。
>
> **检查点 C（2026-09-21）**：`A-001` self → `A-002` independent `conditional`/**开放 required = 0**（独立复跑 `PASS 104.81s` 与附件一致；25 条分类 oracle + 形状反例）→ 4 条 recommended 全部 fixed（`D-003`：sentinel 探针字面化、44 张分母表存在性检查、`matrixClassify` 收窄 + 15 例常驻 oracle、台账对齐）→ 修正后复跑一致（`PASS 105.07s`）→ `A-003` 关门记录。
>
> **移交 R3-D**：Root 六条判据的退出矩阵（判据 1–5 证据分别来自 `GOAL-002` / `GOAL-003`–`005` / `GOAL-006` / 本目标）、self + independent 关门审计、**用户确认关门**（判据 6）；产品级「支持哪些组合」的书面承诺按 `A-002` 建议留作 R3-D 的 P-004。Root 仍 `active · 2/3`。

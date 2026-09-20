---
id: GOAL-007-r3-pg-cross-version-restore-matrix
doc: execution
status: active
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# 执行记录 · GOAL-007-r3-pg-cross-version-restore-matrix

## 执行索引

| E-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| E-001 | 2026-09-21 | R3-C 立项（用户预确认 slug；范围取自 `D-018` §2 第 4 项） | recorded | `02-execution/E-001-goal-created.md` |

## 事实边界

> 本目标承接 **R3-C**：PG 15/16/17 跨版本 `pg_dump`/`pg_restore` 组合矩阵（逐组合 supported/unsupported 落盘）+ 判据 4「升级后恢复」有界核对，并据此关闭或 residual 化 `I-041-004`。**尚未实施**：组合定义、矩阵实测、升级后恢复核对、审计均为待办；`progress: 0/3`。R3-A/B 已在 `GOAL-006`（`done · 3/3`）完成并关门，本目标不得重开其范围。

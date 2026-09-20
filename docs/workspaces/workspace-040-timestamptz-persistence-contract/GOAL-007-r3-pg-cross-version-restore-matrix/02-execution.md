---
id: GOAL-007-r3-pg-cross-version-restore-matrix
doc: execution
status: active
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-21
updated: 2026-09-21
version: 0.2.0
---

# 执行记录 · GOAL-007-r3-pg-cross-version-restore-matrix

## 执行索引

| E-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| E-001 | 2026-09-21 | R3-C 立项（用户预确认 slug；范围取自 `D-018` §2 第 4 项） | recorded | `02-execution/E-001-goal-created.md` |
| E-002 | 2026-09-21 | 检查点 A：工具兼容行为前置实测 + 矩阵定义与判定口径冻结（`D-001`，关闭 `I-041-009`） | recorded | 见下方事实边界与 `01-decision/D-001` |

## 事实边界

> 本目标承接 **R3-C**：PG 15/16/17 跨版本 `pg_dump`/`pg_restore` 组合矩阵（逐组合 supported/unsupported 落盘）+ 判据 4「升级后恢复」有界核对，并据此关闭或 residual 化 `I-041-004`。
>
> **检查点 A 已完成（2026-09-21）**：
> - 实测环境版本（非沿用文档）：容器 server/client **15.19 / 16.15 / 17.11**；常驻 server **15.4**；归档格式版本 **1.14 / 1.15 / 1.16**。
> - 探测得到两条工具规则：`pg_dump` 要求 client major ≥ server major；`pg_restore` 要求 dumper major ≤ client major ≤ 目标 server major，并识别出两个独立失败模式（归档格式门、server 端 `transaction_timeout` GUC 门）。
> - 判定口径与组合边界已先于矩阵本体冻结（`D-001` §2–§3）；证据 `attachments/r3c-pg-tool-compatibility-probe-v0.1.md`。
> - `I-041-009`（required）→ **verified**。
>
> **尚未实施（检查点 B/C）**：真实 VP-040 迁移链在 16/17 上的应用、真实 schema 的 9 格 dump 与 27 格 restore 逐格记录、升级后恢复的形状校验（§4 四项）、`I-041-004` 收口、self + independent 审计。
>
> `progress: 1/3`。R3-A/B 已在 `GOAL-006`（`done · 3/3`）完成并关门，本目标不得重开其范围。

---
id: GOAL-008-r3-exit-matrix-and-root-closeout
doc: execution
status: done
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-21
updated: 2026-09-21
version: 0.2.0
---

# 执行记录 · GOAL-008-r3-exit-matrix-and-root-closeout

## 执行索引

| E-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| E-001 | 2026-09-21 | R3-D 立项（用户预确认 slug；范围取自 `D-018` §2 第 5 项） | recorded | `02-execution/E-001-goal-created.md` |
| E-002 | 2026-09-21 | 检查点 A：退出判据证据矩阵 + 判据 5 反向核验 + 真实路径实测 | recorded | `02-execution/E-002-r3d-exit-matrix-and-criterion5-sweep.md` |
| E-003 | 2026-09-21 | 用户报告 `dev.cmd start` 失败：复现、根因、修复、端到端与启动器实测 | recorded | `02-execution/E-003-dev-startup-artifact-dir-defect.md` |
| E-004 | 2026-09-21 | 检查点 C：用户确认关门（判据 6）→ Root `done · 3/3`；dev 环境改指专用库 | recorded | `01-decision/D-001-root-closeout-user-confirmation.md`；Root `D-020` |
| E-005 | 2026-09-21 | dev 专用库落地与实测（`schema_ui_dev` 迁移至 head 87；`dev.cmd start/stop` 实测通过） | recorded | `02-execution/E-005-dev-dedicated-db-and-root-closeout.md` |

## 事实边界

> 本目标承接 **R3-D**：Root 六条成功标准的逐条证据矩阵 + 判据 5 反向核验 + 残留/例外清账 + self/independent 关门审计 + **用户确认关门**。**已关门**（2026-09-21，`done · 3/3`）。
>
> **检查点 A（2026-09-21）**：退出矩阵落盘 `attachments/r3d-root-exit-criteria-matrix-v0.1.md`；判据 1–5 **满足**、判据 6 **部分满足**（待用户确认）；判据 5 由编排器亲自扫描核验；判据 3/4 的真实 PG/SQLite 路径实测非 skip。
>
> **检查点 B（2026-09-21）**：`A-001` self → `A-002` independent 关门审计（`conditional`/**开放 required = 0**，自行复跑四测试与判据 5 扫描、逐个核对跨目标关门意见、明确同意交给用户确认关门且拒绝自行标 `done`）→ `A-003` 响应（6 条 recommended 全闭合或登记）。随后用户报告 `dev.cmd start` 失败：`A-004` 记录新 **required `F-I-101` → `fixed`**（`recoveryArtifactsDir` 相对路径 → `docker run -v` exit 125；全分支绝对化 + 边界 fail-closed + 4 条回归 + 一次性库与**启动器本体**实测）→ `A-005` independent 定向复审 **`pass`**/开放 required = 0（独立复现旧失败与新修复、扫过 13 处同类隐患）→ `A-006` 响应（`F-I-102`/`F-I-103` fixed）。
>
> **检查点 C（2026-09-21）**：用户书面确认关闭 Root（`D-001`；Root `D-020`）→ 本目标 `done · 3/3`、Root `done · 3/3`、六条成功标准勾选、投影同步；dev 环境按用户裁决改指专用库 `schema_ui_dev`。
>
> **边界**：VP-040 的波次关闭 / Vision Review 属决策层（`/vision`）动作，**不在**本次关门范围。

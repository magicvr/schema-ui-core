---
id: GOAL-007-r3-pg-cross-version-restore-matrix
title: R3 · PostgreSQL 15/16/17 pg_restore 跨版本矩阵与升级后恢复有界核对（R3-C）
status: active
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-21
updated: 2026-09-21
version: 0.2.0
progress: 1/3
plan_refs:
  - VP-040-timestamptz-persistence-contract
primary_plan: VP-040-timestamptz-persistence-contract
serves_summary: 承接 Root R3 的 R3-C：在固定版本的 PostgreSQL 15/16/17 临时容器上，用 PostgreSQL 自带 client 工具对 VP-040 转换后的库做 pg_dump/pg_restore 组合矩阵，逐组合记录 supported / unsupported，并完成判据 4「升级后恢复」的有界核对，据此关闭或书面 residual 化 I-041-004。
---

# GOAL-007 · PostgreSQL 15/16/17 pg_restore 跨版本矩阵与升级后恢复有界核对

## 概述

承接 `GOAL-001-timestamptz-persistence-contract` 的 **R3 阶段 R3-C 检查点**（Root `D-018` §2 第 4 项）。R3-A/B（`GOAL-006`）已 `done · 3/3`；本目标只做**跨版本备份/恢复矩阵与升级后恢复核对**，不改 Store/迁移/描述符/wire。

**子目标编号与 slug 经用户 2026-09-21 预确认**（AGENTS §11；Root `D-019` §2）：`GOAL-007-r3-pg-cross-version-restore-matrix`。

## 范围

| # | 交付 | 依据 |
|--:|------|------|
| 1 | **组合定义**：固定版本 PostgreSQL 容器（15/16/17）与宿主/容器 client 工具的 server×client 组合清单，含每组合的驱动方式与可复现命令 | `D-018` §2 第 4 项 |
| 2 | **逐组合实测记录**：每个组合的 `pg_dump` / `pg_restore`（或等价路径）结果，**supported / unsupported 均逐条落盘**（命令、版本串、退出码、关键输出） | `D-018` §2 第 4 项；`I-041-004` |
| 3 | **升级后恢复有界核对**：在旧版本上生成产物、在新版本上恢复并校验（VP-040 转换后的 canonical 形状与 C3 校验路径），记录有界范围 | `D-018` §2 第 4 项；判据 4 |
| 4 | `I-041-004` 关闭或用户书面 residual 化 | `D-018` §5 |
| 5 | self + grok independent 审计落盘、required 合法闭合 → 静默关门 | `D-018` §4 R3-C |

## 非目标（本子目标**不**做）

- **不**重开 R1/R2 的冻结决策、canonical SQL/checksum、已落码 descriptor（v1–v87 不可变）。
- **不**改 Store 物理类型、codec 语义、C3 Port 形状与三类产物区分（R2 已冻结）。
- **不**改公共 wire 输出/输入（R3-A/B 已关门）。
- **不**引入 ORM、第三数据库、Redis/MQ/多实例/A3。
- **不**把 Docker 临时容器的跨版本结果当作**生产就绪**证据（`D-017` §3 约束②：其定位是 CI/reproducibility）。
- **不**做 R3-D（退出矩阵/关门）——待本目标完成后再立项。
- **不**新增备份调度/权限/远端存储/保留策略/UI。

## 红线

- 破坏性 migration / 恢复动作只可作用于**一次性或专用测试** database，且必须可丢弃（`D-017` 用户裁决约束）。
- 宿主既有实例的方案与数据不得被本目标修改；容器与宿主库隔离。
- 不以「矩阵跑通」代替独立审计结论；required 未合法闭合不得关门。
- 逐组合结果必须**逐条**记录，禁止用「全部支持」概括未实测组合。

## 成功标准（本子目标检查点，用于 `progress` 派生）

| 检查点 | 判据 | 状态 |
|--------|------|------|
| **A** | 组合定义与驱动落盘（容器版本、client 版本、驱动命令、可复现入口）；每组合的预期判定口径（supported/unsupported 的定义）明确 | **completed**（`D-001` §1–§3、§5 冻结矩阵两轴 9+27 格、三类判定与驱动机制；前置实测 `attachments/r3c-pg-tool-compatibility-probe-v0.1.md`；`I-041-009` → verified） |
| **B** | 全部组合实测并逐条记录（含 unsupported 的原因与退出码）；升级后恢复有界核对完成；`I-041-004` 关闭或书面 residual | pending |
| **C** | self + grok independent 审计落盘、required 合法闭合 → 静默关门 | pending |

`progress: 1/3` 由 A～C 等权派生；**不**放行阶段、**不**关闭 finding、**不**推导 `done`。

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 | 影响门禁 | 状态 | 证据 |
|----|------|----------|----------|------|------|
| I-041-004 | non-blocking（继承，`D-018` §5） | PG 15/16/17 跨版本 `pg_restore` 兼容矩阵 | B | **open**（本目标 B 关闭或 residual） | `D-018` §5；`D-017` §3 约束② |
| I-041-009 | required（本目标新增） | 「supported / unsupported」的判定口径与组合边界（是否需用户裁决某些组合不作为支持承诺） | A | **verified（2026-09-21）** | `01-decision/D-001-r3c-matrix-definition-and-criterion.md` §2–§3；前置实测附件 |

## 已知环境事实（立项时）

- 常驻 PostgreSQL **15.4**（`apps/api/configs/.env` 的 `PG_TEST_*`；破坏性动作只可作用于一次性/专用测试 database）。
- 宿主**无** `psql`/`pg_dump`/`pg_restore` 二进制；`docker` 可用；已实测 `postgres:15-alpine` 内 `pg_dump` 为 **15.19**。
- 本机已有镜像：`postgres:15-alpine`、`postgres:16`、`postgres:17-alpine`（组合定义前需复核，缺失版本按需拉取）。
- 上述为立项时快照；组合定义（检查点 A）须以实测版本串为准，不得沿用本表。

## 父目标

- `GOAL-001-timestamptz-persistence-contract`（R3 检查点）

## 关门条件

**A～C 全部完成**、`03-audit` 的 self 与 independent 意见落盘、required 合法闭合后才可静默关门；关门**不**等于 Root 关门（Root 关门仍须 R3-D 与用户确认）。

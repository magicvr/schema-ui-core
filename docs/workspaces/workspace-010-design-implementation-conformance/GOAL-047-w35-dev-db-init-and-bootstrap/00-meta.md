---
id: GOAL-047-w35-dev-db-init-and-bootstrap
title: W35 · dev/test 数据库初始化快捷方式与自举修复
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-09-21
updated: 2026-09-21
version: 0.2.0
progress: 2/3
plan_refs:
  - VP-010-design-implementation-conformance
primary_plan: VP-010-design-implementation-conformance
vision_ref: schema-ui-core-admin-foundation@0.4.0
serves_summary: 用户报告实例重置后 .\dev.cmd start 启动失败（PG 目标库不存在）后立项：为 dev/test 提供明确的一次性初始化快捷方式，修复既有 e2e-pgset 在 DB_NAME 不存在时无法自举的缺陷，补齐 README/QUICKSTART 初始化章节，并让启动错误对操作者可操作。
---

# GOAL-047 · W35 · dev/test 数据库初始化快捷方式与自举修复

## 概述

用户在 `[workspace-040]` Root 关门后重置了常驻 PostgreSQL 实例并报告（2026-09-21）：

> 我重置了一下贡献pg库，现在我们是否需要做一下数据库初始化处理才能让程序自动启动，目前直接执行 `.\dev.cmd start` 是报错的。同时这是否指向一个我们可能需要处理待办：我们需要给出明确的初始化数据库的快捷方式。如果是的话，我们需要把它登记到总路线图当中。

诊断（编排器原样复现，2026-09-21）：

1. **PG 路径要求目标库已存在**：`openStore` 只对配置库 `Ping` 后应用 catalog，**从不** `CREATE DATABASE`（SQLite 无需初始化）。实例重置后：
   `LIFECYCLE_START_FAILED [core.auth-session]: open store: postgres ping: FATAL: database "schema_ui_dev" does not exist (SQLSTATE 3D000)`
2. **既有 helper 无法自举**：仓内已有 `apps/api/cmd/e2e-pgset create|drop|verify|list`，但（a）只面向 `schema_ui_e2e_*`；（b）`maintenanceDSN()` 连接的是 `DB_NAME` —— 当 `DB_NAME` 本身不存在时，该工具自身同样 3D000，**无法创建那个库**（须以 `DB_NAME=postgres` 覆盖才能工作）。
3. **无文档、无快捷方式**：`dev.cmd` 无 init 子命令；README/QUICKSTART 无 PG 初始化步骤；启动错误不提示如何初始化。

本目标立项依据：用户 2026-09-21 P-004 裁决 —— 结构上**在 workspace-010 新加一个子目标**（不新开 VP）；范围选**完整**（初始化快捷方式 + 自举修复 + 文档 + 启动错误可操作化）。编号与 slug 经用户确认：`GOAL-047-w35-dev-db-init-and-bootstrap`。

登记出处：`docs/vision/roadmap.md` §未决项统一登记 §一「待立项 · dev/test PG 数据库初始化快捷方式」（commit `934c60f5`）。

## 范围

| # | 交付 | 依据 |
|--:|------|------|
| 1 | **修复 `e2e-pgset` 自举**：维护连接不再依赖「`DB_NAME` 已存在」；`DB_NAME` 不存在时回退到一定存在的维护库（`postgres` → `template1`），使工具能在空实例上创建第一个库 | 复现证据（§概述 2） |
| 2 | **一次性初始化快捷方式**：一条命令建齐 **dev 与 test** 库（`dev.cmd init-db`，或等价 `cmd/dbsetup`；同时接受显式库名/附加库名） | 用户范围裁决「完整」 |
| 3 | **启动错误可操作化**：PG 目标库不存在（SQLSTATE 3D000）时，错误信息给出下一步（指向 #2 的命令），而不是裸 `FATAL` | 用户范围裁决「完整」 |
| 4 | **文档**：README / QUICKSTART 增「初始化数据库」章节（何时需要、命令、多方言差异、失败排查） | 用户范围裁决「完整」 |
| 5 | **回归与实测**：删除库 → 跑 #2 → `dev.cmd start` 全绿的可复跑证据；#1 的单元/集成测试 | P-002 |
| 6 | self + grok independent 审计落盘、required 合法闭合 → 静默关门 | 项目惯例 |

## 非目标（本目标**不**做）

- **不**改变「应用不自动建库」的既有语义（不把 `CREATE DATABASE` 塞进启动路径）——自动建库会掩盖拼错的库名/连错实例，且与 `W21` 启动身份判定取向冲突。
- **不**改迁移链、canonical SQL/checksum、codec、wire 合同或备份 Port（与 VP-040 已关门范围无关）。
- **不**引入 docker/容器作为 PG 的必需依赖（初始化只依赖数据库连接与 CREATEDB 权限）。
- **不**做 GUI/运维面板；不做权限/保留策略/远端存储。
- **不**改 `.env` 的语义或键名；不把秘密写入仓库。
- **不**新开 VP（用户已裁决：本波次挂 workspace-010）。

## 红线

- 初始化命令必须**幂等**：库已存在时给出明确结果而非报错（`create` 语义），且**不得**删除或覆盖既有库（`drop` 必须是显式独立动作）。
- 维护连接回退**只**用于「连接维护库」这一步，不得让业务连接悄悄连到别的库。
- 不得把 `.env` 中的秘密输出到日志/错误信息/回复。
- 不得以「测试通过」代替独立审计结论；required 未合法闭合不得关门。

## 成功标准（本子目标检查点，用于 `progress` 派生）

| 检查点 | 判据 | 状态 |
|--------|------|------|
| **A** | `e2e-pgset` 自举修复（`DB_NAME` 不存在时可创建第一个库）+ 初始化快捷方式落码（dev 与 test 库一次建齐、幂等）+ 启动错误给出可操作提示；均有可执行测试 | **completed**（`D-001`；`internal/pgsetup` 测试；`TestMatrixClassify`/startup hint；`dev.cmd init-db` reset 实测） |
| **B** | 文档章节落盘（README/QUICKSTART）；**实测**：删除两个库 → 快捷方式 → `dev.cmd start` 全绿（API `/readyz 200` + Web 200）；全仓回归绿 | **completed**（`E-002`；reset → init-db → dev.cmd start/stop；PG test path；Go 64/64） |
| **C** | self + grok independent 审计落盘、required 合法闭合 → 静默关门 | pending |

`progress: 2/3` 由 A～C 等权派生；**不**放行阶段、**不**关闭 finding、**不**推导 `done`。

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 | 影响门禁 | 状态 | 证据 |
|----|------|----------|----------|------|------|
| I-047-001 | non-blocking | 初始化快捷方式的 CLI 形态细节（子命令名、是否支持附加库名/多方言参数） | A | **verified**（`D-001`）：`dev.cmd init-db` → `cmd/dbsetup`，默认 dev+test，支持 `--dev-only`/`--test-only`/附加安全库名 | 用户 2026-09-21 裁决；`01-decision/D-001`；`E-002` |
| I-047-002 | non-blocking | 是否让 e2e 专用库也纳入同一快捷方式（当前 `e2e-pgset` 负责 `schema_ui_e2e_*`） | B | **verified（本目标边界）**：不合并；e2e 生命周期继续由 e2e 自管；`e2e-pgset` 复用自举层但不纳入默认 dev/test 集 | `cmd/e2e-pgset`；`apps/web/README.md` §e2e-postgres；`D-001` |

## 父目标

- `GOAL-001-design-implementation-conformance`（workspace-010 Root）

## 关门条件

**A～C 全部完成**、`03-audit` 的 self 与 independent 意见落盘、required 合法闭合后才可静默关门。

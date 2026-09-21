---
id: D-020-root-closeout-user-confirmation
doc: decision-entry
status: accepted
parent: null
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# D-020 · Root 关门用户确认（判据 6）与 dev 环境指向裁决

## 决定的来源

- 用户 **2026-09-21** 在 Root 关门确认包（`GOAL-008/03-audit/A-006` §4）中书面裁决：**确认关闭 Root**。
- 前置事实：退出判据矩阵 `GOAL-008/attachments/r3d-root-exit-criteria-matrix-v0.1.md`（判据 1–5 满足、判据 6 待确认）；self 审计 `GOAL-008/A-001`；independent 关门审计 `A-002`（`conditional`/开放 required = 0）与定向复审 `A-005`（**`pass`**/开放 required = 0，独立复现 `F-I-101` 的旧失败与新修复）。
- 用户报告的 `.\dev.cmd start` 启动失败已按用户要求**先修复**（`F-I-101`，required → `fixed`），并经独立复现与启动器本体实测（API `/readyz 200` + Web `200`）。

## 1. 关门裁决（判据 6）

**用户选择：确认关闭 Root。** 据此：

- `GOAL-008-r3-exit-matrix-and-root-closeout` → `done · 3/3`；
- `GOAL-001-timestamptz-persistence-contract`（Root）→ **`done · 3/3`**，六条成功标准勾选；
- `goal-tree.md` 树/表/叙述、Root `00-meta.md`、`docs/vision/roadmap.md` 投影与 `workspace.md` 纲领段同步。

**关门依据**（不放行任何未闭合项）：

| 项 | 状态 |
|----|------|
| 判据 1–5 | 满足（逐条可核对证据，见矩阵；含真实 PG/SQLite 恢复与升级路径，本轮实测非 skip） |
| 判据 6 | **用户 2026-09-21 书面确认**（本条） |
| 跨目标开放 required | **0**（GOAL-002～008 逐个核对） |
| 未闭合 recommended | `GOAL-008/A-005` 的 `F-I-102`/`F-I-103` 已 `fixed`（`A-006`）；`F-I-006`（vision/workspace 投影）随本条同步 |
| residual | R1 `D-021` 的 `F-I-005` 已按 `fixed` 闭合（`GOAL-002/A-048`）；无未闭合 residual |

## 2. dev 环境指向裁决（用户 2026-09-21）

**用户选择：改为专用 dev 库。** 处置：

- `apps/api/configs/.env`（gitignored）的 `DB_NAME` 由 `postgres` 改为 **`schema_ui_dev`**（专用 dev 库），使 `dev.cmd start` 不再把完整迁移链跑到共享常驻实例的**维护库** `postgres` 上；
- 在常驻实例上创建 `schema_ui_dev` 并按需迁移（`dev.cmd start` 首次启动即完成迁移）；
- `PG_TEST_*` 与其它键（`ADMIN_PASSWD`、mail/telegram 键）**不变**——测试路径仍按原约定；
- 该文件**不入库**，变更记录在本条与本目标执行台账。

## 3. 本次关门**不**包含的动作（边界）

- **VP-040 的波次关闭 / Vision Review 不在此列**：VP 状态与 Vision Review 属**决策层**（`/vision`）动作。本次只把 `docs/vision/roadmap.md` 的**投影**更新为「工作区 Root 已关门」的事实陈述，**不**把 VP-040 的 `status` 改为 `closed`、**不**新建 VRev。用户如要关闭该 VP，走 `/vision`。
- **不**重开任何冻结决策、canonical SQL/checksum（v1–v87 不可变）、codec、C3 Port、wire 合同或跨版本矩阵结论。
- **不**把容器矩阵结果改写成生产就绪证据（`D-017` §3 约束②）。

## 4. 未选方案

- **用户未选「暂不关门」**：`dev.cmd start` 缺陷已修复并独立验证，未留下阻断项。
- **保持 `.env` 指向 `postgres`**：用户未选；继续指向共享维护库会让每次 dev 启动在该库上跑迁移，风险不可控。
- **由编排器自行把 Root 标 `done` 而不问用户**：违反判据 6 与 P-004；编排器在 `A-002`/`A-005`/`A-006` 中均明确拒绝自行标 `done`，等待本条确认。

## 5. 影响与边界

- 本决策**只**关闭工作区 Root 与 R3-D；工作区 `goal-tree` 显示全部 8 个目标 `done`。
- 后续若要新增工作（新波次/新目标），按 `docs/architecture/workspace-protocol.md` 与 `/vision` 的判定树处理，不在本 Root 内续做。

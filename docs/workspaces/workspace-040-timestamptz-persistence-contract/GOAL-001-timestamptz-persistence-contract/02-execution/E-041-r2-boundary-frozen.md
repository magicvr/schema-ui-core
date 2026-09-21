---
id: E-041-r2-boundary-frozen
doc_type: goal-execution-entry
status: recorded
date: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-041 · R2 边界与信息门禁冻结

## 事实

1. **用户 2026-09-20 裁决**：确认关闭 `GOAL-002`，**并进入 R2**；四条 recommended（F-I-008/009/025/028）随 R2 处理。
2. **R2 边界落盘**：Root 新增 **`D-016-r2-boundary.md`**（`status: accepted`），按 `docs/architecture/workspace-protocol.md`「新建阶段子目标前须先在 Root 冻结边界与信息门禁」的要求，冻结：
   - **目标**：把 R1 冻结的合同落成可执行代码——15 个 conversion descriptor、共享时间 codec、各模块仓储读写改造——并**首次记录真实 `MigrationChecksum`**。
   - **范围（10 项）**：SQLite 15 descriptor（F-5）、PG 显式 DDL（毫秒族整数拆分式）、canonical SQL + 真实 checksum 记录、共享 codec（`apps/api/internal/temporal`）、仓储读写与谓词改造（无过渡期）、测试改写与**金额列断言拆分**、边界测试重定向、`rebuildOperationLog` fail-closed 断言、Backup Port 类型表面与 provider、未发布 baseline 的发布前调整。
   - **非目标**：不重开 Root `D-004`；不改 v1–v72；不引 ORM/第三库/Redis/MQ；不泄漏驱动类型；**不做 VP-020 回归矩阵**（R3）；不做备份调度/鉴权/远端存储/UI；不把 residual 或临时容器验证当已验证事实。
   - **完成判据 M1～M4**：M1 codec + 单测；M2 15 descriptor 落码 + **真实 checksum 已记录** + v1–v72 不变；M3 仓储与谓词改造 + **双方言**回归 + 金额列保持 `bigint`；M4 `D-021` residual 三项完成并经 independent 复审 → residual 关闭 + R2 关门审计。
   - **信息门禁**：新增 Root 级编号段 **`I-041-NNN`**——`I-041-001`（Go codec 公共 API 形态，M1 前）、`I-041-002`（公共 wire formatter 改造落点与是否属 R2，M3 前）、`I-041-003`（PG 可执行验证环境：本机无常驻 PG，`D-021` 只授权临时容器，M3 前）均 **`open` required**；`I-041-004`（PG 跨版本兼容矩阵）为 `non-blocking` deferred 至 R3。
   - **渐进子目标策略**：不预创建全部子目标；首个候选为 **`GOAL-003-r2-codec-and-descriptor-m1-m2`**（M1/M2），**slug 须经用户确认后落盘**（AGENTS §11 硬约束，禁止静默默认 slug）。
3. **Root 状态同步**：`00-meta.md` 路线图 R2 由 `pending` → **`active`** 并补入 R2 边界摘要；`01-decision.md` 决策索引补 D-016；`goal-tree.md` 纲领路线图同步。
4. **本轮未创建任何 R2 子目标**——按上述策略，待用户确认 slug 后立项。

## 证据

- Root 决策：`GOAL-001/01-decision/D-016-r2-boundary.md`；索引：`GOAL-001/01-decision.md`。
- Root 路线图：`GOAL-001/00-meta.md`；工作区树：`goal-tree.md`。
- 既有约束对位：Root `D-004`（模块归属）、`D-014`（未发布 baseline）；child `D-017`（checksum 约定）、`D-018`（行拷贝/无过渡期）、`D-019`（F-5 + 两次重建 + `rebuildOperationLog` 断言 + PG 显式 DDL）、`D-020`（测试载体与重定向）、`D-021`（residual 范围与复审触发）。
- 本轮 `apps/**` **未修改**。

## 状态评估

- **Root `progress: 1/3`**（R1 completed；R2 `active` 但**未完成**；R3 pending）。
- **R2 未放行生产 schema 变更**——每次迁移仍须遵守 `D-017`–`D-021` 与 Root `D-004` 的不变量；`I-041-001`～`003` 在对应门禁前关闭。
- **下一步**：就首个 R2 子目标的 slug 与边界向用户确认，然后按 `D-016` §6 立项（M1/M2）。

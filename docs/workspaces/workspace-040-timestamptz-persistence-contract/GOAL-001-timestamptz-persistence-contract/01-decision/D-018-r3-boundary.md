---
id: D-018-r3-boundary
doc: decision-entry
status: accepted
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-018 · R3 阶段边界与信息门禁

## 决定的来源

- **R2 已关门**（`GOAL-005` `done · 3/3`，独立关门审计 `A-004` `pass` / open required = 0，2026-09-20），Root 纲领路线图 R2 → `completed`，Root `progress: 2/3`。
- 按 `docs/architecture/workspace-protocol.md`，**新建阶段子目标前必须先在 Root 冻结该阶段边界与信息门禁**——本决策即 R3 的边界冻结。
- 既有约束：Root `D-005`（wire 输入兼容）、`D-009`（wire 非 DB 例外范围）、`D-016` §3（R3 非目标：**公共 wire formatter 实施归 R3**）、`D-017` §3 约束②（Docker 版本矩阵属 CI/release reproducibility，不是 `I-041-003` 的关闭条件）、child `D-021`（F-I-005 residual，已在 R2 内按 `fixed` 闭合）。
- 输入清单：`GOAL-002/attachments/r1-public-wire-inventory-v0.1.md`（§C2/R3 closure requirements 1–5）。

## 1. R3 目标

**读写/时区回归、公共 wire 输出与输入合同落码、备份有界核对与跨版本矩阵、退出矩阵与 Root 关门**（Root 纲领路线图 R3 原话），并把 Root 成功标准的**判据 3、4、6** 收口。

## 2. R3 范围（in scope）

| # | 交付 | 载体 |
|--:|------|------|
| 1 | **公共 wire formatter 实施**：一个 shared fixed-6 UTC formatter（`YYYY-MM-DDTHH:MM:SS.ffffffZ`）替换 `rfc3339Milli` 与 inventory 列出的全部 inline 布局；输入 parser 按 `D-005` 兼容 0/3/6/9 位小数与 `+00:00` 等价 offset、**拒绝无时区** | `apps/api/internal/handler/rfc3339.go` 及 inventory 的 Go 面 |
| 2 | **Go 与 Web fixture 同步**（inventory 逐项）+ **单位族矩阵**：秒 / 毫秒 / 可空 / sentinel 转换 各至少一个 endpoint | `internal/handler`、`apps/web` 测试；`cmd/server` |
| 3 | **VP-020 展示/输入时区回归矩阵**：会话/用户时区设置与 UTC 存储的展示 round-trip（`I-040-004`） | Web 组件/单测（载体见 `I-041-007` 裁决） |
| 4 | **备份有界核对 + PG 跨版本矩阵**：C3 harness 在 PG 15/16/17（固定版本临时容器）上的 `pg_dump`/`pg_restore` 组合，**逐组合记录 supported / unsupported**（`I-041-004`），并复核判据 4 的「升级后恢复」路径 | `internal/backup`、证据附件 |
| 5 | **退出矩阵 + 关门**：Root 六条判据的逐条证据矩阵、self + grok independent 关门审计、**用户确认关门**（判据 6） | 本 Root 台账 + 子目标台账 |

## 3. R3 **非目标**（out of scope，不得借 R3 实施）

- **不**重开 R1/R2 的冻结决策、已落码 descriptor、canonical SQL/checksum（v1–v87 不可变）。
- **不**改 Store 的物理类型、codec 语义、C3 的 Port 形状与三类产物区分（R2 已冻结）。
- **不**引入 ORM、第三数据库、Redis/MQ/多实例/A3（Root 红线）。
- **不**把 `pgtype`/`pgx`/SQLite 驱动类型泄漏进 handler/模块公共契约。
- **不**把 Docker 临时容器的跨版本结果当作生产就绪证据（`D-017` §3 约束②）。
- **不**扩展 VP-020 的展示/输入能力本身（R3 只做**回归矩阵**，不新增时区功能）。
- **不**把「非 DB」文本（邮件正文、audit `detail` JSON、任意 payload 字段）纳入固定 6 位输出（`D-009`）。

## 4. R3 完成判据（检查点，用于 progress 派生）

| 检查点 | 判据 |
|--------|------|
| **R3-A** | shared fixed-6 formatter 落码并替换 inventory 的 Go 面；输入端兼容矩阵（0/3/6/9 位 + `+00:00` + 拒绝无时区）有可执行测试；Go/Web fixture 同步 |
| **R3-B** | 单位族矩阵（秒/毫秒/可空/sentinel 各至少一个 endpoint）+ VP-020 会话时区展示 round-trip 通过；`I-040-004` 关闭 |
| **R3-C** | PG 15/16/17 跨版本 `pg_restore` 矩阵逐组合记录（supported/unsupported 均落盘）；判据 4 的升级后恢复 bounded cross-check 完成；`I-041-004` 关闭或书面 residual |
| **R3-D** | 退出矩阵（六条判据）落盘；self + independent 关门审计通过、开放 required = 0；**用户确认关门** |

> `progress` 由 R3-A～D 等权派生；**不**放行阶段、**不**关闭 finding、**不**推导 `done`。

## 5. R3 信息需求与门禁（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 状态 |
|----|------|-----------------|----------|--------------|------|
| I-040-004 | required | VP-020 展示/输入与 UTC 存储的回归矩阵如何覆盖会话时区 | R3-B | R3-B 前 | **open**（载体由 `I-041-007` 裁决） |
| I-041-004 | non-blocking（继承） | PG 15/16/17 跨版本 `pg_restore` 兼容矩阵 | R3-C | R3 前 | **open**（R3-C 收集并逐组合记录） |
| I-041-007 | required（本目标新增） | VP-020 矩阵的**载体与验收口径** | R3-B | R3-B 前 | **verified（用户 2026-09-20 P-004 裁决）**：Go 单测锁 wire 形状与 parser 兼容；Web 用**组件/单测**锁会话时区展示 round-trip；**不**引入浏览器 e2e 依赖 |
| I-041-008 | required（本目标新增） | 公共 wire 输出的**破坏性**：现网/前端是否已有依赖 3 位小数的契约或 fixture；若需兼容期，范围与时长 | R3-A | R3-A 前 | **collecting**（`D-005` 输入兼容已冻结；输出由 `D-003` 冻结为固定 6 位） |

- **到期 open required 阻断对应门禁**；本表新增编号段 `I-041-NNN`（Root 级）继续使用。

## 6. 渐进子目标策略

- 按 `D-004`/`D-016` 的做法，R3 **不预创建全部子目标**；先冻结边界，随后**按检查点渐进立项**。
- 首个候选子目标：**`GOAL-006-r3-wire-formatter-and-unit-family-matrix`**（R3-A + R3-B），对应 formatter 实施、fixture 同步与单位族/时区矩阵。
- R3-C（跨版本矩阵 + 备份核对）与 R3-D（退出矩阵 + 关门审计）待 R3-A/B 完成后再立项。
- **子目标命名与 slug 须经用户确认后落盘**（AGENTS §11 硬约束：禁止静默默认 slug）。

## 未选方案

- **把 R3 全部拆成 4 个子目标一次性立项**：违反 P-001「按阶段创建」，且 R3-C 的范围取决于 R3-A/B 的实际形态。未采用。
- **把 wire formatter 留在 R2**：`D-016` §3 已按用户裁决划归 R3（R2 只做 Store/持久化层）；重开该裁决会破坏已关门的 R2 边界。未采用。
- **用 Docker 容器作为「生产就绪」的跨版本证据**：`D-017` §3 约束②明确其只作 CI/release reproducibility。未采用。

## 影响与边界

- 本决策只冻结 **R3 的边界与门禁**；**不**放行任何实现，**不**关闭 `I-040-004`/`I-041-004`。
- Root 关门（判据 6）仍须**用户确认**；本决策**不**改变该要求。

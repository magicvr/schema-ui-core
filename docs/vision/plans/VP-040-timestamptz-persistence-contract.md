---
doc_type: vision-plan
id: VP-040-timestamptz-persistence-contract
title: DB 时间列 timestamptz 持久化合同
status: active
vision_ref: schema-ui-core-admin-foundation@0.4.0
lead_workspace: workspace-040-timestamptz-persistence-contract
created: 2026-09-19
updated: 2026-09-20
version: 0.2.2
parent: null
---

# VP-040 · DB 时间列 timestamptz 持久化合同

## 状态、激活与关门门禁

| 项 | 值 |
|-----|-----|
| status | **`active`**（2026-09-20 · v0.2.2 · 1 个 delivery 区 · `workspace-040-timestamptz-persistence-contract`） |
| 组合位置 | **架构分支 · C1**（`RES-T03-tz` / `RT-T03`）；承接 VP-035 标为「现在修」、因当时红线禁止改 schema 而只登记的时间列合同 |
| 计划阶段 Vision Review | [VRev-101](../reviews/VRev-101-vp039-vp040-planned.md) self `pass`（0 required；本 VP 为同审查的停放意图） |
| 激活门禁 | **已满足（2026-09-20）**：① [VP-039](VP-039-version-maintenance-diagnostics.md) 已 `closed`；② `I-040-001` 已由用户 P-004 选择 R1 方向（PG `timestamptz(6)` + SQLite fixed-6 UTC RFC3339 `TEXT`，逐列证据仍待冻结）；③ 架构类 freshness `6197e802` → `b0a6789b` PASS；④ 激活就绪 self Review [VRev-104](../reviews/VRev-104-vp040-timestamptz-persistence-contract-activation.md) `pass`；⑤ slug/Root 已由 `/govern` scaffold。**本文件仍不是实现完成或关门证据。** |
| 基础设施边界 | 不消耗 A3 多实例、Redis、MQ、搜索引擎 trigger；不引入 ORM / 第三库；不重开 VP-013 方言决策 |
| 与 VP-039 | **正交、串行**。VP-039 波次已于 2026-09-19 `closed`；本 VP 已于 2026-09-20 激活；禁止把 schema 迁移并进维护提示 VP |

## 用户已裁决（2026-09-19 · P-004）

| 项 | 裁决 |
|----|------|
| 结构 | **另立本 VP**，不并入 VP-039，不塞 VP-010 |
| 激活时机 | **VP-039 波次之后**（VP-039 `closed` 或用户书面改序） |
| 方言物理类型 | **未冻结**（`I-040-001`）；SQLite 无真 `timestamptz`，不得把「改一列类型」当成已验证方案 |

## R1 用户裁决（2026-09-20 · P-004）

R1 子目标 `[workspace-040] GOAL-002-r1-contract-and-denominator-freeze` 已承接用户方案选择：PostgreSQL 字面 `timestamptz(6)`；SQLite 固定 6 位 UTC RFC3339 `TEXT`；全部绝对时刻列纳入分母，ID/duration/step/version/计数/金额/flag 排除；SQLite 与 PostgreSQL 各自原地转换；不提供 SQLite→PostgreSQL 产品级搬运器；语义 sentinel `0` 转为 `NULL`；公共 API 时间输出统一为 6 位微秒 RFC3339 UTC `Z`，入站兼容合法 RFC3339 变体后规范化；C3 建立统一 Backup SPI/Service，由 SQLite/PG native provider 提供 metadata、verification 与 restore-to-new-db，不做调度/权限/远程存储/保留策略/UI。上述是 R1 方向裁决，不替代逐列 inventory、转换失败策略与 self/independent 审计。

## 意图

Admin 时区 / 数字 / 货币**展示与输入**已由 VP-020 交付，但持久化层时间列仍普遍为 SQLite 兼容的 `INTEGER` epoch。VP-035 R2/R3 核对：`apps/api` 内 `timestamptz` 命中 0；`RT-T03` 已由本 VP 承接并进入 `active`。这不是符合性漏做（从未写入已交付分母），而是架构分支正在实施的持久化合同。

本 VP 冻结并交付 **Store 时间列合同**：生产权威 PostgreSQL 使用 `timestamptz`（或与之合同平等的物理类型）；SQLite 内嵌默认必须合同平等、不得残缺；双方言走同一不可变迁移台账；handler / 模块公共契约继续只打本模块 Repository，禁止把驱动时间类型泄漏进公共面。

本 VP 是架构合同，不是 Admin 时区 UX 重做，也不是 PITR / 多实例 / ORM。

## 首波范围与边界

| 范围 | 本 VP 首波 | 不在本 VP |
|------|-----------|-----------|
| 合同 | 时间列逻辑语义（UTC 绝对时刻）、PG 物理类型、SQLite 对等物理类型、读写编解码 | 把 SQLite 假装成有原生 `timestamptz`；第三库 |
| 迁移 | 纳入首波分母的既有 `*_at` / 过期类时间列，双方言 checksum 台账可 apply | 一次性改所有 INTEGER（含金额、开关、version、bot_id）；无回滚搬运器承诺须在 R1 显式声明 |
| 代码面 | Store / 模块 Persistence 编解码与回归；公共契约不出现 `pgtype` / driver 时间类型 | 重开 VP-013 端口形状；把业务 handler 改成直接扫 `*sql.Tx` |
| 与 VP-020 | 消费已交付的展示/输入时区合同；核对持久化时刻与展示时区不漂移 | 重做 locale 数字/货币；改用户时区偏好存储格式（除非 R1 证明必须） |
| 基础设施 | 单进程 + 双方言基线 | Redis / MQ / 多实例 / KMS / PITR |

## 与相邻 VP / 路线图的边界

| VP / 方向 | 关系 |
|-----------|------|
| **VP-013** | **消费**双方言端口与全局 checksum 台账，**不重开** A1；不引入 ORM |
| **VP-020** | 展示/输入已交付；本 VP 只补持久化层。漂移问题归本 VP 回归，不重开 VP-020 |
| **VP-035** | 评估来源；本 VP 实施 C1，不重做 as-built 矩阵 |
| **VP-039** | 产品面维护提示；时间列合同不进 039 退出分母 |
| **VP-010** | 不是 as-designed 漏实现；若迁移中发现符合性文档分叉，转 VP-010 文档卫生 |
| **VP-009** | 迁移安全/注入/备份面缺陷归 009 |
| **RT-P05 / PITR** | 仍 gated；本 VP 只要求既有 dump/restore 路径在类型变更后仍能启动鉴权（范围由 R1 冻结） |

## 方向级退出判据

在同时满足下列方向时，本 VP **可以**有界或完整关门（证据必须在工作区目标内）：

1. **合同冻结**：PG 物理类型、SQLite 对等物理类型、UTC 语义、NULL/零值、编解码、禁止泄漏进公共契约——均已书面冻结（`I-040-001`）。
2. **分母迁移**：R1 冻结的时间列双方言均已迁移 + checksum；未纳入分母的 INTEGER（金额、flag、version 等）显式排除。
3. **读写正确性**：写入的绝对时刻与读回一致；VP-020 时区展示不因存储形状改变而漂移；双方言回归（含至少一条 PG 路径）。
4. **备份面有界核对**：在 R1 冻结范围内，既有 SQLite 快照 / PG dump 路径要么可升级后恢复，要么有书面 residual（无产品搬运器须点名）。
5. **范围保持**：未引入 ORM / 第三库 / Redis / MQ / 多实例；未把 Admin 维护提示或业务域混入。
6. **证据与审计**：退出矩阵与必要独立意见已落盘，开放 required = 0，并经用户确认关门。

## 纲领路线图（实现层由 `/govern` 承接；工作区已建立，R1/R2/R3 尚未实施）

```text
R1 合同与分母冻结：方言物理类型、列清单、零值、备份 residual
  → R2 双方言迁移 + Store 编解码
  → R3 读写/时区回归 + 备份有界核对 + 证据与关门
```

## 信息需求（P-005）

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|--------------|------|----------|----------|------------------|------|-------------|-------------|
| I-040-001 | SQLite 用什么物理类型与 PG `timestamptz` 合同平等？（TEXT RFC3339 / INTEGER epoch+约定 / 其他） | required | 阻断 R1 冻结与 R2 迁移；用户已选方向，逐列证据未闭 | R1 | 对照 VP-013 历史合同与用户 D-002；核对 `timestamptz(6)` / fixed-6 RFC3339 TEXT 编解码 | collecting | R1 冻结前复核 | 用户 P-004：PG `timestamptz(6)`；SQLite fixed-6 UTC RFC3339 `TEXT`；sentinel 0 → NULL；待逐列证据 |
| I-040-002 | 哪些列进首波分母？全部绝对时刻列是否全覆盖？ID/duration/step/version/计数/金额/flag 如何排除？ | required | 阻断 R1/R2 | R1 | 全仓扫描 compiled catalog 与 runtime repository，形成逐列 inventory | collecting | R1 冻结前复核 | 用户 P-004：全部绝对时刻列纳入；非时间整数类别排除；待 inventory |
| I-040-003 | 存量库升级策略与备份 residual？无产品搬运器是否再次声明？ | required | 阻断 R1 方案与判据 4 | R1 | 对照 VP-013/016 dump 路径；设计双方言原地转换与失败恢复 | collecting | R1 冻结前复核 | 用户 P-004：SQLite/PG 各自原地转换；不提供 SQLite→PG 产品搬运器；具体策略待 R1 |
| I-040-004 | 与 VP-020 展示合同的回归矩阵（会话时区、UTC 存储） | required | 阻断 R3 | R1 | 复用 VP-020 验收用例，补存储形状变更对照 | open | — | — |
| I-040-005 | 激活前架构类 freshness；且 VP-039 已 `closed` 或用户书面改序 | required | **阻断激活** | 激活前 | `/vision` 核对 VP-039 status + 五域 freshness | verified | — | VRev-104 `pass`：VP-039 `closed`，freshness `6197e802` → `b0a6789b` PASS，slug 已确认并开区 |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-040-timestamptz-persistence-contract | GOAL-001-timestamptz-persistence-contract | delivery | 2026-09-20 | `active` · Root `active · 0/3`；R1/R2/R3 由 `/govern` 按路线图承接 |

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| — | — | — | — | — |

## 规划修订短史

| date | change |
|------|--------|
| 2026-09-19 | 初创 `planned` v0.1.0 · 0 区 · 停放。用户确认：C1 另立本 VP，不并入 VP-039，不塞 VP-010；激活硬门禁 = VP-039 波次之后（或书面改序）。计划阶段 self = [VRev-101](../reviews/VRev-101-vp039-vp040-planned.md)。 |
| 2026-09-20 | 用户指令走流程激活：`I-040-005` verified；激活 self = [VRev-104](../reviews/VRev-104-vp040-timestamptz-persistence-contract-activation.md) `pass`；VP-040 `planned → active` v0.2.0，lead `workspace-040-timestamptz-persistence-contract` 交 `/govern` 开区。随后用户 P-004 裁决 R1 方向：PG `timestamptz(6)` + SQLite fixed-6 UTC RFC3339 `TEXT`；全部绝对时刻列纳入分母；双方言各自原地转换；sentinel 0 → NULL；不提供 SQLite→PG 产品搬运器。VP-040 修订为 v0.2.1，逐列证据由 GOAL-002 承接。 |
| 2026-09-20 | 用户 P-004 追加公共 wire 裁决：所有公共时间输出统一为 6 位微秒 RFC3339 UTC `Z`；需同步 formatter、parser、fixtures 与 VP-020 回归。VP-040 当前版本保持 v0.2.1，D-003/E-004 落在 workspace-040。 |
| 2026-09-20 | 用户 P-004 裁决 C3 建立统一 Backup SPI/Service：SQLite/PG native provider + metadata/verification/restore-to-new-db；事务 rollback 优先；不含调度、权限、远程存储、保留策略、KMS/TLS/UI。VP-040 修订为 v0.2.2，D-006/E-007 落在 workspace-040。 |

## 声明

本文件是已确认的 Vision Plan 意图，不是 Goal 五件套、实现事实或 progress 权威。激活不等于实现许可之外的 schema/迁移完成；R1/R2/R3 仍由工作区目标与 Goal 审计承接。

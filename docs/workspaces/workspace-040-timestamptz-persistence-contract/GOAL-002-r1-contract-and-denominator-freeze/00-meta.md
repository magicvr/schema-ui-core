---
id: GOAL-002-r1-contract-and-denominator-freeze
title: R1 · 时间合同与分母冻结
status: done
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 1.0.0
progress: 4/4
plan_refs:
  - VP-040-timestamptz-persistence-contract
primary_plan: VP-040-timestamptz-persistence-contract
serves_summary: 承接 Root R1：逐列盘点绝对时刻、冻结 PostgreSQL timestamptz(6) + SQLite 固定 6 位 RFC3339 TEXT 合同、NULL/零值与双方言原地转换边界，并在 self + grok independent 审计通过后为 R2 放行。
---

# GOAL-002 · R1 · 时间合同与分母冻结

## 概述

本子目标承接 `[workspace-040-timestamptz-persistence-contract] GOAL-001-timestamptz-persistence-contract` 的 R1 阶段。用户已在 2026-09-20 通过 P-004 选择：PostgreSQL 使用字面 `timestamptz(6)`；SQLite 使用 UTC RFC3339 固定 6 位小数的 `TEXT`；所有绝对时刻列进入分母；两方言各自原地转换；不提供 SQLite→PG 产品级搬运器；sentinel `0` 按语义转为 `NULL`。

本子目标只冻结合同与分母，不开始大规模 schema 迁移。R2 只能在本目标的 required 信息、self 审计与 grok independent 审计合法闭合后启动。

## 红线

- 不把 VP-013 历史 PG `BIGINT` / SQLite `INTEGER` 合同静默当作 VP-040 终态；它们是迁移来源。
- 不把所有整数列误归为时间；金额、flag、version、计数器、ID 前缀、duration、TOTP step 明确排除。
- 不在未冻结逐列单位、精度、NULL/默认值与转换失败策略前修改 migration DDL。
- 不引入 ORM、第三数据库、Redis/MQ/A3 或应用内 SQLite→PG 搬运器。

## 成功标准

- [ ] C1：全仓 compiled catalog 与运行时读写 inventory 完成；每列有旧单位、目标物理类型、精度、NULL/默认值、读写路径与证据。
- [ ] C2：合同冻结：PG `timestamptz(6)` + SQLite 固定 6 位 UTC RFC3339 `TEXT`；sentinel 0 → NULL 的逐列规则落盘；统一 6 位微秒 RFC3339 wire；未选方案与影响明确。
- [ ] C3：原地转换与失败/回滚策略冻结；不提供跨引擎产品级搬运器；统一 Backup SPI/Service（方言 native provider、metadata、verification、restore-to-new-db、失败边界）与 R2/R3 residual 明确。
- [ ] C4：Root self 审计 + 本地 grok build（grok 4.6 · high）independent 审计完成；required findings 合法闭合；Root R1 可标 completed。

## 纲领路线图

| 检查点 | 目的 | 状态 |
|---------|------|------|
| C1 | 逐列 inventory 与单位/类型/读写路径证据 | **completed**（A-006 independent accepted） |
| C2 | 物理合同、精度、NULL/零值、wire 与未选方案冻结 | **completed（frozen）**：90 列逐列合同、谓词 exact SQL、逐表 SQLite+PG DDL（44 张表）、可执行边界测试均落盘；**F-I-002 由 A-046 判 `fixed`** |
| C3 | 双方言原地转换、备份依赖与失败策略冻结 | **completed（frozen）**：`r1-c3-backup-recovery-boundary-v1.0-fc.md`（三类产物区分、双 token、调用点、harness、错误分类）；**F-I-004 经 A-042 closed、A-044/A-046 复审残留已修**；转换/失败策略经 `D-019`（F-5）与 `D-021` 落盘 |
| C4 | self + grok independent 审计、响应与 R1 放行 | **completed**：关门向 independent = **A-046**（**开放 required = 0**）；编排响应 = **A-047**；**用户 2026-09-20 书面确认关门** |

**`progress: 4/4`**（C1～C4 全部 completed；C4 的用户确认已于 2026-09-20 取得）。progress 不替代信息门禁或审计结论，**也不放行阶段**。

## 关门记录

- **关门日期**：2026-09-20（用户书面确认）
- **关门时开放 required**：**0**
- **关门向 independent 审计**：A-046（`conditional`，判定 F-I-002 `fixed`、F-I-005 `accepted-residual`，无新 required）
- **编排响应**：A-047（`pass`）
- **残余风险**：F-I-005 = **`accepted-residual`**（范围穷举三项 + 复审触发 + 失效条件，见 child `D-021`）——**不得读作哈希已验证**
- **仍开 recommended（不阻断关门，随 R2 处理）**：F-I-008、F-I-009（属 R3）、F-I-025（已修待复审确认）、F-I-028（已实际改写待复审确认）
- **放行边界**：本次关门**只放行设计面**；R2 的生产 schema 变更仍需 R2 自身验收（含 `D-021` 的 residual 复审触发点 = R2 首次记录任一 v73+ 哈希时）

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 | 影响门禁 | 状态 | 证据 |
|----|------|----------|----------|------|------|
| I-040-001 | required | PG `timestamptz(6)` 与 SQLite 固定 6 位 RFC3339 TEXT 的逐列编解码/精度合同 | C2/R2 | **verified** | 90 列逐列合同 `r1-c2-per-column-conversion-contract-v1.0-fc.md`；谓词 exact SQL；E1/E2 表达式经 SQLite 3.51/3.53 与 **PG 15.19/16.15/17.11** 实测；F-I-002 由 **A-046** 判 `fixed` |
| I-040-002 | required | 全部绝对时刻列分母与排除列清单 | C1/C2/R2 | **verified** | `r1-time-column-inventory-v0.3.md`：**90 列 / 44 张表**（「倒数第二段 = 表名」机械去重；A-044 独立复算一致）；排除列清单已列明；F-I-001/F-I-003 经 A-006/A-030 判 closed |
| I-040-003 | required | SQLite/PG 各自原地转换、失败恢复与备份依赖 | C3/R2/R3 | **verified**（R3 侧另有 residual） | `r1-c3-backup-recovery-boundary-v1.0-fc.md` 经 **A-042 判 F-I-004 closed**；转换/失败策略见 `D-019`（F-5 子女先行）与 `D-021`（PG 容器验证）；PG 跨版本兼容矩阵与 restore harness 的**执行**属 R3 |
| I-040-004 | required | VP-020 展示/输入回归矩阵 | R3 | **open**（**不阻断 R1**） | Root I-040-004；R1 仅登记接口，R3 执行；对应 recommended F-I-009 |

## 父目标

- `GOAL-001-timestamptz-persistence-contract`

## 关门条件

R1 子目标只有在 C1～C4 全部完成、Goal `03-audit` 的 self 与 independent 意见均已落盘、required finding 合法闭合后，才可由编排器静默关门并放行 Root R2。
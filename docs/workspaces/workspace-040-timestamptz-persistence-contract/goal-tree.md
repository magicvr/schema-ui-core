---
title: 目标树 · workspace-040-timestamptz-persistence-contract
status: active
created: 2026-09-20
updated: 2026-09-20
parent: null
version: 0.5.0
workspace_id: workspace-040-timestamptz-persistence-contract
---

# 目标树 · DB 时间列 timestamptz 持久化合同

> 工作区：`workspace-040-timestamptz-persistence-contract`
> canonical：`docs/workspaces/workspace-040-timestamptz-persistence-contract/`
> Root：`GOAL-001-timestamptz-persistence-contract`（**`active · 1/3`**）
> primary_plan：`VP-040-timestamptz-persistence-contract`（**`active`** v0.2.3）

## 树

```text
GOAL-001-timestamptz-persistence-contract [active] (1/3) · 纲领容器
├── GOAL-002-r1-contract-and-denominator-freeze [done] (4/4) · R1 合同与分母冻结
├── GOAL-003-r2-codec-and-descriptor-m1-m2 [done] (4/4) · R2 M1/M2 codec 与 descriptor
├── GOAL-004-r2-repository-and-predicate-rewrites [done] (3/3) · R2 M3 仓储与谓词改造
└── GOAL-005-r2-backup-port-and-closeout [active] (2/3) · R2 M4 Backup Port 与阶段关门
```

## 纲领路线图

```text
R1 合同与分母冻结 [completed · GOAL-002 · 4/4]
   → R2 双方言迁移 + Store 编解码 [active · 边界已冻结（Root D-016）· M1/M2 = GOAL-003、M3 = GOAL-004 · M2/M3 按 Root D-017 合并提交]
      → R3 读写/时区回归、备份核对、证据与关门 [pending]
```

> R1 已按用户裁决渐进建立为 `GOAL-002-r1-contract-and-denominator-freeze`（**`done · 4/4`**）；R2/R3 仍不预创建。Root `progress: 1/3` 由 Root `00-meta.md` 的三个显式检查点派生。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-timestamptz-persistence-contract | DB 时间列 timestamptz 持久化合同 | null | **active** | 1/3 | 2026-09-20 |
| GOAL-002-r1-contract-and-denominator-freeze | R1 · 时间合同与分母冻结 | GOAL-001-timestamptz-persistence-contract | **done** | 4/4 | 2026-09-20 |
| GOAL-003-r2-codec-and-descriptor-m1-m2 | R2 · 共享 codec 与 15 个 conversion descriptor（M1/M2） | GOAL-001-timestamptz-persistence-contract | **done** | 4/4 | 2026-09-20 |
| GOAL-004-r2-repository-and-predicate-rewrites | R2 · 仓储读写与谓词改造（M3） | GOAL-001-timestamptz-persistence-contract | **done** | 3/3 | 2026-09-20 |
| GOAL-005-r2-backup-port-and-closeout | R2 · Backup Port 与阶段关门（M4） | GOAL-001-timestamptz-persistence-contract | **active** | 2/3 | 2026-09-20 |

## 说明

- Root `active · 1/3`——**R1 已完成并关门**（2026-09-20 用户书面确认；关门向 independent = A-046，开放 required = 0）。
- **R1 关门边界**：只放行**设计面**。F-I-005 = **用户书面 `accepted-residual`**（范围 + 复审触发 + 失效条件见 child `D-021`），**不得读作哈希已验证**；R2 的生产 schema 变更仍须经 R2 自身验收（residual 复审触发 = R2 首次记录任一 v73+ 哈希时）。
- **M2/M3 次序**（**Root `D-017`**，用户 2026-09-20 裁决）：M2 落码后按 `D-018` 无过渡期使全仓测试红 → **M2 不单独提交**；M3 `GOAL-004-r2-repository-and-predicate-rewrites`（用户确认 slug）已立项并开工，转绿后与 M2 一并提交。`I-041-003` 已按常驻 PostgreSQL 15.4 实测证据关闭（约束：破坏性 migration 只可作用于一次性/专用测试 database）。M4 `GOAL-005-r2-backup-port-and-closeout` 已按用户确认 slug 立项（`active · 0/3`）：`D-016` 第 9 项 Backup Port 类型表面与 provider + `D-021` residual 收尾 + R2 关门审计。
- R2/R3 的门禁不得由激活状态或 progress 投影替代。
- 状态、progress、parent 或新增子目标发生变化时，必须同步本文件树与状态表。任何进度值都不放行阶段、不关闭 finding、不推导 `done`。

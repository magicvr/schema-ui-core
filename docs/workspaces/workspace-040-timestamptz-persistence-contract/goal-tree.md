---
title: 目标树 · workspace-040-timestamptz-persistence-contract
status: active
created: 2026-09-20
updated: 2026-09-20
parent: null
version: 0.1.0
workspace_id: workspace-040-timestamptz-persistence-contract
---

# 目标树 · DB 时间列 timestamptz 持久化合同

> 工作区：`workspace-040-timestamptz-persistence-contract`
> canonical：`docs/workspaces/workspace-040-timestamptz-persistence-contract/`
> Root：`GOAL-001-timestamptz-persistence-contract`（**`active · 0/3`**）
> primary_plan：`VP-040-timestamptz-persistence-contract`（**`active`** v0.2.0）

## 树

```text
GOAL-001-timestamptz-persistence-contract [active] (0/3) · 纲领容器
└── GOAL-002-r1-contract-and-denominator-freeze [active] (0/4) · R1 合同与分母冻结
```

## 纲领路线图

```text
R1 合同与分母冻结 [active · GOAL-002 · 0/4]
   → R2 双方言迁移 + Store 编解码 [pending]
      → R3 读写/时区回归、备份核对、证据与关门 [pending]
```

> R1 已按用户裁决渐进建立为 `GOAL-002-r1-contract-and-denominator-freeze`；R2/R3 仍不预创建。Root `progress: 0/3` 由 Root `00-meta.md` 的三个显式检查点派生。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-timestamptz-persistence-contract | DB 时间列 timestamptz 持久化合同 | null | **active** | 0/3 | 2026-09-20 |
| GOAL-002-r1-contract-and-denominator-freeze | R1 · 时间合同与分母冻结 | GOAL-001-timestamptz-persistence-contract | **active** | 0/4 | 2026-09-20 |

## 说明

- Root 初始 `active · 0/3`；激活与开区不代表任何 schema、迁移、编解码或回归已完成。
- R1 required 信息必须在方案冻结前关闭；R2/R3 的门禁不得由激活状态或 progress 投影替代。
- 状态、progress、parent 或新增子目标发生变化时，必须同步本文件树与状态表。

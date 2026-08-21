---
title: 目标树 · workspace-015-observability
status: active
created: 2026-08-21
updated: 2026-08-21
parent: null
version: 1.0.0
workspace_id: workspace-015-observability
---

# 目标树 · 可观测性

> 工作区：`workspace-015-observability`
> canonical：`docs/workspaces/workspace-015-observability/`
> Root：`GOAL-001-observability`（**active** · 1/5）
> primary_plan：`VP-015-observability`（**active** · 架构 A4）

## 树

```text
GOAL-001-observability [active]   · 可观测性（指标导出 + OpenTelemetry）
└── GOAL-002-metrics-export-contract [done 3/3] · R1 导出合同与配置面冻结（D-001 + 配置面；A-001 self pass）
```

R2～R5 子目标按纲领串行，待前序阶段收口后逐段立项。下一阶段：R2 指标 scrape 接入（依赖 R1 合同，已就绪）。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-observability | 可观测性（指标导出 + OpenTelemetry） | null | active | 1/5 | 2026-08-21 |
| GOAL-002-metrics-export-contract | R1 · 指标导出合同与配置面冻结 | GOAL-001-observability | done | 3/3 | 2026-08-21 |

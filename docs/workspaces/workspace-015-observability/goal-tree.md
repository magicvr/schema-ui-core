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
> Root：`GOAL-001-observability`（**active** · 3/5）
> primary_plan：`VP-015-observability`（**active** · 架构 A4）

## 树

```text
GOAL-001-observability [active]   · 可观测性（指标导出 + OpenTelemetry）
├── GOAL-002-metrics-export-contract [done 3/3] · R1 导出合同与配置面冻结（D-001 + 配置面；A-001 self pass）
├── GOAL-003-metrics-scrape-endpoint [done 4/4] · R2 指标 scrape 接入（internal/obs + live 冒烟；A-001 self pass）
├── GOAL-004-otel-traces [done 4/4]         · R3 OTel traces 接入（I-002 闭合；OTLP 导出实证；A-001 self pass）
└── GOAL-005-requestid-correlation [active] · R4 与 request-id 关联（D-001 已落盘，实现中）
```

R5 子目标待 R4 收口后立项。下一阶段：R5 双路径证据 + Root 关门（含独立审计）。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-observability | 可观测性（指标导出 + OpenTelemetry） | null | active | 3/5 | 2026-08-22 |
| GOAL-002-metrics-export-contract | R1 · 指标导出合同与配置面冻结 | GOAL-001-observability | done | 3/3 | 2026-08-21 |
| GOAL-003-metrics-scrape-endpoint | R2 · 指标 scrape 接入 | GOAL-001-observability | done | 4/4 | 2026-08-21 |
| GOAL-004-otel-traces | R3 · OpenTelemetry traces 接入 | GOAL-001-observability | done | 4/4 | 2026-08-22 |
| GOAL-005-requestid-correlation | R4 · 与 request-id 关联 | GOAL-001-observability | active | 0/4 | 2026-08-22 |

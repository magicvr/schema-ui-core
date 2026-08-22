---
title: 目标树 · workspace-015-observability
status: active
created: 2026-08-21
updated: 2026-08-22
parent: null
version: 1.1.0
workspace_id: workspace-015-observability
---

# 目标树 · 可观测性

> 工作区：`workspace-015-observability`
> canonical：`docs/workspaces/workspace-015-observability/`
> Root：`GOAL-001-observability`（**done** · 5/5）
> primary_plan：`VP-015-observability`（**closed** · 架构 A4；2026-08-22 组合层关门，VRev-034）

## 树

```text
GOAL-001-observability [done 5/5]   · 可观测性（指标导出 + OpenTelemetry）
├── GOAL-002-metrics-export-contract [done 3/3] · R1 导出合同与配置面冻结（D-001 + 配置面；A-001 self pass）
├── GOAL-003-metrics-scrape-endpoint [done 4/4] · R2 指标 scrape 接入（internal/obs + live 冒烟；A-001 self pass）
├── GOAL-004-otel-traces [done 4/4]         · R3 OTel traces 接入（I-002 闭合；OTLP 导出实证；A-001 self pass）
├── GOAL-005-requestid-correlation [done 4/4] · R4 与 request-id 关联（I-005 闭合；判据测试锁定；A-001 self pass）
└── GOAL-006-dual-path-evidence [done 4/4]    · R5 双路径证据与 Root 关门准备（E-002 live 证据；A-001 self pass）
```

Root 关门审计：A-001 self `pass` + A-002 independent `conditional`（F-001/F-002 → **fixed**；F-003 → 文档化残余；F-004/F-005 → **fixed**）→ A-003 响应闭环。VP-015 已于 2026-08-22 经 `/vision` 有界组合层 `closed`（VRev-034；VR-036）。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-observability | 可观测性（指标导出 + OpenTelemetry） | null | done | 5/5 | 2026-08-22 |
| GOAL-002-metrics-export-contract | R1 · 指标导出合同与配置面冻结 | GOAL-001-observability | done | 3/3 | 2026-08-21 |
| GOAL-003-metrics-scrape-endpoint | R2 · 指标 scrape 接入 | GOAL-001-observability | done | 4/4 | 2026-08-21 |
| GOAL-004-otel-traces | R3 · OpenTelemetry traces 接入 | GOAL-001-observability | done | 4/4 | 2026-08-22 |
| GOAL-005-requestid-correlation | R4 · 与 request-id 关联 | GOAL-001-observability | done | 4/4 | 2026-08-22 |
| GOAL-006-dual-path-evidence | R5 · 双路径证据与 Root 关门准备 | GOAL-001-observability | done | 4/4 | 2026-08-22 |
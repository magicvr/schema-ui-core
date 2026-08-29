---
title: 目标树 · workspace-024-distribution-formalization
status: active
created: 2026-08-29
updated: 2026-08-29
parent: null
version: 0.2.0
workspace_id: workspace-024-distribution-formalization
---

# 目标树 · 分发形态正式化（cli+包 对外服务化收口）

> 工作区：`workspace-024-distribution-formalization`（**active** · 2026-08-29 开区）
> canonical：`docs/workspaces/workspace-024-distribution-formalization/`
> Root：`GOAL-001-distribution-formalization`（**active** · 1/7）
> primary_plan：`VP-024-distribution-formalization`（active v0.2.0）

## 树

```text
GOAL-001-distribution-formalization [active 4/7]  · 分发形态正式化（cli+包 对外服务化收口：VP-022/023 go 后残余 7 项 + 方法 B 置顶）
├── GOAL-002-r1-serve-shell [done 5/5]  · R1 serve 壳闭环（schema-ui serve · HTTP 壳 + config 装载 + assembly 服务器面 · RT-D02 接线 · 判据 #1）
├── GOAL-003-r2-public-release-channel [done 4/4]  · R2 公开发布通道（npmjs @magicvr 六包 + apps/api/v0.4.0 + golden-field 无凭据消费 · 判据 #2）
├── GOAL-004-r3-compose-cicd [done 4/4]  · R3 compose/CI 实跑（compose 全服务 + consumer-regression 免凭据重构 + 信号级 drain harness · 判据 #3）
└── GOAL-005-r4-fork-comparison [done 4/4]  · R4 fork 对照计时（同一演进集 v0.3.0→v0.4.0：fork 同步 vs 包 bump 实测对比 · 判据 #4）
```

> 下一波：GOAL-006-r5-six-package-granularity（R5 六包形态细化 · done 5/5 · 判据 #5/#6）

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-distribution-formalization | 分发形态正式化（cli+包 对外服务化收口） | null | active | 4/7 | 2026-08-29 |
| GOAL-002-r1-serve-shell | R1 serve 壳闭环（schema-ui serve · HTTP 壳 + config 装载 + assembly 服务器面 · RT-D02 接线） | GOAL-001-distribution-formalization | done | 5/5 | 2026-08-29 |
| GOAL-003-r2-public-release-channel | R2 公开发布通道（npmjs @magicvr 六包 + apps/api/v0.4.0 + golden-field 无凭据消费） | GOAL-001-distribution-formalization | done | 4/4 | 2026-08-29 |
| GOAL-004-r3-compose-cicd | R3 compose/CI 实跑（compose 全服务 + consumer-regression 免凭据重构 + 信号级 drain harness） | GOAL-001-distribution-formalization | done | 4/4 | 2026-08-29 |
| GOAL-005-r4-fork-comparison | R4 fork 对照计时（同一演进集 v0.3.0→v0.4.0：fork 同步 vs 包 bump 实测对比） | GOAL-001-distribution-formalization | done | 4/4 | 2026-08-29 |
| GOAL-006-r5-six-package-granularity | R5 六包形态细化（renderer external 化 + ui 纯原子断言 + 冻结面 v1.4.0） | GOAL-001-distribution-formalization | active | 0/5 | 2026-08-29 |
---
title: 目标树 · workspace-014-object-storage
status: active
created: 2026-08-21
updated: 2026-08-21
parent: null
version: 0.4.0
workspace_id: workspace-014-object-storage
---

# 目标树 · 对象存储适配器

> 工作区：`workspace-014-object-storage`
> canonical：`docs/workspaces/workspace-014-object-storage/`
> Root：`GOAL-001-object-storage`（**active**；R1–R2 完成，R3–R5 未开始 → 计 2/5）
> primary_plan：`VP-014-object-storage`（**active** · 架构 A2）

## 树

```text
GOAL-001-object-storage [active]   · 对象存储适配器（S3 兼容 + 本地盘内嵌）
├── GOAL-002-object-port-freeze [done]   · R1 端口与配置面冻结
├── GOAL-003-object-s3-driver [done]   · R2 S3 兼容接入
└── GOAL-004-object-families-migration [active]   · R3 三类落盘收口走端口
```

R1/R2 已结项（各自 self + independent 审计通过，开放 required 均 0）。下一步按纲领串行进入 R3（三类落盘收口走端口），届时创建 GOAL-004。R3～R5 子目标按阶段渐进创建。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-object-storage | 对象存储适配器（S3 兼容 + 本地盘内嵌） | null | active | 2/5 | 2026-08-21 |
| GOAL-002-object-port-freeze | R1 内核对象存储端口与配置面冻结 | GOAL-001-object-storage | done | 4/4 | 2026-08-21 |
| GOAL-003-object-s3-driver | R2 S3 兼容接入（驱动 + readyz 扩依赖） | GOAL-001-object-storage | done | 3/3 | 2026-08-21 |
| GOAL-004-object-families-migration | R3 三类落盘收口走端口 | GOAL-001-object-storage | active | — | 2026-08-21 |

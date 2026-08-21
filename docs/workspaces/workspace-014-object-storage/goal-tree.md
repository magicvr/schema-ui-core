---
title: 目标树 · workspace-014-object-storage
status: active
created: 2026-08-21
updated: 2026-08-21
parent: null
version: 0.3.0
workspace_id: workspace-014-object-storage
---

# 目标树 · 对象存储适配器

> 工作区：`workspace-014-object-storage`
> canonical：`docs/workspaces/workspace-014-object-storage/`
> Root：`GOAL-001-object-storage`（**active**；R1 完成，R2–R5 未开始 → 计 1/5）
> primary_plan：`VP-014-object-storage`（**active** · 架构 A2）

## 树

```text
GOAL-001-object-storage [active]   · 对象存储适配器（S3 兼容 + 本地盘内嵌）
├── GOAL-002-object-port-freeze [done]   · R1 端口与配置面冻结
└── GOAL-003-object-s3-driver [active]   · R2 S3 兼容接入
```

R1 已由 GOAL-002 承载并结项（端口冻结 + 本地适配器 + 配置面；self A-001 pass、independent A-002 conditional→F-001 fixed 闭合，开放 required 0）。下一步按纲领串行进入 R2（S3 兼容接入），届时创建 GOAL-003。R2～R5 子目标按阶段渐进创建。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-object-storage | 对象存储适配器（S3 兼容 + 本地盘内嵌） | null | active | 1/5 | 2026-08-21 |
| GOAL-002-object-port-freeze | R1 内核对象存储端口与配置面冻结 | GOAL-001-object-storage | done | 4/4 | 2026-08-21 |
| GOAL-003-object-s3-driver | R2 S3 兼容接入（驱动 + readyz 扩依赖） | GOAL-001-object-storage | active | — | 2026-08-21 |

---
id: E-002
goal: GOAL-016-r3-s09-data-permission
title: S1 方案冻结完成（数据权限设计）
date: 2026-08-15
status: recorded
parent: GOAL-016-r3-s09-data-permission
created: 2026-08-15
updated: 2026-08-15
version: 1.0.0
---

# E-002 · S1 方案冻结完成（2026-08-15）

## 事实

- D-002 落盘：数据范围模型（all/self，org 按 I-004 排除）、过滤下推（resourceFilter→ResourceEntity.List 边界，RowScopeProvider 可选接口）、端点与权限键（data-permission.read/write）、迁移 0027/0028、协议对照（本地鉴权扩展）、S2 清单。
- I-001~I-004 全部闭合（证据 = D-002 各节）；progress 0/5 → 1/5（S1 检查点）。
- 关键证据：handler/resources.go（工厂 + requirePermission）、users_repository.go（SQL 组装）、kernel/profile.go（ProfileAdmin 内容扩展）、protocol-inventory D-PERM / ADR-0004。
- 审计：A-003 self（S1 方案）。

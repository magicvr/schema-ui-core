---
id: D-001
goal: GOAL-013-nav-order-config
title: 立项边界：导航顺序（方案 A）
date: 2026-08-14
status: accepted
parent: GOAL-013-nav-order-config
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-001 · 立项边界（导航顺序方案 A）

## 1. 背景与归属

- 用户 2026-08-14 裁决：方案 A（声明式集中排序）——默认排序由维护者在制作模块时排好（产品级冻结清单），使用者通过配置文件覆盖（重启生效）。业界主流 = 启动前配置写死（WordPress position / AntD 路由序 / Spree / Medusa），管理员页面排序为少数增强。
- 归属 workspace-11（导航结构是 admin 功能模块交付的一部分）；配置覆盖载体依赖 workspace-10 W7（YAML 配置体系）。

## 2. 现状（已实测）

- `sortNavigation`：Order 优先，冲突时 NodeID 字母兜底。
- 现有 Order 冲突：Order 1 = Users/Settings/Account；Order 2 = Roles/Activity/Notifications → 实际顺序 Dashboard → Account → Settings → Users → Activity → Notifications → Roles…，Users/Roles 被挤后，不符合直觉。

## 3. 关键设计约束（S1 冻结细化）

- 默认清单顺位（草案）：Dashboard → Users → Roles → Settings → Activity → Account → Notifications → File library → Data dictionary → System monitoring → Scheduled tasks → Recycle bin。
- 覆盖语义：全量清单（非差异补丁）；清单未列 NodeID 追加末尾（新模块不消失）。
- 非法配置：回退默认 + 告警日志（倾向，S1 确认）。
- 维护规则：新模块（S-05~S-14）落地时更新默认清单（测试锁定）。

## 4. 信息就绪

I-001/I-002/I-003（S1）、I-004（S4）见 00-meta；均 open。

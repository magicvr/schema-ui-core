---
id: E-001
goal: GOAL-013-nav-order-config
date: 2026-08-14
status: recorded
parent: GOAL-013-nav-order-config
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-001 · 立项（导航顺序）

## 事实

- 2026-08-14：用户裁决方案 A 立项（归 workspace-11）。业界查证：声明式/注册时写死顺序是主流（WordPress position、AntD Pro 路由序、Spree/Medusa 注册序）；管理员页面排序是少数增强（WP 插件、Atlassian pin）。
- 实测现状：Order 冲突清单（Users/Settings/Account=1、Roles/Activity/Notifications=2）→ 字母兜底导致顺序反直觉。
- 依赖：workspace-10 W7（YAML 配置体系）提供覆盖载体；默认清单部分不依赖，可先行。
- 五件套 + D-001 就位；信息项 I-001~I-004 登记（open）。
- 未产生代码变更。

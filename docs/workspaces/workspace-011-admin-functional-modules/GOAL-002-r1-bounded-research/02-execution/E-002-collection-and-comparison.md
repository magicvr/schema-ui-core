---
id: E-002
goal: GOAL-002-r1-bounded-research
title: S2/S3 · 候选池收集 + 基架对照
date: 2026-08-14
status: recorded
parent: GOAL-002-r1-bounded-research
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# E-002 · S2/S3 · 候选池收集 + 基架对照

## S2 · 候选池收集（来源证据）

- 业界通用 admin 能力样本：admin 面板综述（Appwrite）、电商 admin 参考（React e-commerce admin）、Go admin 框架（simple-admin-core）、企业后台模板惯例（vue-element-admin / Ant Design Pro 生态）、多商户平台面板（eshop-plus 类）。
- 业务领域样本：用户点名（订单 / 钱包典型，类目 / 通知入池）+ 业界电商后台惯例（商品、营销、物流等）。
- 来源锚点与 web 检索链接登记于 I-011-001 `2。

## S3 · 基架对照（已覆盖 C-01～C-11）

- 对照已交付基架（users/roles/settings/activity/operationlog 模块、代表页集、协议面 v2.8.0）→ 11 项已覆盖（C-01～C-11），不重复立项。
- 核实关键边界：改密 API 存在（users CRUD 内、吊销 access token）→ 个人中心缺口 = 自助 UI + 会话管理；无导出端点；overview 仅 demo（dev.examples）→ 生产 Profile 无 dashboard；无通知模块；上传仅控件级。

## 证据

- I-011-001-tiered-inventory.md（分档清单附件）；handler/users.go、account/session.go、模块名册、页面集（本会话核实）。

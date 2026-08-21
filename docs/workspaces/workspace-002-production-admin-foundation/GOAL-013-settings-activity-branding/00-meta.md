---
id: GOAL-013-settings-activity-branding
title: Settings 品牌与 Activity 操作日志只读面
status: done
created: 2026-08-04
updated: 2026-08-04
parent: GOAL-001-production-admin-foundation
version: 0.1.0
progress: 4/4
---

# GOAL-013 · Settings 品牌与 Activity 操作日志只读面

## 概述

用户裁决：恢复并**实装** Admin 必需的 Settings 与 Activity（操作日志）能力，而非删除导航。Settings 提供站点标题 + Logo URL（文本，无 upload）；Activity 提供操作日志只读列表与详情。

## 成功标准

- [x] **S1 · Settings API + 持久化**：`site_settings` 迁移；`GET /api/branding` 公开；`GET/PATCH /api/settings` 鉴权；标题必填默认 `Schema UI Core`；logo 可空。
- [x] **S2 · Branding 接线**：Shell/Login 显示标题；logo 空则不显示 UI logo 与 branding favicon；`document.title` 用站点标题。
- [x] **S3 · Activity 只读 API + Schema 页**：`/api/operations` list+detail、`operations.read`；settings/activity Schema + manifest 导航。
- [x] **S4 · 回归**：api `go test`/`vet` 全绿；web vitest 全绿 + `tsc -b`。

## 父目标

- [GOAL-001-production-admin-foundation](../GOAL-001-production-admin-foundation/00-meta.md)

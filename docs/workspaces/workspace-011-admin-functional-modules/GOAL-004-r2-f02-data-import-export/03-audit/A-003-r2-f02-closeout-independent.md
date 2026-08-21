---
id: A-003
goal: GOAL-004-r2-f02-data-import-export
title: S5 · 关门 independent 审计（grok-4.6 · security/data · conditional）
date: 2026-08-14
source: independent
scope: S5 关门（提交 `39a1671` 声称范围；security/data 高影响门禁）
verdict: conditional
parent: GOAL-004-r2-f02-data-import-export
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# A-003 · S5 关门 independent 审计（grok build 代贴）

> provider：grok-4.6（本地 CLI · 只读）。原意见全文见会话记录；本文件为代贴摘要 + 完整 findings 台账。

## verdict：**conditional**

导入写路径绕过既有角色委派边界（F-001 **required**）——必须与 `POST /api/users` 相同的 `authorizeRoleAssignment` 边界；其余为 recommended/info。

## findings（完整台账）

| id | severity | scope | 内容 | 编排器处置 |
|----|----------|-------|------|------------|
| F-001 | **required** | 导入委派边界 | 导入写路径绕过 `authorizeRoleAssignment`：持有 `data.import` 的非 admin 可导入 `roles=admin` 或超持权限角色 | **fixed**（D-004 / E-004）：行级委派检查 + 用例 |
| F-002 | recommended | 结构坏行 | 引号不闭合等行级解析错 → 硬失败无报告无审计 | **fixed**：收集进 errors[] + 停止解析 + applied>0 必审计 |
| F-003 | recommended | BOM | 导入不剥离 UTF-8 BOM（Excel 首列 `﻿username` → 整表失败） | **fixed**：表头前剥离 |
| F-004 | recommended | 未知列 | 未知列未计入告警 | 留痕（告警列表归 R3；D-002 更新） |
| F-005 | recommended | 测试缺口 | 匿名 401 / editor 导出 200 / data.export 审计 / RFC 4180 转义 / 413 / INVALID_CSV / 委派用例 | **fixed**：补 8 用例 |
| F-006 | recommended | 错误契约 | 导出超限仍用 `INVALID_PAGE_SIZE`（目录文案 100，导出上限 10000） | **fixed**：新增 `INVALID_EXPORT_LIMIT` 码 |
| F-007 | info | 审计健壮性 | 导出审计吞写失败 | **fixed**：slog 记录 |
| F-008 | info | 外籍文件口径 | 方案写 404，实现 403（403 语义更准） | D-002 口径修正 |
| F-009 | recommended | 公式注入 | 导出用户可控字段未防 `= + - @` 公式注入 | **fixed**：前导 `'` 防护 + 用例 |
| F-010 | info | 只读限制 | 未复跑测试/未核 SHA | E-004 复跑 |

## 审计问题对照（8/8）

1. 权限 fail-closed：成立（匿名 401、无键 403、模块关闭 404）；PolicyAdminEditor 合理（读类数据面）。
2. fileId owner 校验：成立（外籍 403、缺失 404）；类型/大小服务端独立判定。
3. 行级校验 + 不回滚：主路径成立；**写路径缺委派边界（F-001）**；结构坏行硬失败（F-002）。
4. 导出字段：列白名单固定（users/roles 显式字段清单，无 password_hash）。
5. CustomAction 白名单：未知 handler fail-closed；URL 仅白名单；文件名由 handler 推导（无注入）。
6. 迁移 0015：CHECK 超集重建，账本 15 条冻结校验通过。
7. 操作日志：import 有测试；导出无单测（F-005 fixed）。
8. 错误码契约：+3 码已入 frozen 集 + catalog（F-006 再 +1）。

## 建议

编排器：F-001 以 **fixed** 闭合（D-004），复跑全量测试后关门。

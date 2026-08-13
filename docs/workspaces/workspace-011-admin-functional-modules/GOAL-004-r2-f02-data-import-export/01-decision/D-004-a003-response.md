---
id: D-004
goal: GOAL-004-r2-f02-data-import-export
title: A-003 响应 — F-001 required fixed + recommended 全落地（无 P-004 裁决冲突）
date: 2026-08-14
status: accepted
parent: GOAL-004-r2-f02-data-import-export
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# D-004 · A-003（grok independent · conditional）响应

> P-004 检查：F-001 为 required 且与 self（A-002「无 required」）不一致。本响应选择 **fixed**——委派边界是既有关卡（GOAL-011 `7.2），导入不应成为旁路；修复成本低、可核对。

## 处置

| finding | 处置 | 修复 |
|---------|------|------|
| F-001 required | **fixed** | `importRoleAssignmentError`：行级 `roles.assign` / admin 仅限 admin / 不得授予超持权限角色；未知角色键下放 `CreateUserManagement` 报 `INVALID_ROLE_REF`；用例（import-only 角色导入 `roles=admin` → 行失败，另一行成功） |
| F-002 recommended | **fixed** | 结构坏行 → errors[] + 停止解析；applied>0 必写审计 |
| F-003 recommended | **fixed** | 表头前剥离 UTF-8 BOM |
| F-004 recommended | 留痕 | 未知列告警列表归 R3（D-002 已注明） |
| F-005 recommended | **fixed** | 补 8 用例（匿名 401 / editor 导出 200 / export 审计 / RFC 4180 / 公式注入 / 413 / INVALID_CSV / 委派） |
| F-006 recommended | **fixed** | 新增 `INVALID_EXPORT_LIMIT` 错误码（契约 + catalog） |
| F-007 info | **fixed** | 导出审计失败 slog 记录 |
| F-008 info | 留痕 | D-002 口径修正（外籍 403） |
| F-009 recommended | **fixed** | 导出公式注入防护（`= + - @` 前导 `'`）+ 用例 |
| F-010 info | 留痕 | E-004 全量复跑 |

## 验证

`go test ./... -count=1` 复跑全绿后关门（E-004）。

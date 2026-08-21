---
id: D-004
goal: GOAL-003-r2-f01-dashboard
title: A-003 响应 — F-001 required fixed + recommended 全落地（无 P-004 裁决冲突）
date: 2026-08-14
status: accepted
parent: GOAL-003-r2-f01-dashboard
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# D-004 · A-003（grok independent · conditional）响应

> P-004 检查：F-001 为 required 且 self（A-002）未发现。选择 **fixed**——渲染路径丢键是真实缺陷（dashboard 为生产 home，双语是 VP-007 已承诺面），修复成本低、可核对。

## 处置

| finding | 处置 | 修复 |
|---------|------|------|
| F-001 required | **fixed** | `parseRenderNode`：text 保留 `textKey`、statCard 保留 `labelKey`（类型 + 解析 + 单测）——data-display 同源受益 |
| F-002 | **fixed** | e2e 4 文件 home 断言改 dashboard（demo overview 不变） |
| F-003 | **fixed** | 结构/i18n 分母 += dashboard.json + fragment |
| F-004 | **fixed** | SM-007 增加 `homePageRef` 断言 |
| F-005 | **fixed** | `menu_dashboard` Permission 留空（policy 授予） |
| F-006/F-007/F-009 | 留痕 | E-003/D-002 措辞 |
| F-008 | **fixed** | StaticDevSession `menu_dashboard: true` |

## 验证

`go test ./... -count=1` + `npm test` 复跑全绿后关门（E-004）。

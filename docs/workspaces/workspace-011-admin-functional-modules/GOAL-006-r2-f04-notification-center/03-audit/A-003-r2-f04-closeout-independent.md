---
id: A-003
goal: GOAL-006-r2-f04-notification-center
title: S5 · 关门 independent 审计（grok-4.6 · security/data · conditional）
date: 2026-08-14
source: independent
scope: S5 关门（提交 `a065288` 声称范围；security/data 面）
verdict: conditional
parent: GOAL-006-r2-f04-notification-center
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# A-003 · S5 关门 independent 审计（grok build 代贴）

> provider：grok-4.6（本地 CLI · 只读）。原意见全文见会话记录；本文件为代贴摘要 + 完整 findings 台账。

## verdict：**conditional**

安全/数据主路径（owner 隔离、四类钩子、裁剪、主开关 API、已读模型、迁移账本、错误码、铃铛 fail-open）成立；**F-001 required**：范例页「通知设置」写不进去（schema `select` 提交字符串，API 只收 JSON bool，且无回读）。

## findings（完整台账）

| id | severity | scope | 内容 | 编排器处置 |
|----|----------|-------|------|------------|
| F-001 | **required** | 设置页 | `select` 提交字符串 vs API 只收 bool；无 GET settings 回读 | **fixed**（D-004 / E-004）：GET/PATCH settings（兼收 bool 与 `"true"/"false"` 字符串，表单面返回字符串）+ recordSource 回填 + 用例 |
| F-002 | recommended | 事件幂等 | 再 disable 已停用 / unlock 未锁定仍发通知 | **fixed**：仅真实转移发（handler 预读状态）+ 用例（re-disable 不重复） |
| F-003 | recommended | best-effort | 通知失败不打日志 | **fixed**：`slog.Error`（与 operationlog 同纪律） |
| F-004 | recommended | 测试缺口 | 缺跨用户 404 / 开关抑制钩子 / 锁内第 6 次失败不重复 | **fixed**：补 4 用例 |
| F-005 | recommended | 铃铛 a11y | 徽标 aria-hidden；失败文案误用 `feedback.actionCompleted` | **fixed**：aria-label 带未读数；`shell.notifications.unavailable` 独立文案 |
| F-006 | info | 模块关仍可落行 | 事件钩子不依赖模块启用（迁移全局） | 留痕（E-003/D-002：行可能落库但 UI 面不可用，可接受） |

## 审计问题对照（8/8）

1. owner 边界 ✅（外籍/未知 id → 404 无 oracle；F-004 补跨用户真实 id 用例）。
2. 事件钩子 ✅（触发点正确；F-002 补转移幂等；F-003 补日志）。
3. 裁剪 ✅（500+2 保未读有覆盖）。
4. 主开关 ✅（API/钩子成立；F-001 修页面写入路径后闭环）。
5. 已读模型 ✅（幂等 read、read-all 计数）。
6. 迁移 0016/0017 ✅（FK/CHECK/索引；账本 17 条）。
7. 铃铛 ✅（fail-open、React 转义无 XSS；F-005 补 a11y/文案）。
8. 错误码 ✅（+2 码入契约 + catalog）。

## 建议

编排器：F-001 以 **fixed** 闭合（D-004），复跑全量测试后关门。

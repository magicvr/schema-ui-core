---
doc_type: vision-roadmap
title: 愿景组合编排
status: active
created: 2026-07-31
updated: 2026-07-31
parent: null
version: 0.2.0
---

# 组合编排 · Schema UI Core Admin 基架

本文件索引已落盘的 VP 与用户确认的后续方向；它不是 Goal 路线图，也不汇总 progress%。

## 已落盘意图

| 顺序 | VP | 意图 | 前置 | 状态 |
|------|----|------|------|------|
| 1 | [VP-001-mvp-admin-foundation](plans/VP-001-mvp-admin-foundation.md) | 初始化 React + Go Admin MVP，覆盖固定协议来源、核心账号权限与协议范例验证。 | 无 | active（lead: workspace-001-mvp-admin-foundation） |

## 已确认的后续方向（尚未建立 VP）

| 顺序 | 方向 | 与前序关系 | 建立 VP 前的约束 |
|------|------|------------|------------------|
| 2 | 前端产品化：以 Linear + Vercel Dashboard 为参考，使用 Tailwind CSS 与 shadcn/ui 风格组件，支持浅色 / 深色模式。 | 在 MVP 基架可用后推进。 | 需确认其独立退出判据与是否复用 VP-001 的工作区。 |
| 3 | 双线演进：维护最小 MVP 基架分支与较完整的 Admin 功能分支，后者逐步纳入钱包、订单、类目、通知等通用能力。 | 以可 fork 的 MVP 基线为前提。 | 需明确分支命名、兼容承诺、回合并策略和首批功能范围。 |

后续方向在建立对应 `VP-00N-*.md` 前不是可引用的 `primary_plan`。

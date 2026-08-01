---
doc_type: vision-roadmap
title: 愿景组合编排
status: active
created: 2026-07-31
updated: 2026-08-01
parent: null
version: 0.4.0
---

# 组合编排 · Schema UI Core Admin 基架

本文件索引已落盘的 VP 与用户确认的后续方向；它不是 Goal 路线图，也不汇总 progress%。

## 已落盘意图

| 顺序 | VP | 意图 | 前置 | 状态 |
|------|----|------|------|------|
| 1 | [VP-001-mvp-admin-foundation](plans/VP-001-mvp-admin-foundation.md) | 初始化 React + Go Admin MVP，覆盖固定协议来源、核心账号权限与协议范例验证。 | 无 | **closed**（2026-08-01；lead: workspace-001-mvp-admin-foundation；三条退出判据经 R6 工作区 Q2 证据满足，用户确认关门） |
| 2 | [VP-002-production-admin-foundation](plans/VP-002-production-admin-foundation.md) | 在 I-PROTO-001 冻结子集之上，交付可直接 fork 使用的生产级 Schema 驱动 Admin 基架：Renderer、真实认证、持久化权限、CRUD 与工程化启动。 | 继承 VP-001 协议验证基线 | **planned**（暂未绑定工作区） |

## 已确认但尚未纳入新 VP 的后续方向

| 顺序 | 方向 | 与前序关系 | 建立 VP 前的约束 |
|------|------|------------|------------------|
| 3 | 订单、钱包、类目、通知等业务能力 | 以 VP-002 的 Admin 基架为前提。 | 需另行建立 VP，明确业务范围与独立退出判据。 |

VP-002 已吸收此前“前端产品化”和“双线演进”方向中的基架部分；剩余业务能力在建立对应 `VP-00N-*.md` 前不是可引用的 `primary_plan`。

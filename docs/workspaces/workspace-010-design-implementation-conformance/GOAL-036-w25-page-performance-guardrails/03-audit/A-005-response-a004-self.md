---
title: A-005 · 响应 A-004 独立复核意见（self · 编排响应）
source: self
status: recorded
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-036-w25-page-performance-guardrails
version: 0.1.0
scope: 响应 A-004（independent · pass · F-009 / F-010 recommended）
verdict: pass（响应结论：recommended 全闭）
---

# A-005 · 响应 A-004 独立复核意见（2026-08-23，self）

响应对象：`A-004-correction-recheck-independent.md`（ox-alpha /audit，`pass`——全部修正关闭证据独立复跑名实相符）。采纳 F-009/F-010 两条 recommended，均按 `fixed` 处置；原 verdict 与原文不改写。

## 关闭证据表

| Finding | 状态 | 证据路径 |
|---------|------|----------|
| **F-009（recommended · low）** GOAL-036 三个索引文件（01-decision / 02-execution / 03-audit）frontmatter 仍 `status: active`，与 00-meta `done` 不一致 | **fixed** | 三文件 frontmatter 统一为 `status: done`（version 顺延 0.2.0/0.2.0/0.4.0）；与 GOAL-037 姊妹索引（done）实践一致 |
| **F-010（recommended · low-med）** `TestSQLiteDSNPragmas` 未钉 `_txlock=immediate`，该连接面不变量仅有概率性网络 | **fixed** | `store_wal_test.go` DSN want 列表补 `_txlock=immediate`（含注释引用 F-010）；`go test ./internal/store/ -run TestSQLiteDSNPragmas` ok——与 F-002 同精神：连接面五参数全钉死 |

## 备注观察回应（A-004/A-002 记录，非 finding）

- `newID` 契约产品/替身双处手工镜像（A-002 备注①）：维持现状，现有排序单测（`TestNewIDSameMillisecondOrdering` 等）为概率性网络；若未来再改 id 格式应同步两处——本响应记录为长期维护提示，不立项。
- 浏览器 e2e 可选背书（A-004 建议 2）：机制层已被单测栅栏覆盖，非必需；如用户需要可另约复跑。

## 仍开放项

无（两项 recommended 已 fixed；无 required；不重开门禁）。

## 结论

A-004（pass）的两条 recommended 已按 `fixed` 合法闭合；GOAL-036 维持 `done 6/6` 不变。
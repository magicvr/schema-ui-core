---
id: GOAL-005-requestid-correlation
title: R4 · 与 request-id 关联
status: done
parent: GOAL-001-observability
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
progress: 4/4
---

# GOAL-005 · R4 与 request-id 关联

## 概述

承载 Root [GOAL-001-observability](../GOAL-001-observability/00-meta.md) 纲领阶段 **R4**：闭合 I-005（request-id / correlation 如何写入 span：属性名、是否 baggage），让 VP-015 退出判据 2 的「能与现有 request-id / correlation 关联」可核对。范围仅限 span 属性与 baggage 注入；不重开 VP-012 的 correlation/错误包络。

## 成功标准（检查点）

- [x] D-001 闭合 I-005：属性名冻结（`correlation.request_id`）、baggage 键（`request-id`）与注入方式、关联判据可核对（checkpoint `8b52f2d`）
- [x] Wrap 落点：span 带 `correlation.request_id` 且与请求上下文一致（真实 requestid 中间件链测试锁定）；baggage 注入可核对；无 id 时静默跳过（checkpoint `bc5e196`）
- [x] 测试覆盖（传入 id / 生成 id 路径共享同一断言 / baggage / 无 id）+ vet/build/全仓 test 全绿（checkpoint `bc5e196`）
- [x] 自审 A-001 pass 后关门（A-001 self pass，开放 required = 0）

`progress` = 完成检查点数 / 4。当前 **4/4**。关门审计：A-001（self，pass）。

## 信息就绪与未知项

I-005（required，最晚 R4 接入前）由本目标 D-001 关闭。无其他新增未知。

## 父目标

- [GOAL-001-observability](../GOAL-001-observability/00-meta.md)

## 台账布局

三个可追加台账使用同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`。
---
id: GOAL-005-requestid-correlation
title: R4 · 与 request-id 关联
status: active
parent: GOAL-001-observability
created: 2026-08-22
updated: 2026-08-22
version: 0.1.0
progress: 0/4
---

# GOAL-005 · R4 与 request-id 关联

## 概述

承载 Root [GOAL-001-observability](../GOAL-001-observability/00-meta.md) 纲领阶段 **R4**：闭合 I-005（request-id / correlation 如何写入 span：属性名、是否 baggage），让 VP-015 退出判据 2 的「能与现有 request-id / correlation 关联」可核对。范围仅限 span 属性与 baggage 注入；不重开 VP-012 的 correlation/错误包络。

## 成功标准（检查点）

- [ ] D-001 闭合 I-005：属性名冻结、baggage 键与注入方式、关联判据
- [ ] Wrap 落点：span 带 `correlation.request_id` 且与请求上下文一致；baggage 注入可核对；nil/无效请求 id 行为确定
- [ ] 测试覆盖（传入 id / 生成 id / baggage 提取）+ 全仓 vet/build/test 全绿
- [ ] 自审 A-001 pass 后关门

`progress` = 完成检查点数 / 4。当前 0/4。

## 信息就绪与未知项

I-005（required，最晚 R4 接入前）由本目标 D-001 关闭。无其他新增未知。

## 父目标

- [GOAL-001-observability](../GOAL-001-observability/00-meta.md)

## 台账布局

三个可追加台账使用同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`。
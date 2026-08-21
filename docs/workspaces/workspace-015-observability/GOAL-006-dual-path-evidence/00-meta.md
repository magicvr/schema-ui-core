---
id: GOAL-006-dual-path-evidence
title: R5 · 双路径证据与 Root 关门准备
status: active
parent: GOAL-001-observability
created: 2026-08-22
updated: 2026-08-22
version: 0.1.0
progress: 0/4
---

# GOAL-006 · R5 双路径证据与 Root 关门准备

## 概述

承载 Root [GOAL-001-observability](../GOAL-001-observability/00-meta.md) 纲领阶段 **R5**：收集并落盘 VP-015 退出判据 3/4 的可核对证据——(a) 缺省（无收集器）仍能开发快测、无额外监听；(b) 显式配置下 metrics scrape **与** 至少一条 trace 导出都有可核对证据。同时为 Root 关门提供独立审计入口（grok build /audit）。

## 成功标准（检查点）

- [ ] D-001 证据方案：判据映射、收集步骤、判定标准、工具（otlp-sink）
- [ ] 缺省路径证据：无 observability 配置启动 → 服务可用、无 metrics/traces 额外端口监听
- [ ] 显式路径证据：metrics scrape 系列实测 + 真实 OTLP sink 收到 trace 导出；命令可重复（N-004/N-008）
- [ ] 自审 A-001 pass 后关门

`progress` = 完成检查点数 / 4。当前 0/4。

## 信息就绪与未知项

无新增开放信息项（I-001～I-005 已全部 verified）。实施未知：无。

## 父目标

- [GOAL-001-observability](../GOAL-001-observability/00-meta.md)

## 台账布局

三个可追加台账使用同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`。
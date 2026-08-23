---
id: GOAL-006-dual-path-evidence
title: R5 · 双路径证据与 Root 关门准备
status: done
parent: GOAL-001-observability
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
progress: 4/4
---

# GOAL-006 · R5 双路径证据与 Root 关门准备

## 概述

承载 Root [GOAL-001-observability](../GOAL-001-observability/00-meta.md) 纲领阶段 **R5**：收集并落盘 VP-015 退出判据 3/4 的可核对证据——(a) 缺省（无收集器）仍能开发快测、无额外监听；(b) 显式配置下 metrics scrape **与** 至少一条 trace 导出都有可核对证据。同时为 Root 关门提供独立审计入口（grok build /audit）。

## 成功标准（检查点）

- [x] D-001 证据方案：判据映射、收集步骤、判定标准、工具 otlp-sink（checkpoint `8ddbb60`；工具 `cf9df6c`）
- [x] 缺省路径证据：无 observability 配置启动 → healthz/readyz 200、25081/4318 无监听、启动日志零提及（E-002）
- [x] 显式路径证据：metrics scrape 系列实测 + 真实 OTLP sink 收到 1037 字节 protobuf 导出；命令序列可重复（E-002）
- [x] 自审 A-001 pass 后关门（A-001 self pass，开放 required = 0）

`progress` = 完成检查点数 / 4。当前 **4/4**。关门审计：A-001（self，pass）。

## 信息就绪与未知项

无新增开放信息项（I-001～I-005 已全部 verified）。实施未知：无。

## 父目标

- [GOAL-001-observability](../GOAL-001-observability/00-meta.md)

## 台账布局

三个可追加台账使用同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`。
---
id: GOAL-004-otel-traces
title: R3 · OpenTelemetry traces 接入
status: active
parent: GOAL-001-observability
created: 2026-08-21
updated: 2026-08-21
version: 0.1.0
progress: 0/4
---

# GOAL-004 · R3 OpenTelemetry traces 接入

## 概述

承载 Root [GOAL-001-observability](../GOAL-001-observability/00-meta.md) 纲领阶段 **R3**：闭合 I-002（OTLP 协议、采样默认、未配置 endpoint 的 no-op 语义），落地 tracer provider 生命周期与 HTTP server span（经 R2 已有的 InstrumentedMux 拦截点），未显式配置不得影响 mvp/dev。request-id 关联归 R4。

## 成功标准（检查点）

- [ ] D-001 闭合 I-002：协议选型、采样默认、no-op 语义、span 面、配置键
- [ ] `observability.traces.{enabled,endpoint,sample_ratio}` 进入 Config / 两份 YAML / env 映射，fail-closed 校验有测试覆盖
- [ ] obs tracing 落地：span 创建/属性/状态映射/采样、OTLP/HTTP 导出路径实测（httptest sink 收到 POST）、缺省 no-op 零行为变化；composition 接线 + gofmt/vet/build/test 全绿
- [ ] 自审 A-001 pass 后关门

`progress` = 完成检查点数 / 4。当前 0/4。

## 信息就绪与未知项

本目标的实施前提即 I-002（required，最晚需要阶段 = R3 方案冻结）——由本目标 D-001 关闭。I-005（request-id ↔ span 关联）保持 open，归 R4。

## 父目标

- [GOAL-001-observability](../GOAL-001-observability/00-meta.md)

## 台账布局

三个可追加台账使用同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`。

---
id: GOAL-004-otel-traces
title: R3 · OpenTelemetry traces 接入
status: done
parent: GOAL-001-observability
created: 2026-08-21
updated: 2026-08-22
version: 1.0.0
progress: 4/4
---

# GOAL-004 · R3 OpenTelemetry traces 接入

## 概述

承载 Root [GOAL-001-observability](../GOAL-001-observability/00-meta.md) 纲领阶段 **R3**：闭合 I-002（OTLP 协议、采样默认、未配置 endpoint 的 no-op 语义），落地 tracer provider 生命周期与 HTTP server span（经 R2 已有的 InstrumentedMux 拦截点），未显式配置不得影响 mvp/dev。request-id 关联归 R4。

## 成功标准（检查点）

- [x] D-001 闭合 I-002：OTLP/HTTP 选型、ParentBased+ratio 采样、no-op 语义、span 面、三配置键（checkpoint `0470307`）
- [x] `observability.traces.{enabled,endpoint,sample_ratio}` 进入 Config / 两份 YAML / env 映射，fail-closed 校验有测试覆盖（10 子测试，checkpoint `2ab4ec4`）
- [x] obs tracing 落地：span 形状/属性/状态映射/采样、W3C 传播、OTLP sink 实证交付；live 冒烟（不可达 endpoint 启动/服务正常 + 导出失败仅 WARN）；composition 接线 + vet/build/test 全绿（checkpoint `2ab4ec4`）
- [x] 自审 A-001 pass 后关门（A-001 self pass，开放 required = 0）

`progress` = 完成检查点数 / 4。当前 **4/4**。关门审计：A-001（self，pass）。

## 信息就绪与未知项

本目标的实施前提即 I-002（required，最晚需要阶段 = R3 方案冻结）——由本目标 D-001 关闭。I-005（request-id ↔ span 关联）保持 open，归 R4。

## 父目标

- [GOAL-001-observability](../GOAL-001-observability/00-meta.md)

## 台账布局

三个可追加台账使用同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`。

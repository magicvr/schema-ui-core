---
id: GOAL-003-metrics-scrape-endpoint
title: R2 · 指标 scrape 接入
status: active
parent: GOAL-001-observability
created: 2026-08-21
updated: 2026-08-21
version: 0.1.0
progress: 0/4
---

# GOAL-003 · R2 指标 scrape 接入

## 概述

承载 Root [GOAL-001-observability](../GOAL-001-observability/00-meta.md) 纲领阶段 **R2**：按 [GOAL-002 D-001](../GOAL-002-metrics-export-contract/01-decision/D-001-metrics-export-contract.md) 已冻结的合同，落地专用 metrics listener 与内核系列——`internal/obs` 包（registry + instrumentation + listener 生命周期）、composition 接线（InstrumentedMux + fx lifecycle）、`module_id` 所有权标注。缺省仍全关；本阶段不触碰 traces。

## 成功标准（检查点）

- [ ] D-001 冻结实施接缝：包边界、instrumentation 拦截点、所有权规则、生命周期与失败语义、token 校验方式
- [ ] `internal/obs` 落地：固定系列（build_info / http_requests_total / duration / modules_enabled / Go+process collectors）、Bearer 守卫、route 标签取注册 pattern；单测覆盖
- [ ] composition 接线：enabled=true 可真实启动并 scrape（`suc_build_info` 等可见），contributed 路由携带模块 `module_id`，disabled 缺省零行为变化；gofmt/vet/build/test 全绿
- [ ] 自审 A-001 pass 后关门（3 前置全过）

`progress` = 完成检查点数 / 4。当前 0/4。

## 信息就绪与未知项

无新增信息项。继承 Root I-001/I-003/I-004（已 verified）；A-001 N-002 建议（route 取 pattern、所有权取 ContributionIdentity、不影响 readyz）作为本目标输入。

## 父目标

- [GOAL-001-observability](../GOAL-001-observability/00-meta.md)

## 台账布局

三个可追加台账使用同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`。

---
id: GOAL-003-metrics-scrape-endpoint
title: R2 · 指标 scrape 接入
status: done
parent: GOAL-001-observability
created: 2026-08-21
updated: 2026-08-21
version: 1.0.0
progress: 4/4
---

# GOAL-003 · R2 指标 scrape 接入

## 概述

承载 Root [GOAL-001-observability](../GOAL-001-observability/00-meta.md) 纲领阶段 **R2**：按 [GOAL-002 D-001](../GOAL-002-metrics-export-contract/01-decision/D-001-metrics-export-contract.md) 已冻结的合同，落地专用 metrics listener 与内核系列——`internal/obs` 包（registry + instrumentation + listener 生命周期）、composition 接线（InstrumentedMux + fx lifecycle）、`module_id` 所有权标注。缺省仍全关；本阶段不触碰 traces。

## 成功标准（检查点）

- [x] D-001 冻结实施接缝：包边界、instrumentation 拦截点、所有权规则、生命周期与失败语义、token 校验方式（checkpoint `ef33b40`）
- [x] `internal/obs` 落地：固定系列、Bearer 守卫、route 标签取注册 pattern；7 组单测覆盖（checkpoint `5ba04c5`）
- [x] composition 接线：enabled=true 真实启动并 scrape（live 冒烟实测 suc_* 全系列），contributed 路由携带模块 `module_id`（probe=admin.probe / health=core 实证），disabled 缺省零行为变化；gofmt/vet/build/test 全绿（checkpoint `5ba04c5`）
- [x] 自审 A-001 pass 后关门（A-001 self pass，开放 required = 0）

`progress` = 完成检查点数 / 4。当前 **4/4**。关门审计：A-001（self，pass）。

## 信息就绪与未知项

无新增信息项。继承 Root I-001/I-003/I-004（已 verified）；A-001 N-002 建议（route 取 pattern、所有权取 ContributionIdentity、不影响 readyz）作为本目标输入。

## 父目标

- [GOAL-001-observability](../GOAL-001-observability/00-meta.md)

## 台账布局

三个可追加台账使用同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`。

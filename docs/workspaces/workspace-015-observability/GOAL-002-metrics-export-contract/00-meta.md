---
id: GOAL-002-metrics-export-contract
title: R1 · 指标导出合同与配置面冻结
status: done
parent: GOAL-001-observability
created: 2026-08-21
updated: 2026-08-21
version: 1.0.0
progress: 3/3
---

# GOAL-002 · R1 指标导出合同与配置面冻结

## 概述

承载 Root [GOAL-001-observability](../GOAL-001-observability/00-meta.md) 纲领阶段 **R1**：把 VP-015 的 I-015-001 / I-015-003 收敛为可实施的导出合同（暴露面形态、绑定/鉴权、系列与标签契约、基数边界、本波分母），并把 `observability.metrics.*` 配置面落进 config 加载器（YAML + env 插值 + fail-closed 校验 + 测试）。

指标 scrape 接入本体（listener、registry、instrumentation）归 R2 子目标；本目标不改 HTTP 行为，缺省零变化。

## 成功标准（检查点）

- [x] D-001 冻结导出合同：暴露面形态、鉴权、内核最小系列、标签白名单与秘密禁令、基数规则、本波分母（Store/对象存储/Job 出局）、readyz 边界（checkpoint `499f97d`）
- [x] `observability.metrics.{enabled,addr,auth_token}` 进入 Config / 两份 YAML / env 映射，fail-closed 校验有测试覆盖（checkpoint `45489f4`，13 子测试）
- [x] `go build ./...` 与 `go test ./internal/config/...`（及全仓无 FAIL）通过；gofmt 干净（本切片文件）

`progress` = 完成检查点数 / 3。当前 **3/3**。关门审计：A-001（self，pass，开放 required = 0）。

## 信息就绪与未知项

本目标不另立信息项；继承 Root [GOAL-001-observability](../GOAL-001-observability/00-meta.md) 的 I-001（指标面合同）、I-003（分母）、I-004（readyz），三者均由本目标 D-001 以证据关闭。I-002（tracing 语义）不在本目标范围（R3 前 closure 即可），因此本波**不**添加 `observability.traces.*` 配置键，避免冻结未定语义。

## 父目标

- [GOAL-001-observability](../GOAL-001-observability/00-meta.md)

## 台账布局

三个可追加台账使用同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter、摘要和条目索引；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 文件。

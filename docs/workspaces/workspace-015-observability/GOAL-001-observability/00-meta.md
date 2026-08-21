---
id: GOAL-001-observability
title: 可观测性（指标导出 + OpenTelemetry）
status: active
parent: null
created: 2026-08-21
updated: 2026-08-21
version: 0.2.0
progress: 1/5
plan_refs:
  - VP-015-observability
primary_plan: VP-015-observability
serves_summary: 交付架构 A4：Prometheus 类指标导出 + OpenTelemetry traces；无收集器仍为 dev/mvp/快测默认。不承载 Sentry / 剖析 / A3 / A5 / Admin 监控页或业务域。
---

# GOAL-001 · 可观测性（指标导出 + OpenTelemetry）

## 概述

本 Root 承载 [VP-015-observability](../../../vision/plans/VP-015-observability.md)（**`active`**）的实现：在已有结构化日志 / request-id 与 `/healthz` `/readyz` 之上，补齐可导出的内核指标面与 OTLP traces。

**边界**：不强制本地默认改成必须有 Prometheus / collector / Jaeger；不承接 Sentry、连续剖析、Admin 监控页或业务域表。安全 finding → VP-009；符合性 gap → VP-010。

## 纲领路线图（P-001）

| 阶段 | 内容 | 先后 | 状态 |
|------|------|------|------|
| R1 | **导出合同与配置面冻结**：指标 scrape 路径/端口/绑定鉴权、基数、内核 vs 模块最小集合、标签不得含秘密（I-001）；Store/对象存储/Job 是否进本波（I-003）；缺省无收集器。承载子目标：[GOAL-002-metrics-export-contract](../GOAL-002-metrics-export-contract/00-meta.md) | 起点 | 已完成（GOAL-002 done 3/3；D-001 + 配置面落地） |
| R2 | **指标 scrape 接入**：Prometheus 类 pull 面；系列携带 `module_id`；未显式配置不成为启动硬依赖；是否扩 `readyz`（I-004 部分）。承载子目标：待立项 | 依赖 R1 | 未开始 |
| R3 | **OpenTelemetry traces 接入**：OTLP 协议/采样/no-op（I-002）；HTTP 请求至少可出 span；未配置 endpoint 不得挡住 mvp/dev。承载子目标：待立项 | 依赖 R1 | 未开始 |
| R4 | **与 request-id 关联**：属性名 / baggage（I-005）；退出 2 的关联判据可核对。承载子目标：待立项 | 依赖 R3 | 未开始 |
| R5 | **双路径证据**：默认无收集器仍能开发快测 + 显式配置下 metrics scrape **与** 至少一条 trace 导出。承载子目标：待立项 | 依赖 R2/R4 | 未开始 |

`progress` = 已完成阶段数 / 5。当前 `1/5`（R1 已完成）。

## 成功标准（方向级）

1. 指标导出面落地；至少一条内核或已启用模块路径可 scrape；系列携带 `module_id`。
2. OTLP traces 可导出；HTTP 请求至少可出 span，并能与现有 request-id / correlation 关联。
3. 未配置收集器时本地/Compose 默认仍能开发与快测。
4. 显式配置后 metrics scrape **与** 至少一条 trace 导出都有可核对证据。
5. 未进入 A3 / A5 / Admin 功能 / 业务域；未改 Charter；未假装交付 Sentry / 剖析 / Grafana 产品。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 指标面：scrape 路径/端口、绑定/鉴权、基数、内核 vs 模块最小集合、标签不得含秘密 | R1 方案冻结 / 实施 | R1 合同冻结 | R1 决策 | verified（2026-08-21） | — | GOAL-002 D-001 §1–§6（`GOAL-002-metrics-export-contract/01-decision/D-001-metrics-export-contract.md`）；对应 VP I-015-001 |
| I-002 | required | Tracing：OTLP HTTP vs gRPC、采样默认、未配置 endpoint 的 no-op | R3 方案冻结 / 实施 | R3 接入前 | R3 决策 | open | — | 对应 VP I-015-002 |
| I-003 | required | Store / 对象存储 / Job 是否进本波分母（HTTP span 已由 VP 退出 2 冻结） | R1 方案冻结 | R1 合同冻结 | R1 决策 | verified（2026-08-21，出局） | — | GOAL-002 D-001 §7：不进本波分母；对应 VP I-015-003（已收窄） |
| I-004 | required | `/metrics` 或 OTLP 是否进入 `readyz`；默认建议未显式配置则不扩 | R2/R3 方案冻结 | R2/R3 接入前 | R2/R3 决策 | verified（2026-08-21，提前闭合：均不进 readyz） | — | GOAL-002 D-001 §8；对应 VP I-015-004 |
| I-005 | required | request-id / correlation 如何写入 span（属性名、是否 baggage） | R4 方案冻结 | R4 关联前 | R4 决策 | open | — | 对应 VP I-015-005 |

## 父目标

- null（Root；Charter `schema-ui-core-admin-foundation@0.2.0` / VP-015）

## 台账布局

新目标为三个可追加台账创建同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter、摘要和条目索引；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 文件。

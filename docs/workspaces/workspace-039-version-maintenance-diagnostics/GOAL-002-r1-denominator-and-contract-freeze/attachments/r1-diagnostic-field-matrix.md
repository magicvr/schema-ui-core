---
doc_type: freeze-matrix
id: r1-diagnostic-field-matrix
parent: GOAL-002-r1-denominator-and-contract-freeze
status: frozen
created: 2026-09-19
updated: 2026-09-19
version: 1.0.0
---

# 冻结矩阵 · 诊断与版本可见性（I-039-001 / I-039-003）

## 产品面

| 面 | 谁可见 | 字段 / 行为 | 不进首波 |
|----|--------|-------------|----------|
| Shell 横幅 | 任何已登录会话 | 仅精确 `runtimeMode` 文案 | commit、模块、DB 大小 |
| Shell 版本入口 | `monitoring.read` | `pkg/version.Version` + QUICKSTART 升级链接 | `Commit`、`BuiltAt`、六包版本 |
| system-monitoring 页 | `monitoring.read`（admin 默认集） | 既有 status 行全字段 | 新页、Grafana、Sentry |
| `/healthz` | 公开 | 既有 version/commit | 不作为 Admin 产品入口 |

## status 行分母（监控页，已交付，保持）

`status`, `availabilityMode`, `ready`, `version`, `commit`, `uptimeSeconds`, `moduleCount`, `modules`, `dbSizeBytes`

排除：Store/对象/Job 指标、剖析、otlp 解析。

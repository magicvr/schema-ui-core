---
doc_type: recon
id: r1-recon-i-039-003-diagnostics
parent: GOAL-002-r1-denominator-and-contract-freeze
status: draft
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# 侦察 · I-039-003 诊断字段分母

> 只读事实，不是冻结决策。派生建议随 C2/C1 用户裁决一并冻结。

## 已交付字段

`GET /api/system-monitoring/status` 单行（需 `monitoring.read`）：

| 字段 | 含义 | 建议首波 |
|------|------|----------|
| `status` | `ok` / `unavailable` / `not-ready`（进程内 ping + ready 门） | **纳入** |
| `availabilityMode` | 原样 `runtime.mode` | **纳入**（区分 read-only） |
| `ready` | bool | **纳入** |
| `version` / `commit` | `pkg/version` | **纳入** |
| `uptimeSeconds` | 进程启动以来 | **纳入** |
| `moduleCount` / `modules` | 当前 plan | **纳入**（已在监控页） |
| `dbSizeBytes` | SQLite 文件大小；PG 下可能为 0 | **纳入**（既有页已有；不新解释） |

`/healthz` `/readyz`：公开探活，字段子集（status/timestamp/version/commit）。不把探活做成 Admin 产品页。

## 明确排除

- Grafana / scrape UI、Sentry、连续剖析（VP-015 非目标 / residual）
- Store / 对象 / Job 指标进分母（I-015-003 residual）
- 新诊断模块或改 Profile 默认集（`I-039-004`）

## 派生建议（待与 P-004 一并确认）

诊断「报告」= **既有 system-monitoring 页**，不新建页面。Shell 只加模式横幅（及可选版本入口，见 I-039-001）。

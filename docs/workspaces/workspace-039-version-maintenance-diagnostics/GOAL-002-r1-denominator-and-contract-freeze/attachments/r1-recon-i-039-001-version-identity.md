---
doc_type: recon
id: r1-recon-i-039-001-version-identity
parent: GOAL-002-r1-denominator-and-contract-freeze
status: draft
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# 侦察 · I-039-001 版本身份与升级入口

> 只读事实，不是冻结决策。

## 权威字段

| 来源 | 字段 | 现状 |
|------|------|------|
| `apps/api/pkg/version` | `Version` / `Commit` / `BuiltAt` | 默认 `"0.1.0"` / `"unknown"` / `"unknown"`；`make build` 经 ldflags 注入 |
| `GET /healthz` | `version` / `commit` | 公开探活，无鉴权；不返回 `BuiltAt` |
| `GET /api/system-monitoring/status` | `version` / `commit` | 需 `monitoring.read`；与 `pkg/version` 同一变量 |
| Prometheus/OTel（VP-015） | version labels | 已交付可观测面，**不进**本 VP 诊断分母 |

仓库根 **无** `CHANGELOG.md`。升级操作文档在 `QUICKSTART.md`（`schema-ui upgrade` · cli+包路径）。六包版本钉在 QUICKSTART 注释，与 `pkg/version` **不是**同一分母。

## 可见性

- `admin.system-monitoring` 仅 **admin** 默认集（mvp/demo 不含）。
- 公开 `/healthz` 已泄露 version/commit（运维探活既成事实）。
- Shell 目前无版本 chip / 升级入口。

## 待 P-004（升级入口）

见编排器提问。候选不在本文件裁决。

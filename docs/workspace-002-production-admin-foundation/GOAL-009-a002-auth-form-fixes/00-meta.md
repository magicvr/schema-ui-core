---
id: GOAL-009-a002-auth-form-fixes
title: A-002 · 缺陷修复（表单提交门禁与认证失效）
status: active
created: 2026-08-03
updated: 2026-08-03
parent: GOAL-001-production-admin-foundation
version: 0.1.0
progress: 0/4
---

# GOAL-009 · A-002 · 缺陷修复（表单提交门禁与认证失效）

## 概述

承接 [Root A-002 审计响应](../GOAL-001-production-admin-foundation/03-audit.md)（independent · 2026-08-03 · fail）：修复两条 required 缺陷——表单 Schema/reaction 错误只展示、不阻断提交（F-002-002），以及 refresh 重试仍 401 / login 后 `/me` 失败时认证状态不丢失（F-002-003）；recommended 三条（F-002-004~006）作为可选加分。Root A-002 的 F-002-001（通用适配层）由 `GOAL-010-a002-schema-adapter` 承接。

## 成功标准

- [ ] **S1 · 表单校验错误阻断提交（F-002-002）**：`apps/web/src/renderer/render.tsx` 在 `gate.errors` / `reaction.errors` 非空时禁用提交按钮，`handleSubmit` 开头拒绝提交；新增「错误显示后请求未发出」回归测试。
- [ ] **S2 · 认证失效状态清理（F-002-003）**：`apps/web/src/account/auth-client.ts` refresh 重试仍 401 时 `clearTokens()` 并触发 `onAuthLost`；login 后 `/me` 失败回滚已存 token 并以登录失败呈现；补二次 401 与 `/me` 失败测试。
- [ ] **S3 · 回归与构建证据**：web `vitest run` 全绿（含新增用例）、`tsc -b` 与生产构建干净；`apps/api` `go test ./...` 全绿（若涉及）。
- [ ] **S4 · 关门审计与 Root finding 关闭**：阶段/关门审计（self + 视需要 `/audit` finding-closure）复核 F-002-002/003 关闭证据，Root 03-audit 对应 finding 按 `fixed` 合法闭合。

> 可选加分 **S5 · recommended 顺带（F-002-004~006，不进进度分母）**：登录页 seed 文案按环境门控（F-002-004）；生产 JWT secret 最小长度/熵校验（F-002-005）；`/healthz` 增加轻量 SQLite/迁移检查或区分 liveness/readiness（F-002-006）。是否纳入由用户决定。

## 派生进度

`progress` 由 S1～S4 四个核心检查点等权派生（`0/4` 起）；S5 可选加分不进分母。检查点不替代审计 finding 或关门结论。

## 信息需求

当前未识别新的 required 信息项：缺陷关闭路径已由 A-002 明确（见 S1/S2 检查点），无需额外信息收集。若实施中发现新的关键未知，按 P-005 登记并回流。

## 依赖与边界

| 项 | 说明 |
|----|------|
| 父目标 | [GOAL-001-production-admin-foundation](../GOAL-001-production-admin-foundation/00-meta.md)（A-002 响应；Root D-014） |
| In | F-002-002 表单提交门禁、F-002-003 认证失效清理、对应回归测试、关门审计；可选 F-002-004~006 |
| Out | F-002-001 通用适配层（归 GOAL-010）；其他业务实体接入；Root / VP-002 关门（Root 层面独立裁决） |

## 父目标

- [GOAL-001-production-admin-foundation](../GOAL-001-production-admin-foundation/00-meta.md)

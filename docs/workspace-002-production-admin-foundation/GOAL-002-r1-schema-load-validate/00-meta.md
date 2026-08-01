---
id: GOAL-002-r1-schema-load-validate
title: R1 · Schema 加载、校验与统一错误面
status: done
created: 2026-08-01
updated: 2026-08-01
parent: GOAL-001-production-admin-foundation
version: 0.4.0
progress: 4/4
---

# GOAL-002 · R1 · Schema 加载、校验与统一错误面

## 概述

在不扩大 `I-PROTO-001 v0.1.3` 的前提下，交付运行时 **Schema 页面文档加载**、**加载时结构校验（D-VAL）** 与 **统一可观察错误面**，使 manifest 中的 `page.schemaUrl` 从「仅建模」变为可消费管线。本目标不切换 App 默认渲染分支（见 GOAL-003），也不交付业务代表性页面全集（见 GOAL-004）。

依据 Root **D-004**（I-001 矩阵 §4 In）与 **D-005**（R1 子目标拆分）。

## 成功标准

- [x] 存在可调用的页面 Schema 加载入口：按已解析的 `schemaUrl`（含路由参数模板展开）获取 page 文档 — `apps/web/src/protocol/load-page.ts` `loadPageDocument`（复用 `resolveSchemaUrl`）；Go `GET /api/schema/{pageId}`（`apps/api/internal/handler/schema.go`）
- [x] 加载路径强制结构校验（至少 page/node 或等价 Ajv 入口）；非法文档 fail-closed，不进入渲染 — `validatePageDocument`（`runtime-schema-validate.ts`，构建期导入 pinned `docs/schemas/*`）；校验失败抛 `PAGE_SCHEMA_INVALID`，加载器不返回文档
- [x] 无效 Schema、网络/解析失败、校验失败暴露统一错误结构（可测、可在 UI 或宿主层观察） — `PageSchemaError`（`code`/`url`/`issues`），覆盖 `PAGE_LOAD_FAILED`/`PAGE_NOT_FOUND`/`PAGE_PARSE_FAILED`/`PAGE_SCHEMA_INVALID`/`PAGE_ID_MISMATCH`；宿主层可捕获观察
- [x] 有自动化测试覆盖：成功加载样例、校验失败、未知/非法 body 的拒绝路径 — `apps/web/src/protocol/load-page.test.ts`（10 项）+ `schema_test.go`

## 派生进度

`progress: 4/4` 由上方 4 条成功标准等权派生；**不**放行 Root R1 检查点（R1 需 002+003+004）。状态 `done`（2026-08-01 经 A-001 independent + A-002 self 关门审计与用户确认）；`progress` 不替代 Root R1 勾选。

## 信息需求

| ID | 问题 / 所需信息 | 级别 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据或结论 |
|----|-----------------|------|----------|----------|-----------------|------|-------------|------------|
| `I-002-001` | Schema 静态资源放置策略（Vite public / 构建内嵌 / mock fetch）最小方案？ | required | 实施路径选择 | 实施开始前 | 对照现有 manifest `schemaUrl` 与 dev server；记决策 | **verified** | 已关闭（D-003） | 用户裁决：Go `GET /api/schema/{pageId}` 端点提供页面文档（manifest `schemaUrl` 不变）；证据 `apps/api/internal/handler/schema.go` + `apps/web/src/protocol/load-page.ts` + 两侧测试 |
| `I-002-002` | 错误码/错误包络是否与现有 ManifestError / RenderError 对齐？ | non-blocking | 错误面一致性 | 验收前 | 对照 `app-manifest` 与 `render` 错误类型 | **concluded** | 已复核（2026-08-01，A-002 响应） | 复核结论：`PageSchemaError` 包络与 `ManifestError`/`RenderError` 约定一致（code+message+位置锚点）；`PAGE_*` 域内专用码不冲突；`url`/`issues` 为 D-003 明示设计 |

## 依赖与边界

| 项 | 说明 |
|----|------|
| 父目标 | `GOAL-001-production-admin-foundation` |
| 同阶段 | 可与 GOAL-003 / GOAL-004 规划并行；**实施上** GOAL-003 默认依赖本目标加载器可用 |
| In | 加载器、D-VAL 串联、统一错误、测试 |
| Out | App 默认分支切换、代表性业务页全集、真实认证、覆盖表扩域 |

## 父目标

- [GOAL-001-production-admin-foundation](../GOAL-001-production-admin-foundation/00-meta.md)

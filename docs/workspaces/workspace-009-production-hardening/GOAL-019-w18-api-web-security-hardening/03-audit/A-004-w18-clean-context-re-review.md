---
id: A-004-w18-clean-context-re-review
doc: audit-entry
goal: GOAL-019-w18-api-web-security-hardening
date: 2026-09-22
source: independent
auditor: Codex REVIEWER / clean-context / read-only
type: cross-audit-independent
scope: W18 A-002 finding closure, API/Web production boundaries, regression evidence, governance consistency and S6 gate
verdict: conditional
open_required: 2
parent: GOAL-001-production-hardening
version: 0.1.0
---

# A-004 · W18 clean-context independent re-review

## 结论

`verdict: conditional`，开放 required = 2。A-002 F-001（API active MFA enrollment 与 nil verifier 的生产装配边界）已真实修复；Web 回归证据和治理勘误仍有缺口，S6 不得放行。

## Required findings

### F-001 · required · A-002 F-002 的浏览器跨 realm 回归不足

复审时现有 Vitest 只用 foreign realm 的 `Object.create` + 自定义 `url` / `Symbol.toStringTag` 伪造 Request，jsdom 没有独立 iframe Request 构造器，不能证明浏览器原生 Request 互操作。需增加真实 Chromium iframe realm：同源 Request 保留 Authorization 与 session-list refresh、跨源 Request 不带两项凭据、跨源 URL 对象不带凭据、auth endpoint 401 不触发 refresh，并观察真实请求。

### F-002 · required · A-002 F-003 的治理勘误不完整

复审发现 E-002 仍写 S1/S2 完成但进度为 `1/6`（应为 `2/6`），`00-meta.md` 仍引用 Web `124/1510`（当前为 `124/1516`），`03-audit.md` 的 I-003 证据未纳入新分母和新增验证，E-004 因而不能宣称 F-003 已完整修复。需同步 E-002、I-003、审计索引与响应事实；S6/I-004/status/progress 在此之前保持不变。

## 已验证

- `server.Run` 的 JWT 配置与 MFA active enrollment fail-closed 路径、资源清理和无 enrollment 启动路径有代码与测试证据。
- Web 实现已移除 realm-sensitive `instanceof` 和 `String(Request)`，结构化读取 `url` / `href` 的逻辑静态正确。
- API `go test ./... -count=1`、Web `124 files / 1516 tests`、typecheck、diff check 均通过。
- goal-tree 主树与状态表包含 GOAL-019，active、5/6、83% 一致；I-004/S6 仍 open。

## 未验证

未运行真实 PostgreSQL 生产数据库；未执行完整 Playwright 套件或生产部署冒烟。本次 opinion 为只读，未修改代码、goal-tree 或目标状态。

---
id: D-001-w18-scope-and-cross-audit
doc: decision-entry
goal: GOAL-019-w18-api-web-security-hardening
date: 2026-09-22
status: accepted
parent: GOAL-001-production-hardening
version: 0.1.0
---

# D-001 · W18 范围、审计模式与独立 Reviewer

## 决定

本波执行一次 VP-009 W18 `api/web` 安全加固审查，范围覆盖当前 `apps/api`、`apps/web` 的安全问题、真实 bug、跨层契约问题与违反既有架构边界的实现；只修复经证据确认、属于共享基架且可在本波闭环的问题。

审计模式采用 `cross`：

1. 主线程先记录扫描事实、修复事实与 self 审计；
2. 实施完成后，使用**上下文干净**的 REVIEWER 子代理进行独立只读交叉审计；
3. Reviewer 按仓库角色配置使用 `gpt-5.6-sol`、reasoning `high`、read-only；不得继承本线程历史，不得直接修改目标状态、progress 或 goal-tree；
4. required finding 未按 `fixed`、用户书面 `accepted-residual` 或 `user-overruled` 合法闭合前，不推进关门。

## 理由

VP-009 是持续安全与健壮性程序；本波跨越 API、Web、生产装配和跨层信任边界，属于 security/production 高影响范围。用户明确要求修复问题并使用独立 Reviewer 交叉审计，因此不能只做一次自审或把扫描摘要当成验收证据。

## 未选方案

- `self`：不足以满足本波 security scope 的独立验证要求。
- 仅扫描不修复：不满足用户要求，也不能闭合真实的 required finding。
- 直接重构整套架构：当前尚无证据支持，除非扫描发现必须升级为 ARCHITECT 决策，否则保持最小修复范围。

## 影响与边界

- 目标承载：`docs/workspaces/workspace-009-production-hardening/GOAL-019-w18-api-web-security-hardening/`
- 不改变 VP-009、Charter、VP-008 `go` 历史状态；若发现会影响 `go` 消费有效性，按 VP-009/VP-008 既有规则记录并暂停受影响放行。
- 不把纯业务功能、新业务域或未被证据触发的 Redis/MQ/多实例设计纳入本波。

---
id: VRev-091-vp009-charter-0-4-0-realign
doc_type: vision-review
title: VP-009 / workspace-009 / Root Charter 0.4.0 re-align
source: self
scope: alignment（VP-009 / workspace-009 / GOAL-001-production-hardening）
verdict: pass
open_required: 0
status: recorded
vision_ref: schema-ui-core-admin-foundation@0.4.0
date: 2026-09-10
auditor: /vision
created: 2026-09-10
updated: 2026-09-10
parent: null
version: 0.1.0
---

# VRev-091 · VP-009 / workspace-009 / Root Charter 0.4.0 re-align

## 背景与触发

用户 2026-09-10 明确要求：检查所有需要 re-align 到 Charter `@0.4.0` 的 VP 与工作区，并集体升级。扫描发现，仍处于现行推进语义的唯一未对齐交付程序是 **VP-009** 及其 lead 工作区和 active Root；VP-010 已在 [VRev-081](VRev-081-vp010-charter-0-4-0-realign.md) 完成对齐，VP-036 planned 时已直接引用 `@0.4.0`。

本次按 P-006 / alignment §8 做**现行投影 re-align**，不回写已关闭 VP、已完成 Goal 或历史审计时点的旧 Charter 语境。

## 扫描与范围判定

| 对象 | 当前状态 | 扫描结果 | 本轮处置 |
|------|----------|----------|----------|
| 唯一 active Charter | `schema-ui-core-admin-foundation@0.4.0` | 当前链源头 | 不修改 |
| VP-009 | `active` | `vision_ref` 仍为 `@0.3.0` | **需 re-align**，更新为 `@0.4.0` |
| workspace-009 | `active` / `vision_role: delivery` | 正文 Charter 声明仍为 `@0.3.0` | **需 re-align**，更新为 `@0.4.0` |
| Root `GOAL-001-production-hardening` | `active` 长期程序容器 | 正文 Charter 声明仍为 `@0.2.0` | **需 re-align**，更新为 `@0.4.0` |
| VP-010 / workspace-010 / Root | `active` 长期程序 | 已精确引用/声明 `@0.4.0` | 无需重复修改 |
| VP-036 | `planned` | `vision_ref` 已为 `@0.4.0` | 无需修改 |
| 已关闭 VP 与历史 Goal | `closed` / `done` 历史证据 | 仍含其验收时点的旧 `vision_ref` 或 Charter 复述 | 按 alignment 保留历史，不纳入现行 re-align |

## 变更核对

| 对象 | 变更 | 保持不变 |
|------|------|----------|
| VP-009 plan | `vision_ref` `@0.3.0 → @0.4.0`；版本 `0.4.0 → 0.4.1`；更新现行 Charter 复述与规划短史 | `status: active`、意图、范围、退出条件、lead、波次与历史证据 |
| workspace-009 | 当前 Charter 声明 `@0.3.0 → @0.4.0`；版本 `0.17.0 → 0.18.0` | `status`、`root_goal`、`canonical_scope`、`vision_role`、`plan_refs`、`primary_plan`、资料引用 |
| Root `GOAL-001-production-hardening` | 当前 Charter 声明 `@0.2.0 → @0.4.0`；版本 `0.14.0 → 0.15.0` | `status: active`、`parent: null`、progress 语义、`plan_refs`/`primary_plan`、波次与 Goal 审计事实 |

## 对齐结论

- VP-009 → workspace-009 → Root 的现行对齐链已恢复到唯一 active Charter `schema-ui-core-admin-foundation@0.4.0`。
- 本次只修正当前声明和机读引用，不改变 VP-009 的长期程序语义，不重开或关闭任何 Goal，不改变任何 progress。
- `workspace-009` scope 的 Charter strategic 宽阻断已解除；后续若要开安全波次，仍须按 `/govern` 扫描当前 Goal 审计与 P-005 信息门禁。
- 本轮未修改 `apps/**`、协议 pin、依赖锁、Profile 默认集、Manifest 或任何 trigger-gated 架构能力。

## VRev-081 finding 响应

[VRev-081](VRev-081-vp010-charter-0-4-0-realign.md) 的 `V-F120` 为 recommended，指向 VP-009 的独立对齐债务。本轮已按其建议完成 VP-009 / workspace-009 / Root 的现行投影同步，故 `V-F120 → fixed`；原报告 verdict 与 finding 原文保留不改写。

## Verdict

**verdict: `pass`（open required = 0）。** 本次 VP-009 对齐属于既有 Charter strategic 修订后的现行投影修复；没有新的 Charter strategic 变更，也没有新增 re-align 范围。

## Findings

无新增 required 或 recommended finding。

## 声明

- 本 Review source = `self`，不冒充 independent。
- 本 Review 不修改 Charter 目的/边界/非目标，不改变 VP、workspace 或 Goal 的 status/progress。
- 已关闭 VP/Goal 的历史 Charter 语境按原验收时点保留；它们不是本轮现行推进范围。
- 该 re-align 完成后，VP-009 可在后续 `/govern` 轮次中按正常安全程序继续扫描与立项。

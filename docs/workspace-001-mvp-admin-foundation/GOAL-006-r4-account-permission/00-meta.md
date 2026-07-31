---
id: GOAL-006-r4-account-permission
title: R4 · 核心账号与权限
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.1.0
---

# GOAL-006 · R4 · 核心账号与权限

## 概述

在 R2 冻结的 MVP 覆盖基线与 R3 Admin 外壳/导航之上，交付核心账号与权限能力：账号权限最小 API 的设计、`D-PERM` 映射冻结，以及前后端可核对的鉴权链路。R4 方案冻结前验证本目标信息项 `I-006-001`；R4 **实施**前闭合父目标 `I-PROTO-002`。

范围依据：Root 将 R4 定义为「核心账号与权限」；协议资料将 `D-PERM` 映射为账号/权限核心能力，并配套 `permissions-inheritance` 等 behavioral fixture 作为结构/行为验证来源。固定协议版本沿用 `schema-ui-docs` artifact `2.7.0`、source commit `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b`。

## 范围边界

### 纳入

- 账号权限最小 API 的设计决策与 `D-PERM` 结构/行为映射。
- 前后端账号/权限链路的实现与可核对验证（R4 实施阶段）。
- 与 R3 Admin 外壳导航上下文的衔接：真实身份/权限来源替换 R3 的默认空 context。

### 排除

- R5 的协议 Renderer 全量、纳入域范例页与 `mock-app` 业务演示；父目标 `I-PROTO-003` 为 R5 验收/关门 required 门禁。
- 完整权限继承产品化、完整主题/设计系统与「完整协议支持」主张；R3 的 `I-005-*` 验证边界不自动扩展到权限语义。
- 任何未在 `I-PROTO-001` v0.1.3 覆盖表内纳入的协议能力。

## 高层路线图（P-001）

1. **契约发现与信息就绪**：进行中；验证 `I-006-001`（账号权限最小 API 与 `D-PERM` 映射），并在 R4 方案冻结前闭合父目标 `I-PROTO-002`。
2. **方案冻结**：未开始；冻结最小 API、权限模型映射与前后端集成边界。
3. **R4 实施**：未开始；受 `I-PROTO-002` 实施门禁约束。
4. **验证与关门**：未开始；补结构/行为/运行时证据与实施阶段自审后再谈 `done`。

## 成功标准

- [ ] R4 必需信息项 `I-006-001` 已由证据验证，或有用户书面接受的有界 residual；未知项没有被默认为已知。
- [ ] 账号权限最小 API 与 `D-PERM` 映射方案已冻结；父目标 `I-PROTO-002` 在 R4 实施前合法闭合。
- [ ] 前后端账号/权限链路具备可核对的实现与验证路径（结构/行为/运行时证据），R3 导航 context 得到真实身份/权限来源衔接。
- [ ] R5 Renderer/范例与完整协议支持保持边界外；关门前无开放 required finding。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-006-001 | required | 账号权限最小 API 与 `D-PERM` 映射是否完整？ | 方案冻结 / 实施 | 方案冻结前 | 对照 `D-PERM` 与 `permissions-inheritance` 等固定 fixture，设计最小 API 并映射 React/Go 实现路径 | open | 不延期；R4 方案冻结前须验证 | 待确认 |

父目标 `I-PROTO-002`（required，R4 **实施**门禁）、`I-PROTO-003`（required，R5 验收/关门门禁）由 Root 维护，本目标不修改其状态。`I-PROTO-004`（vendor vs pin，non-blocking）仍 open；关闭该项时须补 schema-conformance 等价性校验或显式记录等价范围（GOAL-005 A-007 F-002 跟进项）。

## 父目标

- [GOAL-001-mvp-admin-foundation](../GOAL-001-mvp-admin-foundation/00-meta.md)

## 备注

- 立项日期：2026-07-31，承接 R3 关门（GOAL-005 `done`）。
- Root 纲领进度仍为 `3/6`；本目标立项不抬升 progress，也不放行 R4 实施。
- 实施前闭合 `I-PROTO-002` 与 `I-006-001`；收集信息阶段不实施代码。

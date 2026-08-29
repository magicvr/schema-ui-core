---
id: GOAL-003-r2-public-release-channel
title: R2 · 公开发布通道闭环（npmjs @magicvr 六包 + apps/api/v0.4.0 + golden-field 无凭据消费实证）
status: active
parent: GOAL-001-distribution-formalization
created: 2026-08-29
updated: 2026-08-29
version: 0.1.0
progress: 0/4
---

# GOAL-003 · R2 · 公开发布通道闭环

## 概述

承接 Root R2 与 VP-024 判据 #2：把 R1 的 serve 面与六包**真实公开发布**——npmjs.com `@magicvr/schema-ui-*` 六包（用户裁决：@magicvr 先行 · truthful 公开；@schema-ui 为 org 就绪后正式化候选，见 D-001 §6）、Go `apps/api/v0.4.0` origin tag + 公共 proxy；golden-field 升级到 v0.4.0 + 六包并完成**无凭据消费实证**（npmjs 公开免 token · go proxy 免认证）；发布流程成文（脚本 + 凭据注入点 + scope 迁移 changelog 注记）。核销 go 后清单第 ② 项（npmjs 公开可见性）与 R1 残余 2（registry 级骨架消费）的发布前提。

## 成功标准（可验证检查点）

- [x] C1：npmjs.com 真实发布六包 `@magicvr/schema-ui-{protocol,lib,theme,ui,renderer,shell}`（lib/theme/ui/shell 0.1.0 · protocol/renderer 0.2.0）——授权/scope 变更记录见 D-001 §6（@schema-ui 为 org 就绪后正式化候选）
- [x] C2：Go `apps/api/v0.4.0` origin tag + 公共 proxy `go get` 实证（含 serve 面与 CLI `schema-ui serve`）
- [x] C3：golden-field 升级：`go.mod` → v0.4.0 · web 六包 → npmjs 公开消费；全程 registry 语义（无 replace / 无 file:）；**无凭据**安装/拉取 + 三探针全绿（protocol 2.9 / render 1573B / token 覆盖）
- [x] C4：发布流程成文：npmjs 发布脚本（`.env` `npm_token` 注入点 + 临时 .npmrc，token 不入库/不落盘）+ scope 迁移 changelog 注记（D-001 §6）

## 方案与路线（P-001）

| 阶段 | 内容 | 状态 |
|------|------|------|
| S1 | 发布面定案：scope/授权/脚本形态（D-001 · I-024-001 关闭） | **已关门**（2026-08-29 · 用户裁决：@magicvr 先行 · 真实发布） |
| S2 | npmjs 发布脚本 + dry-run → 真实发布六包；Go v0.4.0 tag + push | **已关门**（2026-08-29 · E-002） |
| S3 | golden-field 升级 + 无凭据消费实证（go get / pnpm install / 探针） | **已关门**（2026-08-29 · E-002） |
| S4 | 证据 + 自审关门（含 R1 残余 2 核销前提确认） | 独立审计 A-002（grok）收尾中 |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-024-001 | required | npmjs 公开发布授权（scope 名 + 凭据方式） | 判据 #2 | R2 | 用户书面授权或裁决 | **verified**（2026-08-29 用户裁决：@magicvr 先行真实发布 · token = 仓库根 `.env` `npm_token`；@schema-ui = org 就绪后正式化候选） | — | 用户裁决（goal round）· D-001 §6 |
| I-024-002 | required | CI 槽位环境（真实 runner / 用户环境等价 + 凭据） | 判据 #3（R3） | R3 | workflow 实跑验证 | **verified 于 R3**（不阻 R2） | R3 前置 | — |
| I-024-003 | required | fork 对照的同一演进集样本（V 演进选择） | 判据 #4（R4） | R4 | 样本设计与 fork 基线 | open | R4 前置 | 待确认 |

## 父目标

- `GOAL-001-distribution-formalization`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger。

## 备注

- 审计模式：S2/S4 涉及 release 实证门禁 = independent（grok build 先例；D-001 已定）；如 grok 不可用则 self + 用户确认。
- 关联：VP-024 判据 #2 · Root 大纲 R2 · R1 残余 2（registry 级消费）· go 后清单第 ② 项。
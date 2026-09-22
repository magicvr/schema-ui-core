---
id: GOAL-019-w18-api-web-security-hardening
title: W18 · api/web 多迭代安全与架构加固
status: done
parent: GOAL-001-production-hardening
created: 2026-09-22
updated: 2026-09-22
version: 0.7.0
progress: 100%
---

# GOAL-019 · W18 · api/web 多迭代安全与架构加固

## 概述

承接 VP-009 的 W18 波次，对当前 `apps/api` 与 `apps/web` 进行一次面向真实风险的全量安全、缺陷与架构边界扫描；对确认属于本波范围的问题实施最小修复，补齐可复现的回归证据，并经 self + 独立 Reviewer 交叉审计后再决定是否关门。

本目标不重开历史 VP/波次，不预设一定存在缺陷；扫描结果、分类、修复范围与残余风险均以本目标台账中的可核对事实为准。

## 纲领路线图

1. [x] S1：建立 api/web 基线并完成候选问题扫描与分类
2. [x] S2：冻结本波修复范围、信息门禁与验证矩阵
3. [x] S3：实施并验证 API 侧 required 修复
4. [x] S4：实施并验证 Web 侧 required 修复
5. [x] S5：完成跨层回归与 self 阶段审计
6. [x] S6：由干净上下文 Reviewer 独立交叉审计；响应意见并完成关门判断（A-007 pass；用户关门裁决待定）

`progress` 仅按以上 6 个显式检查点等权派生；不作为 finding 闭合、阶段放行或目标完成依据。

## 成功标准

- [x] api/web 扫描范围、基线、候选问题及证据路径已落盘；事实与推断已分开
- [x] 本波 required 修复范围、非目标、信息门禁和验证矩阵已冻结
- [x] 属于本波范围的 required 安全/正确性/架构问题已修复并有针对性回归证据；不适用或移交项有理由和触发条件
- [x] API、Web、跨层及相关生产装配路径完成与风险相称的验证；未把既有失败伪装成新成果
- [x] self 审计与独立 Reviewer 交叉审计均已落盘；required finding 按 `fixed` / `accepted-residual` / `user-overruled` 合法闭合
- [x] 目标与 `workspace-009` 的 goal-tree、执行台账和审计台账一致；满足条件后再由用户确认 `done`

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 当前 api 扫描是否存在影响本波的安全、正确性或架构边界问题？ | S2 范围冻结 / S3 / S6 | S2 | SCOUT 扫描、代码与测试对照、必要时最小复现 | verified | — | `apps/api/server/serve.go`：Run 绕过 Config.validate，非 development 缺 secret 时使用公开固定 JWT fallback；见 E-002 / D-002 |
| I-002 | required | 当前 web 扫描是否存在影响本波的安全、正确性或架构边界问题？ | S2 范围冻结 / S4 / S6 | S2 | SCOUT 扫描、代码与测试对照、必要时最小复现 | verified | — | `apps/web/src/account/auth-client.ts`：Request URL 被 String 化导致跨源认证头注入风险；见 E-002 / D-002 |
| I-003 | required | 本波变更的 API/Web/跨层验证分母与生产装配路径是什么？ | S2 / S5 / S6 | S2 | 记录 baseline、变更路径、测试矩阵与运行环境 | verified | — | E-004：API `go test ./...`；Web `npm run test` 124/1516、`npm run typecheck`；针对性 Vitest 与真实 Chromium iframe E2E 已纳入 |
| I-004 | required | 修复后是否仍有未闭合 required finding 或需用户裁决的 residual/overruled？ | S6 关门 | S6 | self + clean-context Reviewer 复核并落盘 | verified | — | A-007 `pass`，open_required=0；无 residual/overruled |

## 当前扫描分类

- API MAJOR：`server.Run` 可绕过 `Config.validate()`，非 development 缺失 JWT secret 时落入公开固定密钥；已由 E-003 修复、Run 级测试锁定，并经 A-007 独立确认 fixed。
- Web MAJOR：`isSameOrigin` / `isSessionListRequest` 对 `Request` 使用 `String(input)`，可能向跨源请求注入当前 Bearer / refresh token；已改为跨 realm 安全解析，并由 Vitest + Chromium E2E、A-007 独立确认 fixed。
- API required（A-002 F-001）：下游 `serve` 组合显式不装配 MFA；已改为启动前检测 active enrollment，命中或检测失败均 fail closed，A-007 独立确认 fixed。
- Web NOTE：预览 blob iframe 当前已有空 sandbox 与 `opener = null`；未发现可控字符串进入 HTML 的证据，本波不把它升级为 finding。

## 审计模式

- 模式：`cross`（VP-009 security scope）
- self：本目标 `03-audit/`，只记录已发生事实
- independent：上下文干净的 REVIEWER 子代理，按角色配置使用 `gpt-5.6-sol`、`high`、read-only；仅写 `source: independent` 意见，不直接修改 status/progress/goal-tree

## 父目标

- `GOAL-001-production-hardening`（VP-009 长期安全与健壮性程序容器）

## 台账布局

本目标使用 `01-decision/`、`02-execution/`、`03-audit/` 三个平铺 ledger 目录；每条决策、执行事实和审计意见单独落盘并由索引引用。

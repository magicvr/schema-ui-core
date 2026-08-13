---
id: GOAL-006-w5-recordview-declared-fields
title: W5 · recordView 声明字段符合性（declared-fields 契约 + dev/文档卫生）
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-08-14
updated: 2026-08-14
version: 0.2.0
progress: 4/4
---

# GOAL-006 · W5 · recordView 声明字段符合性

## 概述

本子目标是 VP-010 / workspace-010 的**第五波**（2026-08-13/14 执行，2026-08-14 补建档案）：**recordView 声明字段（declared-fields）符合性**——渲染器按 schema 声明的字段元数据渲染记录详情（字段标题 / 顺序 / 包含集），并对声明缺失 / 异常采用 **fail-open 兜底**（回退默认渲染，不崩不黑屏）。以 users / roles / activity 三模块 schema 为代表声明，i18n 双语标题同步；伴随开发脚本与 QUICKSTART 文档卫生（同批未归档提交）。

性质：**as-designed（schema 声明）vs as-built（renderer）对齐** 的符合性波次，属 VP-010 程序范围。

## 当前边界

- 范围：`apps/web/src/renderer/render.ts(x)`（recordView 渲染）、三模块 schema（`users.json` / `roles.json` / `activity.json` 声明字段）、i18n 消息（en-US / zh-CN）、相关测试（`render.test.ts` / `render.test.tsx` / `ui-bilingual.test.tsx` / `schema-keys.structural.test.ts`）；dev 脚本（dev.cmd 等待 API ready、按 PID 停止）与 QUICKSTART 排版。
- **不**改变 Profile 默认集、模块矩阵、Manifest 装配语义、协议 pin、认证 / 授权 / 数据隔离 / fail-closed 等已冻结共同门禁语义；**不**私增协议 capability（declared-fields 为本地渲染器契约，缺失时 fail-open 回退，见 I-001）。

## 成功标准与路线图（P-001）

- [x] **S1 · 声明字段实现**：renderer recordView 支持 schema 声明字段（标题 / 顺序）；users / roles / activity schema 增加声明；i18n 双语；测试补齐。（2026-08-13 · commit `7f10fff`）
- [x] **S2 · fail-open 契约修正**：声明字段缺失 / 异常时 fail-open 兜底 + 健壮性缺口修正。（2026-08-14 · commit `a831754`）
- [x] **S3 · dev/文档卫生**：dev 脚本等待 API ready 后启动 Web、stop 按 PID 精确停止；QUICKSTART 修正 dev.cmd 调用前缀与排版。（2026-08-13 · commit `5c309ff` / `c420e5d`）
- [x] **S4 · 回归与关门**：HEAD 全量回归（V-001～V-006 绿；V-007/V-008 因 F-1 基础设施受阻）+ 关门自审（A-001 conditional）。（2026-08-14 · E-004）

`progress: 4/4` 由四个等权检查点派生。本目标于 2026-08-14 关门；F-1（容器 smoke 复现性，跨门禁）移交 freshness review 决策，见 E-004 / A-001。

## 审计策略

本波为共享渲染层整改（可逆、已提交）；回归以 HEAD 冻结矩阵（V-001～V-005）为证据，关门采用 **self** 审计（A-001）。若用户要求，可补 independent 复审。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | recordView 声明字段是否为协议（schema-ui-docs）能力？缺失声明时的契约语义是什么？ | S2 方案 | S2 前 | 核对本地 registry / 上游 fixtures / 提交语义 | verified | — | declared-fields 为**本地渲染器契约**（非协议 capability）；缺失 / 异常时 **fail-open 回退默认渲染**（a831754 提交语义 + 测试断言） |

## 父目标

- [GOAL-001-design-implementation-conformance](../GOAL-001-design-implementation-conformance/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；索引与目录条目共同构成正式记录。

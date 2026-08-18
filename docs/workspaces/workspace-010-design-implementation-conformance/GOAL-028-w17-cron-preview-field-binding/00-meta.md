---
id: GOAL-028-w17-cron-preview-field-binding
title: W17 · Cron 字段绑定与中文 describeCron
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
progress: 3/4
---

# GOAL-028 · W17 · Cron 字段绑定与中文 describeCron

VP-010 / workspace-010 的**第十七波**：承接 [GOAL-024](../GOAL-024-w16-user-perspective-improvements/00-meta.md) A-005 F-004 / A-007 F-003（recommended open）——把 Cron 预览接到创建/编辑任务表单的 `cron` 字段下方，并把 `describeCron` 做成按 `Accept-Language` 协商的中文/英文自然语言，而不是英文 stub。

不重开已关门的 GOAL-024；本波为 Root 下级新波次。

## 当前边界

- **范围**：定时任务创建/编辑模态的 `cron` 字段即时预览；`POST /api/scheduled-tasks/cron/preview` 的 `description` 中文人话（及 en 对等句）；移除页面级独立预览块。
- **非范围**：改 Cron 解析器语义、6 段 cron、Lightbox/导入等 W16 其它残余；不改 Profile 默认集 / 模块矩阵 / Manifest 装配 / 协议 pin。

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：字段绑定方式与 `describeCron` 句式（D-001）。
- [x] **S2 · 实施**：schema + `cron-preview` 绑定 + `describeCron` 本地化（E-002）。
- [x] **S3 · 测试与回归**：Go/Web 定向 + `tsc`（E-003）；未跑 handler 全包 / 全量 vitest / e2e。
- [ ] **S4 · 自审与关门**：self 审计；goal-tree / workspace 同步。

progress: 由四个等权检查点派生；当前 **3/4**。

## 审计策略

| 阶段 | 模式 | 说明 |
|------|------|------|
| S1～S3 | none → 后继兜底 | 可逆 UX / 文案；无门禁语义变化 |
| S4 关门 | self | 常规、边界清楚 |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|
| I-001 | required | 预览如何接到表单字段而不新增协议控件类型 | S1 / S2 | S1 | 对照 FormControls 白名单与 custom-component 注册表 | **verified** | D-001：本地 `afterComponent` 扩展，不新增 protocol field type |

## 父目标

- [GOAL-001-design-implementation-conformance](../GOAL-001-design-implementation-conformance/00-meta.md)

## 溯源

- GOAL-024 A-005 F-004 / A-007 F-003（recommended open）
- GOAL-024 D-004：不写 residual，另开后续波次
- GOAL-027 D-001 §2 原冻结：「表单 Cron 字段下方挂 `cron-preview`」+ 中文描述

## 台账布局

- `01-decision/`：D-NNN；`02-execution/`：E-NNN；`03-audit/`：A-NNN。

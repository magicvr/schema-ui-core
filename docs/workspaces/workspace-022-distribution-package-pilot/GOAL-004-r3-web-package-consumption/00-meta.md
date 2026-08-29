---
id: GOAL-004-r3-web-package-consumption
title: R3 · Web 包闭环
status: active
parent: GOAL-001-distribution-package-pilot
created: 2026-08-29
updated: 2026-08-29
version: 0.1.0
progress: 0/4
---

# GOAL-004 · R3 · Web 包闭环

## 概述

承接 Root R3 检查点与 VP-022 退出判据 #2：验证「空下游 app 仅通过 npm 包组（protocol / renderer / shell / ui）组装，可渲染与主线同一的 schema 页面集；品牌定制仍走 Token 覆盖路径」。S1 扫描已盘点 `apps/web` 六面（protocol/renderer/app–shell/host/components/theme）+ 发现协议 pin 漂移（见 E-001）。

## 成功标准（阶段检查点）

- [ ] **S1 · 拆包边界设计**：I-002 落盘——包清单（protocol/lib/theme/ui/renderer/shell 六包）+ 导出面规则 + peer 依赖耦合矩阵（React/Tailwind/ajv）+ 实施路径评估（E-001 扫描事实 + 设计附件）
- [ ] **S2 · 打包链路与首包**：npm 包可发布形态（tsc/Vite lib + `pnpm pack` 链路）产出 ≥1 包（protocol 优先）；golden web app（`attachments/golden-web/`，仅包依赖）消费验证
- [ ] **S3 · 渲染闭环**：空下游 app 经包组渲染与主线同一 schema 页面集；Token 覆盖（brand.css 机制）定制验证；peer 矩阵实测
- [ ] **S4 · 关门**：self 审计 + 判据 #2 满足声明（供 R5 go/no-go）

## 信息就绪

| ID | 级别 | 所需信息 | 影响门禁 | 最晚阶段 | 状态 | 证据 |
|----|------|----------|----------|----------|------|------|
| I-002 | required | npm 包拆分与 peer 依赖策略（六包边界 / 导出面 / React-Tailwind-ajv 耦合矩阵） | R3 实施 / R5 发布 | R3 | **collecting**（S1 设计草拟中） | `web-package-boundary-design-v0.1.md` |
| I-007 | required | 协议 pin 漂移实况：web `APP_MANIFEST_PROTOCOL_VERSION=2.9`（`81aa1d8`）vs Charter/roadmap `v2.8.0`（`521cff8`）——差异范围与影响 | R4 演练基线 / R5 发布 | R4 | **collecting**（E-001 初证；待 `/vision` 裁决是否 pin bump） | E-001 |

## 父目标

- `GOAL-001-distribution-package-pilot`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；条目见各索引。
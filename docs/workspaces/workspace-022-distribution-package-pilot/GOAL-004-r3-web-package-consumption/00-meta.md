---
id: GOAL-004-r3-web-package-consumption
title: R3 · Web 包闭环
status: done
parent: GOAL-001-distribution-package-pilot
created: 2026-08-29
updated: 2026-08-29
version: 0.3.0
progress: 4/4
---

# GOAL-004 · R3 · Web 包闭环

## 概述

承接 Root R3 检查点与 VP-022 退出判据 #2：验证「空下游 app 仅通过 npm 包组（protocol / renderer / shell / ui）组装，可渲染与主线同一的 schema 页面集；品牌定制仍走 Token 覆盖路径」。S1 扫描已盘点 `apps/web` 六面（protocol/renderer/app–shell/host/components/theme）+ 发现协议 pin 漂移（见 E-001）。

## 成功标准（阶段检查点）

- [x] **S1 · 拆包边界设计**：I-002 落盘——六包清单/导出面规则/peer 耦合矩阵/三路径评估；用户裁决 **B · Vite lib 产物打包**（E-001 + 设计附件 + D-001）
- [x] **S2 · 打包链路与首包**：Vite lib + declaration 链路 → `@schema-ui/protocol` v0.1.0（306 kB 自包含 + d.ts）；golden-web 仅包依赖 `pnpm install` + probe **PASS**（E-002）
- [x] **S3 · 渲染闭环**：用户裁决粗粒度单包（D-002）→ `@schema-ui/renderer` v0.1.0（436.7 kB · React peer）→ golden-web SSR 渲染真实形态 schema 文档 **PASS** + Token 覆盖纪律断言 **PASS** + 能力门控 fail-closed 可观测（E-003）
- [x] **S4 · 关门**：A-001 self `pass`（0 required；F-006 recommended = d.ts 链路 TS5056 留 go 后）；**判据 #2 满足声明成立**；GOAL-004 `done 4/4`

## 信息就绪

| ID | 级别 | 所需信息 | 影响门禁 | 最晚阶段 | 状态 | 证据 |
|----|------|----------|----------|----------|------|------|
| I-002 | required | npm 包拆分与 peer 依赖策略（六包边界 / 导出面 / React-Tailwind-ajv 耦合矩阵） | R3 实施 / R5 发布 | R3 | **collecting**（S1 设计草拟中） | `web-package-boundary-design-v0.1.md` |
| I-007 | required | 协议 pin 漂移实况：web `APP_MANIFEST_PROTOCOL_VERSION=2.9`（`81aa1d8`）vs Charter/roadmap `v2.8.0`（`521cff8`）——差异范围与影响 | R4 演练基线 / R5 发布 | R4 | **collecting**（E-001 初证；待 `/vision` 裁决是否 pin bump） | E-001 |

## 父目标

- `GOAL-001-distribution-package-pilot`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；条目见各索引。
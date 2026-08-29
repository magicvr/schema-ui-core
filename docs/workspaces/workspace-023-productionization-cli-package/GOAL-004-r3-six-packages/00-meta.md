---
id: GOAL-004-r3-six-packages
title: R3 · 六包细化与 d.ts 自动化
status: done
parent: GOAL-001-productionization-cli-package
created: 2026-08-29
updated: 2026-08-29
version: 0.2.0
progress: 4/4
---

# GOAL-004 · R3 · 六包细化与 d.ts 自动化

## 概述

承接 Root R3 与 VP-023 判据 #3：`@magicvr/schema-ui-{protocol,lib,theme,ui,renderer,shell}` 六包独立可发布（renderer 保持自包含 bundle v0.1.0；**新四包 = lib/theme/ui/shell**）+ d.ts 自动化管线（TS5056 修复验证）→ 冻结面 v1.3.0 升格。

## 成功标准（阶段检查点）

- [x] **S1 · 聚合入口与构建脚本**：lib（含 i18n 面）/theme/ui/components 核心/shell/app 四聚合入口 + `scripts/build-lib-packages.mjs`（Vite JS API 循环）
- [x] **S2 · 产物与 d.ts**：四包产物 + 单文件 d.ts（dts-bundle-generator 管线修复尝试）；golden-field 独立安装 + import 探针
- [x] **S3 · 发布与消费**：四包 registry 发布；golden-field 消费验证（含 renderer 0.1.0 共存）
- [x] **S4 · 关门**：冻结面 v1.3.0 + A-001 self `pass` → 判据 #3 满足声明 · F-006 核销

## 信息就绪

| ID | 级别 | 所需信息 | 影响门禁 | 最晚阶段 | 状态 | 证据 |
|----|------|----------|----------|----------|------|------|
| I-023-005 | required | d.ts 管线方案（dts-bundle-generator vs 目录化改名）对 TS5056 的修复有效性 | S2 产物 | S2 | **verified**（改名方案根治 TS5056 · 五包 declaration 全 0） | — | GOAL-004 E-001 |

## 父目标

- `GOAL-001-productionization-cli-package`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；条目见各索引。
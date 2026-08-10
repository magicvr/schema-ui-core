---
id: GOAL-003-s2-ui-schema-bilingual
title: S2 · 固定 UI 与 Schema 分母双语化（titleKey/labelKey 真解析）
status: done
parent: GOAL-001-localization-and-system-settings
created: 2026-08-09
updated: 2026-08-09
version: 0.2.0
progress: 5/5
---

# GOAL-003 · S2 · 固定 UI 与 Schema 分母双语化

## 概述

承接 Root [GOAL-001-localization-and-system-settings](../GOAL-001-localization-and-system-settings/00-meta.md) 的 **S2 阶段**：按 S0 冻结的 F-V029 分母（固定 UI U1～U7 + 双 Profile Runtime Manifest 12 pageId/schema 并集 + M1～M4）完成前端可见文本双语化——固定 UI 文案迁入翻译 catalog；Manifest 与 Schema 驱动文本通过 `titleKey`/`labelKey`/`contentKey` 真解析（Renderer/导航渲染时解析，缺失可观察 + 安全回退）。

**方案依据**：Root D-002 §I-L10N-001（前端 key 解析，2026-08-09 用户书面裁决）+ F-V029 冻结表。本目标只实施与验证，不重新决策。

**范围纪律**：不改写 `schema-ui-docs@v2.7.0` 协议语义；只使用本地 component-registry 已声明的 key 字段（`labelKey`/`titleKey`/`contentKey`，及为缺口 additive 登记的本地扩展字段并记录）；不涉及系统设置四类字段与 API 行为（S3）；不动错误 envelope（S4）。

## 成功标准（可验收 · 等权检查点 · 共 5 项）

- [x] **C1**：固定 UI 面 U1～U7 全部迁入 catalog（登录页、Shell 顶栏/侧栏/用户导航、通用反馈/loading/empty/error/permission denied/success、通用 table/search/form/modal/validation 文案、Settings 面外壳）——`zh-CN`/`en-US` 双 catalog 有键，无硬编码英文残留路径（en-US 为规范基线可保留 en 文本）。
- [x] **C2**：Manifest `titleKey`/`labelKey` 真解析：导航 label 与页面标题在渲染时经 catalog 解析；`label`/`title` 文本 fallback；缺失 key 可观察（事件 + 测试）；Manifest 数据侧为并集内 pageId/nav 项补 `*Key`（additive，协议已声明字段）。
- [x] **C3**：Schema 驱动文本真解析：registry 已声明 `*Key` 的组件（field label、section title、text content、table 列 label、action/submit label、状态文案等）在 Renderer 渲染时解析；协议字面文本保留为 fallback；缺失 key 安全回退不阻断。
- [x] **C4**：M4 缺失翻译流程：制造缺失 key → 记录 locale/key/UI 路径 → 安全文本 fallback → 主流程仍可完成（vitest 断言事件与回退）。
- [x] **C5**：验证：vitest 驱动 shipped 函数全绿（新增 S2 测试覆盖 C1–C4）；`npm run build` 通过；证据矩阵开始填充（F-V029 表格 U 行/页面行证据路径回填）。

## 派生进度展示

`progress: 5/5` 由上方 5 个等权检查点派生；仅为展示，不放行阶段、不关闭 finding。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | S2 实施输入（catalog 结构、manifest/schema key 语义、F-V029 分母）是否齐备 | C1–C5 | 实施前 | 读 Root D-002 + F-V029 表 + component-registry | **closed** | — | Root D-002 + F-V029 已冻结（2026-08-09） |
| I-002 | required | 本地 component-registry 对缺口的 additive key 字段登记（如 placeholder/submitLabel） | C3 | 实施中 | 盘点 registry 声明后按需登记并记录 | **closed** | — | S2 盘点：registry 已声明 labelKey/titleKey/contentKey；action 按钮 label 直用协议 label（加 labelKey additive 登记见 D-001） |

## 父目标

- [GOAL-001-localization-and-system-settings](../GOAL-001-localization-and-system-settings/00-meta.md)（Root；本目标为 S2 阶段子目标）

## 台账布局

本目标使用 ledger 目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter 与条目表；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*`。

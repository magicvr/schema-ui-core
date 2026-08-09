---
id: GOAL-002-s1-locale-core
title: S1 · 多语种核心（locale 解析/资源/切换/格式化）
status: active
parent: GOAL-001-localization-and-system-settings
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
progress: 0/5
---

# GOAL-002 · S1 · 多语种核心

## 概述

承接 Root [GOAL-001-localization-and-system-settings](../GOAL-001-localization-and-system-settings/00-meta.md) 的 **S1 阶段**：实现前端多语种运行时核心——locale 解析/匹配/回退、翻译资源装载、缺失 key 可观察回退、用户切换（localStorage 单通道）、HTML `lang` 与 locale-aware 日期/数字格式化，并用单元测试驱动 shipped 函数。

**方案依据（已 accepted 的 Root 决策）**：D-002 §I-L10N-001（前端 key 解析）、§I-L10N-002（localStorage 单通道 + 优先级）。本目标不重新决策，只实施与验证。

**范围纪律**：不实现任何页面/设置面双语化（S2/S3）；不新增账号资料字段；不引入服务端 locale 协商（S4）；不改写 `schema-ui-docs@v2.7.0` 协议语义。

## 成功标准（可验收 · 等权检查点 · 共 5 项）

- [ ] **C1**：locale 解析纯单元（`zh-CN`/`en-US`/`auto` 解析、匹配、回退优先级：用户显式 → 系统默认（非 auto）→ 浏览器偏好 → `en-US` 安全回退）可单测；登录前后同一优先级一致应用。
- [ ] **C2**：翻译资源装载：`zh-CN`/`en-US` catalog 为纯数据文件；运行时按 locale 装载；装载失败安全降级（不阻断启动）。
- [ ] **C3**：缺失 key 可观察且安全回退：`schema-ui:missing-translation` 事件（detail 含 locale/key/path）+ 回退链（当前语种 → en-US → key 本身）；不渲染为空、不抛异常。
- [ ] **C4**：用户切换：Shell/用户菜单语种切换器（普通用户可达、无需设置权限），localStorage 单通道持久化，登出不清除；切换立即生效并触发 `lang` 更新。
- [ ] **C5**：HTML `lang` 与格式化：`document.documentElement.lang` 跟随有效 locale；日期/数字格式随有效 locale（`Intl.DateTimeFormat`/`NumberFormat`，无自定义格式模板）。
- [ ] **C6**：验证：vitest 驱动真实 shipped 函数（resolver/catalog/switch/format），输出捕获到 scratch 日志；`npm run build` 通过。

## 派生进度展示

`progress: X/5` 由上方 5 个等权检查点派生；仅为展示，不放行阶段、不关闭 finding。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | S1 实施输入（locale 语义、优先级、持久化边界）是否齐备 | C1–C6 | 实施前 | 读 Root D-002 §I-L10N-001/002 | **closed** | — | Root D-002 已 accepted（2026-08-09 用户书面裁决） |

## 父目标

- [GOAL-001-localization-and-system-settings](../GOAL-001-localization-and-system-settings/00-meta.md)（Root；本目标为 S1 阶段子目标）

## 台账布局

本目标使用 ledger 目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter 与条目表；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*`。

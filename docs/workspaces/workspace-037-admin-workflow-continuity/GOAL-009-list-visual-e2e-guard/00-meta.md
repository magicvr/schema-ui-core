---
id: GOAL-009-list-visual-e2e-guard
title: 列表视觉浏览器级回归守卫
status: done
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
version: 1.0.0
progress: 4/4
plan_refs:
  - VP-037-admin-workflow-continuity
primary_plan: VP-037-admin-workflow-continuity
vision_ref: schema-ui-core-admin-foundation@0.4.0
---

# GOAL-009 · 列表视觉浏览器级回归守卫

## 概述

承接 R6 `GOAL-007` 审计 `A-002` 的 **F-003**（recommended）：R6 的列表视觉合同（C5/C7/C8）此前**没有任何持久化的浏览器级回归**——`apps/web/e2e/` 对 `data-filter-*`、`data-saved-view*`、`data-list-page-actions`、`data-table-footer` 的匹配数为 0，全部结论依赖 jsdom 与一次性临时预览页。

这一缺口有直接后果：R6 的 C5 → C7 → C8 连续三轮回归，**每一轮都由用户先于测试发现**。原因是结构性的——R6 的合同本质上是真实布局（计算高度、按钮与输入的贴合、响应式断点实际隐藏了哪些控件），而 jsdom 没有布局引擎、`matchMedia` 被 stub，jsdom 套件只能断言 class 字符串，看不到合同真正描述的属性。

用户于 2026-09-18 指示「先做列表视觉的 e2e 守卫」，本目标即以真实浏览器断言补齐该层守卫，并已于同日以 `done · 4/4` 完成：守卫在 mvp/admin 两 profile 下通过，且经 6 种变异测试证明非空转。

## 范围与边界

- 在 `apps/web/e2e/` 新增持久化 Playwright 规格，覆盖 R6 列表视觉合同中**jsdom 无法观察**的部分：搜索控件配对几何、页面 actions 高度一致性、折叠开关的存在性与 token 来源、列表布局顺序与列表内 footer、视图表单的归属与贴近度。
- 断言风格为**关系型**（高度相等、贴合为负间距、"确实有控件被隐藏"），不是像素快照；token 改值或间距微调不产生噪音，只有破坏合同才失败。
- 必须在 CI 的**两个 profile（mvp 与 admin）**下都通过（CI 浏览器矩阵为 `profile × dialect`）。
- 复用既有 e2e 挂具（`sign-in.ts`、profile 校验、临时 DB），不新增 CI job、不改 `playwright.config.ts` 的契约。

明确非目标：不引入视觉快照/像素 diff 测试；不修改列表实现（除非守卫发现真实缺陷）；不覆盖没有这些表面的页面；不改动 `playwright.config.ts` 的 profile/dialect 契约；不重开 R6 的视觉范围或 GOAL-008 的类型检查范围。

## 高层路线图

1. **C1 · 可达性与基线**：确认 e2e 挂具本机可跑、确认可作为守卫目标的页面与 profile 可达性、采集真实几何基线。已完成，证据见 `D-001`、`E-001`。
2. **C2 · 守卫规格与实现**：实现 `e2e/list-visual-surface.spec.ts`，按合同条目组织断言。已完成，证据见 `E-002`。
3. **C3 · 双 profile 回归与变异验证**：在 mvp/admin 两 profile 下通过，并以变异测试证明守卫非空转。已完成，证据见 `E-003`。
4. **C4 · self 审计与投影**：完成 self 审计、闭合 F-003 并向 Root 投影。已完成，证据见 `A-001`、`E-003`。

## 成功检查点

- [x] C1：e2e 挂具本机可跑（既有规格 1 passed）；确定 roles 页为目标（同时具备 search form 与 toolbar）；确认 mvp/admin 两 profile 均可达 `/users`、`/roles`；采集桌面 1440 与窄屏 700 的真实几何基线。
- [x] C2：`e2e/list-visual-surface.spec.ts` 已实现，覆盖 C7 配对/图标/视图表单归属、C8 高度/token/开关存在性、C5 布局顺序与列表内 footer。
- [x] C3：规格在 `admin` 与 `mvp` 两 profile 下均通过；6 种真实回归变异全部被捕获。
- [x] C4：self 审计完成（`A-001` `pass`，开放 required = 0），F-003 合法闭合，Root 投影已记录。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-009-001 | required | e2e 挂具在本机能否跑通（Go/Chromium/端口）？ | C1/C2 | C1 | 跑既有规格冒烟；检查 `go version`、Chromium 缓存、25080/25173 端口 | verified | 2026-09-18 已完成 | `E-001` |
| I-009-002 | required | 哪个页面同时具备 search form、toolbar 与视图 surface，且在 CI 两 profile 下可达？ | C1/C2 | C1 | 全仓检索 schema 的 `mode:search`/`toolbar`/`filters`；两 profile 实跑探测导航 | verified | 2026-09-18 已完成；roles 页满足，mvp/admin 均可达 | `E-001` |
| I-009-003 | required | 合同各条目的真实几何基线是什么？ | C2 | C2 | 在真实浏览器采集 1440/700 两档的盒模型与计算样式 | verified | 2026-09-18 已完成 | `E-001` |
| I-009-004 | non-blocking | 是否引入视觉快照（pixel diff）测试？ | 范围外 | C2 | 用户未要求；快照对 token 改值极敏感，与“关系型断言”取向冲突 | deferred | 触发：用户明确要求像素级回归时 | `D-001` 未选方案 |

## 父目标

- `[workspace-037-admin-workflow-continuity]` `GOAL-001-admin-workflow-continuity`。

## 台账布局

本目标从第一条记录起使用平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/`，并保留 `attachments/`。

## 备注

- 本目标**不是** Root 的纲领阶段，不改变 Root 六阶段分母与 `progress: 5/6`；它是 Root 下的整改子目标，与 GOAL-008 平行。
- 目标关闭的是 GOAL-007 `A-002 F-003`（recommended）。该 finding 不阻断 R6 关门（R6 已 `done · 8/8`），本目标是主动补齐其回归覆盖。
- R5 `GOAL-006` 的 `R5-I-004` 用户书面关门确认仍开放，本目标不替代、不关闭它。

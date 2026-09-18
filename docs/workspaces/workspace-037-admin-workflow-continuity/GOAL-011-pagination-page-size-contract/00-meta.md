---
id: GOAL-011-pagination-page-size-contract
title: 分页每页条数契约与跳转按钮文案修正
status: done
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
version: 1.1.0
progress: 4/4
plan_refs:
  - VP-037-admin-workflow-continuity
primary_plan: VP-037-admin-workflow-continuity
vision_ref: schema-ui-core-admin-foundation@0.4.0
---

# GOAL-011 · 分页每页条数契约与跳转按钮文案修正

## 概述

用户于 2026-09-18 在使用中报告两个缺陷：

1. **每页条数下拉不生效且默认值名实不符**：下拉默认显示 **10**，但实际生效的是 **20**；选择 **10** 也不能真正生效。
2. **页码跳转的确认按钮文案错误**：显示为「搜索」（英文 `Search`），应为「跳转」（英文同义）。

根因（`D-001`，已由代码与实测确认）：

- 前端 `apps/web/src/renderer/resource.ts` 的 `DEFAULT_PAGE_SIZE = 10`，而服务端 `apps/api/internal/handler/resources.go` 的 `DefaultPageSize = 20`。`buildResourceQuery()` 在 `pageSize === DEFAULT_PAGE_SIZE` 时**省略该参数**，于是：默认显示 10 → 参数被省略 → 服务端按 20 返回（默认名实不符）；在选中 20 之后再选 10 → 仍被省略 → 依旧 20（10 永不生效）。
- 跳转表单的提交按钮复用了 `t("feedback.search")`，语义错位。
- 既有回归 `schema-table.test.tsx` 的 mock 把「未带 pageSize」当作 10（`?? "10"`），把错误契约固化进了测试，因此该缺陷从未被测试拦下。

本目标已于 2026-09-18 以 `done · 4/4` 完成：默认值统一为 20 且 10 真正生效、按钮文案改为「跳转 / Go」、清理同类字面量、新增前后端常量绑定的结构守卫与真实浏览器用例；self 审计 `A-001` `pass`、开放 required = 0。用户随后授权走 Root 目标关闭流程（见 Root `D-017`）。

## 范围与边界

- 修正前后端分页默认值契约：前端默认改为 **20**（与产品要求和服务端一致），并确保 **10 能真正生效**。
- 修正跳转按钮文案：新增 `feedback.jumpToPage`（zh「跳转」/ en「Go」）并用于该按钮；表单的 `aria-label` 仍为「跳至页 / Go to page」。
- 修正回归测试里固化错误契约的 mock 默认值，并补「10 生效」与「默认 20」的断言。
- 新增**结构性守卫**：断言前端 `DEFAULT_PAGE_SIZE` 与服务端 `DefaultPageSize` 一致，防止该耦合再次静默漂移（同类缺陷曾由用户在浏览器中发现，而非测试）。

明确非目标：不改分页交互形态、分页尺寸选项集合（10/20/50/100）、跳转校验逻辑或其他控件文案；不改 `apps/api` 任何代码（服务端 20 是既有正确行为）；不改通知中心等其它自带 `pageSize` 的表面；不重开 R2/R6（分页展示属已交付范围，本目标为缺陷整改子目标）。

## 高层路线图

1. **C1 · 复现与根因**：以可执行证据复现「默认显示 10 / 实际 20」与「10 不生效」，确认客户端与服务端默认值不一致。证据见 `D-001`、`E-001`。
2. **C2 · 修正实现**：统一默认值为 20、修正跳转按钮文案、清理其余把 10 当默认的字面量。证据见 `E-002`。
3. **C3 · 回归与防复发**：补「10 生效 / 默认 20 / 文案正确」断言与前后端契约结构守卫，复跑全量。证据见 `E-003`。
4. **C4 · 审计与投影**：self 审计并投影 Root（不关闭 Root/VP；其关门由用户授权的关门流程单独处理）。证据见 `A-001`、`E-004`。

## 成功检查点

- [x] C1：复现证据成立（改动前 2 条测试失败），根因定位到 `DEFAULT_PAGE_SIZE` 与服务端默认值不一致 + 省略参数的耦合（`E-001`）。
- [x] C2：前端默认值为 20；选择 10 时请求带 `pageSize=10`；跳转按钮文案为「跳转」/「Go」（`E-002`）。
- [x] C3：新增断言与结构守卫通过；修正后的 mock 不再固化错误契约；Vitest 113/1434、typecheck exit 0、e2e 3 passed × admin/mvp（`E-003`）。
- [x] C4：self 审计 `pass`、开放 required = 0，Root 投影已记录（`A-001`、`E-004`）。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-011-001 | required | 服务端分页默认值到底是多少？前端与它是否一致？ | C1/C2 | C1 | 读 `apps/api/internal/handler/resources.go` 并对照前端 `DEFAULT_PAGE_SIZE` | verified | 2026-09-18 已完成 | 服务端 `DefaultPageSize = 20`；前端为 10（`E-001`） |
| I-011-002 | required | 为什么「10」不生效？参数在哪一层被丢弃？ | C1/C2 | C1 | 追踪 `buildResourceQuery` 与调用点，构造复现用例 | verified | 2026-09-18 已完成 | `pageSize === DEFAULT_PAGE_SIZE` 时省略参数（`E-001`） |
| I-011-003 | non-blocking | 其它表面（通知中心等）是否同样受影响？ | 范围外 | C3 | 检查各 `pageSize` 常量与请求构造 | verified | 2026-09-18 已核对：通知中心显式带参、独立常量，不受影响 | `E-001` |
| I-011-004 | required | 跳转按钮文案的正确键与两种语言值？ | C2 | C2 | 查 i18n catalog 现有键 | verified | 2026-09-18 已完成 | 新增 `feedback.jumpToPage`（跳转 / Go）（`E-002`） |

## 父目标

- `[workspace-037-admin-workflow-continuity]` `GOAL-001-admin-workflow-continuity`。

## 台账布局

本目标从第一条记录起使用平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/`，并保留 `attachments/`。

## 备注

- 本目标**不是** Root 的纲领阶段，不改变 Root 六阶段分母；它是 Root 下的整改子目标，与 `GOAL-008`/`GOAL-009`/`GOAL-010` 平行。
- 用户已授权：本目标两个缺陷修正完成后，走 **Root 目标关闭流程**。该授权同时构成 `GOAL-006 R5-I-004` 的书面确认来源，按 P-004 在关闭流程中留痕；本目标自身不关闭 Root/VP。
- 分页与跳转表面属 VP-037 方向级退出判据 #7（单页有效列表仍显示分页区域）的已交付范围；本条为缺陷整改，不重开 R6 的视觉范围。

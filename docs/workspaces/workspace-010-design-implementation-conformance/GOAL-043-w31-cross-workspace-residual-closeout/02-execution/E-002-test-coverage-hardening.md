---
id: E-002-test-coverage-hardening
doc: execution-entry
status: recorded
goal_id: GOAL-043-w31-cross-workspace-residual-closeout
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-design-implementation-conformance
version: 1.0.0
---

# E-002 · 可处理项修复：暗色断言 / 第二页面 / Host-resource 对照（C2）

## 1. `GOAL-009 A-001 F-001` 暗色下开关的计算背景

`apps/web/e2e/list-visual-surface.spec.ts` 的 C8 用例在原有「`.dark` 下 `--control` 有覆盖值且≠浅色」之后，新增暗色下的**计算背景**断言，共四条不变量：

- 根上的 `--control` = 暗色覆盖值；
- 开关自身的 `--control` 继承同一值；
- `--color-control` 主题别名解析到该值（防止构建期被内联成字面色）；
- 开关的 `backgroundColor` 等于该值（真正钉住「消费 token」）。

**实现要点（记录在案的坑）**：开关带 `transition-colors`。第一版直接加 `.dark` 后读 `getComputedStyle` 拿到的是**过渡插值帧**（实测 `oklab(0.955 0 0)`，即浅色），会被误判为缺陷。改为读取前内联 `transition: none`、读完还原，即得到解析后的目标值。这不是产品缺陷，是测试读法问题——已在用例注释中写明。

**变异验证**（证明新增断言非空转）：把开关的类名从 `bg-control` 改成 `bg-control dark:bg-[oklch(0.955_0_0)]`（暗色下硬编码浅色），浅色断言**全部仍通过**，新增暗色断言失败并报出「the toggle must consume --control in dark mode, not a hard-coded light colour」——正是 F-001 所担心的失效模式。

## 2. `GOAL-009 A-001 F-002` 第二页面覆盖

分页契约用例改为**参数化双页面**：`roles` 与 `users` 各跑一次，断言内容不变（默认显示 20、首个列表请求不带 `pageSize`、选 10 后确实发出 `pageSize=10`、跳转按钮文案）。这样「某页自身的接线（数据源/表格 id）回归」不再只靠单页覆盖。导航抽为 `openListPage(page, linkName)`；`openRolesList` 保留为薄封装。

## 3. `GOAL-005 A-002 F-002` / `R5-I-005` Host 终态与普通 resource 反馈的直接对照

- `apps/web/src/app/HostFailureScreen.tsx`：把文案映射表与通用键导出（`HOST_FAILURE_MESSAGE_KEYS`、`HOST_FAILURE_GENERIC_KEY`），**无行为变更**——仅让对照测试能读取该表而不是靠解析源码或复制一份。
- 新增 `apps/web/src/host/resource-feedback-parity.test.ts`（3 用例）：对 8 个两表面共有的条件（maintenance / timeout / offline / rate-limited / authentication / forbidden / not-found / unavailable）逐条对照——
  1. resource 策略必须把这些信号分到同一概念类别，host 表必须对每个条件有显式条目（不得缺失）；
  2. 两侧都不得回退到通用文案，且各自命名空间不同（`hostFailure.*` vs `feedback.*`）、不得相互混用；
  3. 两侧引用的每个 key 必须在 `zh-CN` 与 `en-US` 两个 catalog 中都存在且非空。

**变异验证**：① 从 host 表删掉 `maintenance` → 3 项用例失败并指名 `[ 'maintenance' ]`；② 把 `zh-CN` 的 `hostFailure.maintenance` 改名 → 本地化用例失败并指名 `zh-CN: hostFailure.maintenance`。两者均已还原。

C2 完成（`tsc` 余项见 `E-003`）。回归数据见 `E-004`。

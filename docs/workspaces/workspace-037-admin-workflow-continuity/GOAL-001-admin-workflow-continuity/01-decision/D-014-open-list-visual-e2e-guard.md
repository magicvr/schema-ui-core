---
id: D-014-open-list-visual-e2e-guard
doc: decision-entry
status: accepted
goal_id: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
parent: null
version: 1.0.0
---

# D-014 · 用户指示补做列表视觉 e2e 守卫并开设 GOAL-009

## 决定

用户于 2026-09-18 指示「先做列表视觉的 e2e 守卫」，即同意在推进 R5 关门之前，先补齐 R6 `A-002 F-003`（recommended）所指的持久化浏览器级回归覆盖。

据此开设整改子目标 `GOAL-009-list-visual-e2e-guard`，承接该 finding。它是 **Root 下的整改子目标，不是纲领阶段**，因此不改变 Root 六阶段分母与 `progress: 5/6`。

同日完成 C1～C4 并关门为 `done · 4/4`（`A-001` self 审计 `pass`，开放 required = 0）。

## 范围

- 新增 `apps/web/e2e/list-visual-surface.spec.ts`，在真实浏览器中断言 R6 列表视觉合同中 jsdom 无法观察的部分。
- 必须在 `mvp` 与 `admin` 两 profile 下通过；复用既有 e2e 挂具与 CI job，不新增 job。
- 断言为关系型（相等/贴合/确实隐藏），不引入像素快照；不修改被守卫的实现。

## 理由

R6 的 C5 → C7 → C8 三轮回归**全部由用户在真实浏览器中发现**，测试未捕获任何一轮。根因是 jsdom 的能力边界：无布局引擎（`getBoundingClientRect` 恒为 0）、`matchMedia` 被 stub，因此 jsdom 只能断言 class 字符串，而回归恰恰是"class 改了、几何也错了"。把断言放到真实浏览器，才能观察合同真正约束的属性。

用户选择先做此项而非直接推进 R5 关门，是合理的次序：R5 是组合验收与关门准备，其判据之一包含非目标边界与回归完整性；在关门之前补齐已知的覆盖缺口，可避免把"无浏览器级视觉守卫"这一已知弱点带入关门结论。

## 未选方案

- **引入视觉快照（pixel diff）测试**：对 token 改值与字体渲染差异极敏感，需基线图与平台一致性投入，且与合同的关系型语义不匹配；用户未要求。
- **扩展现有 `w4-long-content-spotcheck.spec.ts`**：该规格 scope 为 W4 长内容截断，混入 R6 视觉合同会使失败诊断指向错误合同。
- **只在 admin 下运行**：CI 矩阵含 mvp；仅 admin 通过在 mvp 腿会失败或被 skip，等于制造新的"看不见的覆盖"。
- **与 GOAL-008 合并为一个目标**：F-003（视觉回归覆盖）与 F-005（类型检查证据口径）性质、验证手段与影响面均不同，合并会使台账混淆两类证据。

## 门禁

本目标已完成。Root/VP 不因本轮关门：R5 `R5-I-004` 用户书面关门确认仍开放。

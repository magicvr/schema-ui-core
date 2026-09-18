---
id: D-001-guard-scope-and-assertion-style
doc: decision-entry
status: accepted
goal_id: GOAL-009-list-visual-e2e-guard
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# D-001 · 守卫范围、目标页与断言风格

## 决定

1. **承接来源**：本目标承接 GOAL-007 `A-002 F-003`（recommended）——R6 列表视觉合同缺持久化浏览器级回归。用户 2026-09-18 指示「先做列表视觉的 e2e 守卫」。
2. **目标页**：`/roles`（`apps/api/modules/roles/schema/roles.json`）。它是唯一**同时**具备 search form（`mode: search`，2 个字段：`q` + `system`）、table `toolbar`（Export / New role）的页面，因而能同时承载 C7 搜索配对、C8 高度一致性、C5 页面 actions 顺序与视图 surface 的全部断言。仅有的另一个带 `filters` 的页面（account 的 sessions 表）只有 1 个筛选项且无 search form，无法触达折叠开关。
3. **断言风格：关系型，非像素快照**。断言"同一行控件高度相等"、"按钮与输入之间无间隙（`<= 0`）"、"窄档确实有控件被隐藏"、"开关背景等于 `--control` 解析值"，而不是断言 `32px` 或 `oklch(0.955 0 0)` 这类字面值。理由：token 改值、间距微调、内容变化都不应产生噪音；只有破坏合同才应失败。token 断言写成"等于解析出的 token 值"还额外钉住了"开关确实由该 token 驱动"这一语义（而不只是"它是某个颜色"）。
4. **profile 策略**：规格必须在 CI 浏览器矩阵的两个 profile（`mvp`、`admin`）下都通过。经实跑探测确认两 profile 均可达 `/roles` 与 `/users`（`mvp` 下侧栏分组需展开后才渲染链接）。
5. **导航与档位**：侧栏是 `hidden lg:block`，因此规格**先在桌面宽度导航**，再按需 `setViewportSize` 收窄——这也是"响应式行为"更真实的触发方式（用户在会话中收窄窗口）。

## 理由

R6 的三轮回归全部由用户在真实浏览器中发现，而不是由测试发现。根因是 jsdom 的能力边界：没有布局引擎（`getBoundingClientRect` 恒为 0），且 `matchMedia` 被 stub（`useTier` 的响应式分档无法真实求值）。jsdom 套件因此只能断言 class 字符串——而 C5 的回归恰恰是"class 改了、几何也错了"。把断言放在真实浏览器里，才能观察到合同真正约束的属性。

选择关系型断言而非像素快照，是因为 R6 合同的语言本身就是关系型的（范例页的语义是"与列配置同高"而非"恰好 32px"、"与输入贴合为同一控件"而非特定像素）。快照会与 token 调优冲突，产生维护噪音，最终被削弱或跳过——那会退回"没有守卫"的状态。

## 未选方案

- **视觉快照 / pixel diff 测试**：对 token 改值与字体渲染差异极敏感，需要基线图与平台一致性投入，且与本合同的关系型语义不匹配；用户未要求。记为 `I-009-004` deferred。
- **用 `/users` 而非 `/roles` 作为目标页**：users 页也同时具备 search form 与 toolbar，但其列与数据更依赖种子规模；roles 页的 2 个搜索字段在窄档恰好产生"隐藏 1 项"的可判定状态，更适合钉住 C8 item 3。
- **扩展现有 e2e 规格（如 `w4-long-content-spotcheck.spec.ts`）**：该规格的 scope 是 W4 长内容截断，混入 R6 视觉合同会让失败诊断指向错误的合同。
- **只在 admin 下运行**：CI 矩阵包含 mvp；若规格仅 admin 通过，会在 CI 的 mvp 腿静默失败或被 `skip`，等于制造新的"看不见的覆盖"。
- **新增独立 CI job**：现有 `browser-e2e` job 已按 `profile × dialect` 跑完 `e2e/`，新规格自动纳入，无需改 CI 契约。

## 门禁

C2 实现后必须在 mvp/admin 两 profile 实跑通过，并在 C3 以变异测试证明守卫非空转（否则与 GOAL-008 F-005 同类：一个从不失败的守卫等同于没有守卫）。C4 自审须核对上述两项证据。

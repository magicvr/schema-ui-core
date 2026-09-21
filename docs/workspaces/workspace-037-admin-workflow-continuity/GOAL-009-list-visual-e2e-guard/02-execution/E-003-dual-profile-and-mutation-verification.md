---
id: E-003-dual-profile-and-mutation-verification
doc: execution-entry
status: recorded
goal_id: GOAL-009-list-visual-e2e-guard
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# E-003 · 双 profile 回归与 6 种变异验证

2026-09-18，完成 C3。守卫"通过"本身不足以证明其有效（GOAL-008 F-005 的教训：一个从不失败的守卫等于没有守卫），故对 R6 已知的每一类回归做变异测试，确认守卫会失败。

## 变异验证（6/6 捕获）

| # | 变异（模拟的回归） | 结果 |
|---|-------------------|------|
| 1 | 把搜索提交按钮移出配对槽（**C5 的原始回归**） | ✅ 失败：`expect(locator).toBeVisible() failed` |
| 2 | toolbar 触发器 `h-8 → h-9`（C8 item 1 高度阶梯） | ✅ 失败：`page-action controls must share one height (got 32, 36, 36, 36, 36)` |
| 3 | 折叠开关改为恒渲染（C8 item 3 空展开） | ✅ 两个测试均失败（桌面档断言开关数应为 0） |
| 4 | 开关样式退回 `bg-background`/`text-muted-foreground`（C8 item 2） | ✅ 失败：`the toggle must use the --control token` |
| 5 | 删除 `.dark` 下的 `--control` 覆盖（两层 token 合同） | ✅ 失败：`--control must have a dark override` |
| 6 | 把视图表单移回页面 actions 行（**C7 item 3 的原始抱怨**） | ✅ 失败：`view form must render inside the view surface` |

每次变异后均还原源文件并校验内容一致；最终工作树仅含本目标新增的规格文件。

## 双 profile 与类型检查

| 场景 | 结果 |
|------|------|
| `APP_PROFILE=admin npx playwright test e2e/list-visual-surface.spec.ts` | **2 passed** |
| `APP_PROFILE=mvp npx playwright test e2e/list-visual-surface.spec.ts` | **2 passed** |
| `npm run typecheck`（`tsc -b && tsc -p e2e/tsconfig.json`） | exit 0 |
| `git status --short` | 仅 `?? e2e/list-visual-surface.spec.ts` |

## 意义

变异 1 与变异 6 分别复现了 **C5 引入的搜索配对回归**与 **C7 修复的视图表单位置问题**——这两个正是 R6 三轮回归中由用户而非测试发现的问题。守卫现已能在真实浏览器中捕获同类回归，F-003 所指的覆盖缺口由此闭合。

C4 自审见 `A-001`。

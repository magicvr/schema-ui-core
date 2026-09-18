---
id: E-002-guard-spec-implementation
doc: execution-entry
status: recorded
goal_id: GOAL-009-list-visual-e2e-guard
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# E-002 · 实现列表视觉浏览器级守卫规格

2026-09-18，完成 C2：新增 `apps/web/e2e/list-visual-surface.spec.ts`（2 个测试，按合同条目组织）。

## 覆盖映射（合同 → 断言）

| 合同 | 断言 |
|------|------|
| C7 / A-003 · W13 T-03 搜索配对 | 提交按钮与关键词输入位于**同一** `data-filter-item`；按钮不在 `data-filter-actions`；`button.left - input.right <= 0`（贴合/重叠，无间隙）；`|verticalOffset| <= 1`；按钮含 `rounded-l-none`、输入含 `rounded-r-none` |
| C8 item 1 高度一致 | `data-list-page-actions` 的全部子项与内部 `button` 高度集合大小为 **1**（失败信息打印实测值）；且 `data-saved-view-columns-trigger` 高度与之相等 |
| C7 item 1 图标 | columns 与 save-view 触发器各含 1 个 `svg` |
| C7 item 3 视图表单归属 | 表单渲染在 `data-saved-views-surface` **内部**、不在 `data-list-page-actions`；与触发按钮垂直间隔 ≤ 一个控件高（40px 上限，规则而非魔数） |
| C5 / D-002 布局顺序 | 筛选面板 → 页面 actions → 列表 surface（文档顺序 + 垂直不重叠）；`data-table-footer` 位于列表 surface 内部；有效列表响应存在 `data-pagination-footer` |
| C8 item 3 空展开抑制 | 桌面档 2 项筛选项时开关**不存在**且无控件被隐藏；收窄后开关**存在**且确实隐藏了控件 |
| C8 item 2 折叠开关语义 token | `--control` / `--control-foreground` 已声明；开关背景**等于解析出的 `--control`**；且**不等于**重置按钮背景；重置按钮与开关不是同一元素；`.dark` 下 `--control` 有非空且不同于浅色的覆盖值 |

## 设计要点

- 断言为关系型；token 断言写成"等于解析出的 token 值"，同时钉住"背景确实由该 token 驱动"这一语义。
- 导航先在桌面宽度完成（侧栏 `hidden lg:block`），再 `setViewportSize` 收窄，以获得真实的响应式分档行为。
- 复用 `e2e/sign-in.ts` 的登录助手与既有挂具，不新增 CI job、不改 `playwright.config.ts`。

## 验证事实

| 场景 | 结果 |
|------|------|
| `APP_PROFILE=admin` 全规格 | **2 passed**（33.2s） |
| `APP_PROFILE=mvp` 全规格 | **2 passed**（35.4s） |
| `npm run typecheck`（含 `tsc -p e2e/tsconfig.json`） | exit 0（新规格类型检查通过） |

C3 的变异验证见 `E-003`。

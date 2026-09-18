---
id: A-001-goal011-self-closeout
doc: audit-entry
status: recorded
goal_id: GOAL-011-pagination-page-size-contract
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# A-001 · GOAL-011 C1～C3 self 审计

- **source**：self
- **日期**：2026-09-18
- **scope**：`GOAL-011-pagination-page-size-contract` 的 C1～C3（复现、根因、修正实现、回归与防复发）
- **verdict**：`pass`

## 成果（可核对）

| 核对项 | 结论 | 证据 |
|--------|------|------|
| 用户报告的两个现象是否被**先复现** | 是。改动前两条测试即失败：`expected '10' to be '20'`、`expected 'Search' to be 'Go'` | `E-001` §1 |
| 根因是否落到具体代码行 | 是。前端 `DEFAULT_PAGE_SIZE = 10` 与服务端 `handler.DefaultPageSize = 20` 不一致，而 `buildResourceQuery` 在该值上省略参数 | `E-001` §2 |
| 「默认改成 20」是否达成 | 是。前端默认 20，与服务端一致；下拉显示 20 | `E-002`、`E-003` |
| 「确保 10 生效」是否有可核对证据 | 是。单元/组件层断言请求含 `pageSize=10` 且渲染 10 行；浏览器层在真实 Go 服务上等待到 `pageSize=10` 的请求 | `E-003` |
| 文案是否修正 | 是。按钮 = `feedback.jumpToPage`（跳转 / Go）；表单标签与输入框语义未变（跳至页 / Go to page） | `E-002`、`E-003` |
| 是否有防复发机制 | 是。新增结构守卫把前端常量与服务端常量绑定，并锁定「非默认尺寸必须上线路」与两个 catalog 的文案键 | `E-003` |
| 是否修掉了**固化缺陷的测试** | 是。`schema-table.test.tsx` 的 mock 由 `?? "10"` 改为 `?? "20"`；`resource.test.ts` 与 `saved-views.ui.test.tsx` 的旧字面量断言改用共享常量 | `E-003` |
| 是否越界 | 否。未改 `apps/api`；未改分页尺寸集合、跳转校验、其它控件文案；未改通知中心等显式带参表面（`I-011-003`） | `E-001` §4、`E-002` |
| 回归是否全绿 | 是。Vitest 113/1434、typecheck exit 0、e2e 3 passed × admin/mvp、`git diff --check` 通过 | `E-003` |

## 偏差与残余

| 项 | 说明 |
|----|------|
| 默认值仍依赖「前后端一致」 | 参数省略是既有优化，因此两常量必须相等；本目标以结构守卫固定该不变式。若将来服务端改默认值，守卫会先失败（这是有意的：宁可要求同步改前端，也不要静默漂移）。 |
| 保存视图会记录默认 pageSize | `currentSavedViewState` 在 `pageSize` 有值时写入；现记录 20。既有保存视图若存过 10，恢复后仍会显示 10（并**真正生效**，因为 10 ≠ 默认）。未改，属既有语义、非本次缺陷。 |
| 通知中心等表面 | 自有 `DEFAULT_PAGE_SIZE = 10` 且显式带参，行为正确；本目标未统一（`I-011-003` 记为范围外，未升格为 finding）。 |
| 分页尺寸选项集合 | 仍为 10/20/50/100。服务端 `error.invalidPageSize` 文案为「1–100 的整数」，与选项集合一致，无需改。 |
| 浏览器层未断言行数变化 | e2e 的 roles 种子行数少于 10，无法用行数证明 10 生效；改用「请求确实带 `pageSize=10`」作为证据（这正是缺陷的判定点）。 |

## 结论

C1～C3 达成，verdict `pass`，开放 required finding = 0。两个用户报告的缺陷均已修正并有跨层证据，防复发守卫已就位。建议进入 C4：投影 Root 并关闭本目标。本条与 `E-004` 均**不**关闭 Root 或 VP-037——用户授权的 Root 关门流程由 `/govern` 按 `GOAL-006` 的 `R5-I-004` 门禁单独执行。

---
id: D-001-pagination-contract-fix
doc: decision
status: accepted
goal_id: GOAL-011-pagination-page-size-contract
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# D-001 · 分页默认值契约与跳转文案修正方案

## 根因（可核对）

| 事实 | 证据 |
|------|------|
| 服务端分页默认值为 **20** | `apps/api/internal/handler/resources.go:40` `DefaultPageSize = 20`（注释：共享列表默认值，W15-F12） |
| 前端客户端默认值为 **10** | `apps/web/src/renderer/resource.ts:3` `export const DEFAULT_PAGE_SIZE = 10;` |
| 等于该常量时参数被**省略** | `resource.ts:256` `if (query.pageSize !== undefined && query.pageSize !== DEFAULT_PAGE_SIZE)` |
| 于是「显示 10 / 实际 20」 | 初始 `localQuery = { page: 1, pageSize: 10 }`（`schema-table.tsx:615`）→ 下拉显示 `query.pageSize ?? 10` = 10 → 请求省略 pageSize → 服务端按 20 返回 |
| 于是「选 10 不生效」 | 选 20 后 `query.pageSize = 20`（发 `pageSize=20`）；再选 10 → 回到省略分支 → 服务端仍 20 |
| 测试为何没拦下 | 既有回归 `schema-table.test.tsx` 的 mock 把缺参当作 10（`?? "10"`），把错误契约写进了测试 |
| 跳转按钮文案错位 | `schema-table.tsx:1589` 提交按钮使用 `t("feedback.search")`（「搜索」/`Search`） |

## 决定

1. **统一默认值为 20**：`apps/web/src/renderer/resource.ts` 的 `DEFAULT_PAGE_SIZE` 改为 `20`，并加注释说明它必须与服务端 `handler.DefaultPageSize` 一致。
2. **清理「把 10 当默认」的其余字面量**：`schema-table.tsx`（初始 query 与下拉回退值）、`render.tsx`（搜索表单写过滤器时的兜底 query）、`activity-export.tsx`（导出兜底 query）改用 `DEFAULT_PAGE_SIZE`，避免默认值再次分裂。
3. **确保 10 生效**：保留「等于默认值即省略参数」的既有优化——因为默认值现在与服务端一致，选择 10 时会显式发送 `pageSize=10`。
4. **跳转按钮文案**：新增 i18n 键 `feedback.jumpToPage`（zh「跳转」/ en「Go」）并用于该提交按钮；表单 `aria-label` 保持「跳至页 / Go to page」，输入框语义不变。
5. **防复发（结构守卫）**：新增结构测试，直接读 `apps/api/internal/handler/resources.go` 解析 `DefaultPageSize`，断言前端 `DEFAULT_PAGE_SIZE` 与之相等；并断言两个 i18n catalog 都存在 `feedback.jumpToPage` 且与 `feedback.search` 不同值。这样该耦合一旦漂移，测试立即失败，而不是等用户再次在浏览器里发现。
6. **回归修正**：把 `schema-table.test.tsx` 的 mock 默认值改为 20（模拟真实服务端），并补「默认显示 20 且首请求不带 pageSize」「选择 10 时请求带 `pageSize=10` 且行数/页数随之变化」「跳转按钮文案为 Go / 跳转」断言。

## 未选方案

- **只把前端常量改成 20，不动其它字面量**：`render.tsx`/`activity-export.tsx` 的裸 `10` 在常量变更后会从「被省略」变成「显式发 10」，搜索表单提交会静默把每页条数改成 10；不采纳。
- **取消「等于默认值即省略」的优化，永远显式发送 `pageSize`**：请求更冗长且会改变大量既有断言的 URL 形态；在默认值已被结构守卫锁死的前提下收益有限，不采纳。
- **把服务端默认值改成 10**：用户明确要求「默认改成 20」，且服务端 20 是既有正确行为；不采纳。
- **只改文案不做契约守卫**：该缺陷的成因是跨层常量漂移，只改值不留守卫会重犯；不采纳。
- **顺带统一通知中心等其它表面的默认值**：超出用户报告范围，且那些表面显式带参、行为正确；不采纳（`I-011-003` 已核对）。

## 边界与不变量

- 不改 `apps/api`（服务端 20 不变），不改分页尺寸选项集合与跳转校验逻辑。
- 不改变「无结果时仍显示分页区域」的 R6 合同，也不重开 R6/R2 的已交付结论。
- 本目标为整改子目标，不改变 Root 六阶段分母与 `progress`；不关闭 Root/VP。

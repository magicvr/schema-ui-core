---
id: E-003-regression-and-guards
doc: execution-entry
status: recorded
goal_id: GOAL-011-pagination-page-size-contract
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# E-003 · 回归与防复发（C3）

## 新增/修正的断言

| 层 | 位置 | 断言 |
|----|------|------|
| 单元 | `renderer/resource.test.ts` | 「省略默认值」改用 `DEFAULT_PAGE_SIZE` 断言；新增「`pageSize=10` 现在必须序列化」（旧断言以字面量 10 锁住了缺陷） |
| 组件 | `renderer/schema-table.test.tsx` | mock 的缺省 pageSize 由 10 改为 **20**（真实服务端契约）；断言下拉默认显示 20、首个请求不带 `pageSize=`；选 10 时**请求含 `pageSize=10` 且表格渲染 10 行**；新增「跳转确认按钮文案为 Go」用例 |
| 组件 | `renderer/saved-views.ui.test.tsx` | 保存视图的默认 pageSize 断言改用 `DEFAULT_PAGE_SIZE`（不再编码旧字面量 10） |
| 结构守卫（新） | `src/pagination-size-contract.guard.test.ts` | 直接解析 `apps/api/internal/handler/resources.go` 的 `DefaultPageSize`，断言前端 `DEFAULT_PAGE_SIZE` 与之**相等**；断言 10/50/100 必须上线路而默认值可省略、且默认值不得再等于 10；断言两个 catalog 都有 `feedback.jumpToPage` 且与 `feedback.search`、`feedback.goToPage` 均不同值，zh = 跳转 / en = Go |
| 浏览器 | `e2e/list-visual-surface.spec.ts`（新增第 3 个用例） | 真实 API + Chromium：下拉显示 **20**；首个列表请求**不带** `pageSize`（服务端默认生效）；选择 10 后**确实发出** `pageSize=10` 的请求且控件显示 10；跳转确认按钮文案匹配 `^(跳转\|Go)$`，表单标签仍为 `^(跳至页\|Go to page)$` |

## 复跑

| 检查 | 结果 |
|------|------|
| `npx vitest run src/pagination-size-contract.guard.test.ts src/renderer/schema-table.test.tsx` | **44 passed** |
| `npm test`（全量 Vitest） | **113 文件 / 1434 测试全部通过** |
| `npm run typecheck`（`tsc -b` + `tsc -p e2e/tsconfig.json`） | exit **0** |
| `npm run test:e2e -- list-visual-surface`（`APP_PROFILE=admin`） | **3 passed**（50.8s，含新增用例） |
| 同上（`APP_PROFILE=mvp`） | **3 passed**（1.2m） |
| `git diff --check` | 通过 |

## 说明

- 浏览器层断言是必要的：该缺陷的可见表现（显示值与生效值不一致）在 jsdom 中不可观察，而既有 jsdom mock 恰好把错误契约写死。新增 e2e 用例直接检查真实请求是否带上 `pageSize=10`，与真实 Go 服务的默认值 20 对齐。
- `e2e/list-visual-surface.spec.ts` 原为 GOAL-009 的交付物；本目标在其中追加一个用例（未改动既有两条），并在文件头的覆盖表中登记 GOAL-011 两行，与 GOAL-010 修改 GOAL-008 守卫文件的先例一致。

C3 完成。C4（self 审计与投影）见 `A-001`、`E-004`。

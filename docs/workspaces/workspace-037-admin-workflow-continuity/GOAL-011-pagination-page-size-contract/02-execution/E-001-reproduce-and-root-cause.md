---
id: E-001-reproduce-and-root-cause
doc: execution-entry
status: recorded
goal_id: GOAL-011-pagination-page-size-contract
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# E-001 · 复现与根因定位（C1）

## 事实

2026-09-18，按用户报告复现两个缺陷并定位根因（不修改实现，先得到失败证据）。

### 1. 复现（改动前，测试失败）

在 `apps/web/src/renderer/schema-table.test.tsx` 补两条断言（默认显示应与生效值一致 = 20；跳转按钮文案应为 Go），把 mock 的「缺省 pageSize」从 10 改为 20（即真实服务端契约）。运行结果：

```
× switches the per-page size (resets to page 1) and jumps to a page
  → expected '10' to be '20'
× keeps the jump confirm button labelled as a jump, not a search
  → expected 'Search' to be 'Go'
Tests  2 failed | 39 passed (41)
```

两条失败正是用户描述的现象：控件显示 10、服务端生效 20；按钮文案为「搜索 / Search」。

### 2. 根因

| 环节 | 事实 | 位置 |
|------|------|------|
| 服务端默认 | `DefaultPageSize = 20` | `apps/api/internal/handler/resources.go:40` |
| 前端默认 | `DEFAULT_PAGE_SIZE = 10` | `apps/web/src/renderer/resource.ts:3` |
| 参数省略规则 | `pageSize !== DEFAULT_PAGE_SIZE` 才写入查询串 | `resource.ts:256` |
| 初始 query | `{ page: 1, pageSize: 10 }` | `schema-table.tsx:615` |
| 下拉显示值 | `query.pageSize ?? 10` | `schema-table.tsx:1501` |
| 跳转按钮文案 | `t("feedback.search")` | `schema-table.tsx:1589` |

推论链：默认 query 为 10 → 请求省略 `pageSize` → 服务端按 20 返回（**显示 10 / 实际 20**）；用户先选 20（发 `pageSize=20`）再选 10 → 10 等于前端默认值 → 再次省略 → 服务端仍 20（**10 不生效**）。

### 3. 为何既有测试未拦下

`schema-table.test.tsx` 的 mock 把「未带 pageSize」解释为 10（`Number(url.searchParams.get("pageSize") ?? "10")`），即把错误契约当成服务端行为；`resource.test.ts` 也以字面量 10 断言「省略默认值」。两者都锁住了缺陷本身（同类问题在 R6 曾以「测试固化回归」的形式出现过）。

### 4. 其它表面核对（`I-011-003`）

- `components/notification-center.tsx`：自有 `DEFAULT_PAGE_SIZE = 10`，但请求**显式**写入 `pageSize`（`params.set("pageSize", ...)`），行为不受本缺陷影响 → 不改。
- `renderer/resource.ts` 的 `DISPLAY_LIST_QUERY` / `EMPTY_RESOURCE_LIST`（`pageSize: 100`）为展示用常量，不参与请求 → 不改。
- `data-permission-scopes.tsx`、`invite-issue-card.tsx`、`telegram-admin-tab.tsx` 显式请求 `pageSize=100`（取全量做选择器）→ 不改。
- `renderer/render.tsx:1456` 与 `components/activity-export.tsx:35` 以字面量 `10` 作兜底 query——在当前常量下恰好被省略，一旦默认值改为 20 就会**显式发出 10**，属必须一并清理的同类耦合（见 `D-001` §2）。

C1 完成；`I-011-001`、`I-011-002`、`I-011-003` 为 `verified`。C2 见 `E-002`。

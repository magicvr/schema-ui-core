---
id: E-002-fix-implementation
doc: execution-entry
status: recorded
goal_id: GOAL-011-pagination-page-size-contract
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# E-002 · 修正实现（C2）

## 实现

按 `D-001` 执行，共 6 个文件（4 个源文件 + 2 个 i18n catalog）：

| 文件 | 变更 |
|------|------|
| `apps/web/src/renderer/resource.ts` | `DEFAULT_PAGE_SIZE` 10 → **20**，并写明它必须等于服务端 `handler.DefaultPageSize`（因为等于默认值时会省略参数）、以及历史偏差的指向（GOAL-011） |
| `apps/web/src/renderer/schema-table.tsx` | 初始 query 与下拉回退值改用 `DEFAULT_PAGE_SIZE`（不再写死 10）；跳转按钮文案 `t("feedback.search")` → `t("feedback.jumpToPage")` |
| `apps/web/src/renderer/render.tsx` | 搜索表单写过滤器时的兜底 query `pageSize: 10` → `DEFAULT_PAGE_SIZE` |
| `apps/web/src/components/activity-export.tsx` | 导出兜底 query `pageSize: 10` → `DEFAULT_PAGE_SIZE` |
| `apps/web/src/i18n/messages/zh-CN.json` | 新增 `"feedback.jumpToPage": "跳转"` |
| `apps/web/src/i18n/messages/en-US.json` | 新增 `"feedback.jumpToPage": "Go"` |

未改 `apps/api`（服务端 20 是既有正确行为）、未改分页尺寸选项集合（10/20/50/100）、未改跳转校验逻辑与表单 `aria-label`（仍为「跳至页 / Go to page」）。

## 效果（服务端契约视角）

| 场景 | 修正前 | 修正后 |
|------|--------|--------|
| 首次加载 | 显示 10；请求省略 pageSize → 服务端 20 | 显示 **20**；请求省略 pageSize → 服务端 20（名实一致） |
| 选 20 | 显式发 `pageSize=20` | 省略（等于默认）→ 服务端 20，行为不变 |
| 选 10 | 省略 → 仍是 20（**不生效**） | 显式发 **`pageSize=10`** → 10 生效 |
| 选 50 / 100 | 显式发送 | 显式发送（不变） |
| 从 20 切回 10 | 不生效 | 生效 |
| 搜索表单提交（兜底 query） | 兜底 10 被省略 → 保持服务端 20 | 兜底 20 被省略 → 保持 20（不引入 10 的静默切换） |
| 跳转确认按钮 | 搜索 / Search | **跳转 / Go** |

C2 完成；`I-011-004` 为 `verified`。C3（回归与防复发）见 `E-003`。

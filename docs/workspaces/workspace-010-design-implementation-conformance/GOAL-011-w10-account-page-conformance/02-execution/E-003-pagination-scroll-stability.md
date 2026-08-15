---
id: GOAL-011-w10-account-page-conformance
doc: execution
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
---

# E-003 · 列表翻页滚动位置保持（用户体验修复）

用户 2026-08-15 反馈：列表翻页时页面视觉定位被强制回到最顶端，体验不佳；期望翻页只变动列表数据、不改变页面滚动位置。

## 根因

`SchemaTable` 的 fetch effect 在**每次 query 变化**（翻页/筛选/排序）时无条件 `setLoading(true)`：

- `DataTable` loading 分支渲染 3 行 skeleton 占位，**旧表格整体卸载**（10 行 → 3 行，高度骤变）
- 新数据到达后表格重新挂载（高度恢复）
- 列表区域高度变化 + 内容整体替换 → 浏览器滚动锚点（scroll anchoring）失效，视口内容被推离 → 用户感知"回顶/跳动"
- 仓库无任何显式 `scrollTo`/`scrollTop` 代码（grep 证实），滚动重置为 DOM 替换副作用

## 修复

`apps/web/src/renderer/schema-table.tsx`（fetch effect）：

- **已有数据时刷新不再进入 loading**：`if (list === null) setLoading(true)`——只有首载（无任何行）显示 skeleton
- 翻页/筛选/排序时旧列表原位保留渲染，新数据到达后**行数据原位替换**（同高度、同结构），浏览器 scroll anchoring 维持视口锚点 → 滚动位置稳定
- 写操作后的 `reloadList` 同样受益（旧列表保持到新数据到达）

## 验证

- 新增回归测试（schema-table.test.tsx "keeps the current rows rendered while paginating"）：可控 fetcher 挂起第 2 页响应，断言翻页在途时**无 loading skeleton**（`data-table-presentation="loading"` 为 null）且旧行仍在；resolve 后新行原位出现（`tbody tr` 行首正则区分页）
- Web 全量 **986/986**（+1）、tsc 0 错误；Go 无改动

## 状态

- 功能语义（翻页/筛选/排序/首载 skeleton/失败 alert）不变，仅刷新呈现路径优化。
- I-001 / I-003（参考样式对齐）仍 open；本目标不变门。

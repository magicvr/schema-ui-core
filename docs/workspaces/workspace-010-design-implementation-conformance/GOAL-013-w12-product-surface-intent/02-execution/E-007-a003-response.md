---
id: E-007-a003-response
doc: execution-entry
goal: GOAL-013-w12-product-surface-intent
date: 2026-08-16
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# E-007 · A-003（列表筛选/表单控件视觉专项）响应与 F-001～F-003 闭合

## 事实

### F-001（recommended · med）· 搜索表单缺少横向工具栏与操作对齐 —— **fixed**

- `render.tsx` `FormInner`：`mode: "search"` 表单改为轻量卡片容器（`rounded-lg border bg-card/60 p-3.5`），字段区由 `FormControls` 以响应式自动网格渲染（1→5 列自适应，`items-end` 对齐）；Search 图标主按钮 + Reset（`RotateCcw`，outline）按钮组通过新的 `actionSlot` 注入网格尾部，与字段同基线。
- Reset 一键清空全部条件并立即重新查询。
- 修复衍生缺陷：`searchFormSubmit` 全部条件清空时未删除继承的旧 filters（spread 残留），已显式 `delete next.filters`。

### F-002（recommended · low）· 基础控件视觉质感 —— **fixed**

- `form-controls.tsx` `SelectField`：原生 select 增加 `appearance-none pr-8` + 绝对定位 `ChevronDown`（全站生效，含 modal/设置页）。
- `BaseInput`：搜索模式关键词输入带放大镜前缀图标（`pl-8`）与有值时的 X 清空按钮（`feedback.clearSearch` aria-label）。
- Checkbox / Radio / Switch（BooleanField）：统一 `accent-primary` + `size-4` 品牌强调。
- 12 页 schema JSON 零改动（渲染器侧透明实现）。

### F-003（recommended · low）· 缺激活筛选可视化与快速清除 —— **fixed**

- `FormInner`：`data-filter-chips` 区渲染非空条件的胶囊标签（字段标签=选项标签，select 值经 `optionList` 映射为人类可读标签），单项 ✕ 立即清除并重查；`清除全部筛选` 链接一键复位。
- i18n：`feedback.reset` / `feedback.clearFilters` / `feedback.removeFilter` / `feedback.clearSearch`（en/zh）。

### 蓝图未落地项（非 finding，留后续波次）

- 折叠/展开高级筛选（字段 ≥4 时）：当前 12 页最多 3 个筛选字段（users），未触发；按 Phase 2 延后。
- DateRange 范围胶囊美化：dateRangePicker 现状保留，Phase 2。


## 修订（用户反馈 2026-08-16）：搜索按键与关键词输入框成对并排

- 用户指出：上一版把「搜索」按钮放在网格末尾按钮组，与对应文本输入框脱开；应按「**有几个文本输入框就有几个搜索按键，且两两成对并排**」布局。
- 实施：`FormControls` 新增 `searchButtonSlot`，search 模式下每个 `input` 字段与其搜索按钮同处一个网格单元（`flex items-end gap-2`，按钮 `h-9 shrink-0` 与输入框同高同基线）；`actionSlot` 只保留 Reset 按钮（网格尾部）。
- 当前 12 页搜索表单每页恰好一个关键词输入，因此每页一个搜索键、紧贴输入框右侧；未来页面若含多个文本输入框，各输入框自动各配一个搜索键（提交同一表单条件）。
- 测试：`search-form-filters.test.tsx` 增加成对断言（submit 按钮与 input 位于同一 grid cell）；全量回归绿。
- 验证：`npx vitest run` 1027/1027；`tsc -b` 0。

## 验证

- `npx vitest run`（apps/web）：全量结果见套件日志（1027 基线 + A-003 扩展断言）。
- `tsc -b`：0 错误。
- 测试扩展：`search-form-filters.test.tsx` 增加 chips 断言 + Reset 清空断言（含 q 与 enabled 均被移除）。

---
id: E-002
doc: execution
status: recorded
parent: GOAL-014-w13-settings-tabs-and-topbar
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# E-002 · S2 实施：T-01～T-04

## T-01 · 设置页功能单元 Tabs

- `apps/api/internal/modules/settings/schema/settings.json`：body 由「单 section 纵向堆叠」重构为 `section[text(描述) → tabs(五个功能单元) → actionButton(恢复默认)]`；五个 tab section 各包一个既有表单（常规/品牌/本地化/外观/安全），表单字段与动作契约零改动。
- `apps/web/src/i18n/messages/en-US.json` + `zh-CN.json`：新增 `schema.settings.tab.{general,branding,localization,appearance,security}` 五组键（双语文案表键集一致）。

## T-02 · 移动端品牌条

- `apps/web/src/app/App.tsx`：新增局部组件 `BrandLink`（logo + 站点标题 + 副标题，含 light/dark logo 变体与首页跳转），header 内两处复用：
  - <lg：`data-shell-region="mobile-brandbar"` 独立品牌行（`lg:hidden`），功能区主行不再被品牌块挤占；
  - ≥lg：品牌链接 `hidden lg:flex` 回到单行布局（与汉堡/抽屉同断点约定）。
- 共享 `handleBrandClick`（homePageRef 首页跳转）；两条品牌链语义一致。

## T-03 · 搜索框组【文本框+搜索键】贴合恒同行

- `apps/web/src/renderer/form-controls.tsx`：`FieldControl` search 配对容器由 `flex items-end gap-2` 改为 `flex flex-nowrap items-end`（零间隙）；`BaseInput` 新增 `paired` 属性，配对时 input 右内角改直角（`rounded-r-none`）。
- `apps/web/src/renderer/render.tsx`：搜索按钮 `rounded-l-none -ml-px`——按钮左缘与 input 右缘 1px 叠边，视觉上贴合为单组件；`flex-nowrap` + 按钮 `shrink-0` + 输入侧 `min-w-0 flex-1` 保证任何页面宽度下恒在同一行。

## T-04 · 顶栏亮暗/语种按键对调

- `apps/web/src/app/App.tsx`：功能区顺序由 `[汉堡, 语种, 亮暗, 铃, 用户]` 改为 `[汉堡, 亮暗, 语种, 铃, 用户]`（亮/暗在左、语种在右）；组件自身 aria-label/行为不变。

## 产物清单

- `apps/api/internal/modules/settings/schema/settings.json`（T-01）
- `apps/web/src/i18n/messages/en-US.json`、`apps/web/src/i18n/messages/zh-CN.json`（T-01）
- `apps/web/src/app/App.tsx`（T-02/T-04）
- `apps/web/src/renderer/form-controls.tsx`、`apps/web/src/renderer/render.tsx`（T-03）

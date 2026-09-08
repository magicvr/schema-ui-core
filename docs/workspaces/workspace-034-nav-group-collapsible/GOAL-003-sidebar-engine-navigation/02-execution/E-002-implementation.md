---
id: E-002-implementation
goal_id: GOAL-003-sidebar-engine-navigation
doc: execution-entry
status: recorded
date: 2026-09-08
parent: GOAL-001-nav-group-collapsible
version: 0.1.0
---

# E-002 · Sidebar Engine 与通用详情抽屉实施

## 已发生事实

- `apps/api/kernel` 的 `NavigationContribution` / `NavigationGroup` 新增可选 `Secondary` 表现元数据，并增加 trim/空白校验；同组元数据 equality 覆盖该字段。`authsession` navigation checksum 保持不包含该表现字段。
- Dashboard Provider 将 `menu_dashboard` 显式注册到 `workspace` Group，并注册 `01` 页面副字符；默认产品分组注册 `WORKSPACE`、`IAM`、`CMS`、`OPS`、`COMMS`、`COMMERCE`，Data dictionary 注册 `SQL` 页面副字符。
- Manifest 聚合新增全量已注册 presentation 投影，校验 fragment 必须属于 enabled module，并把 secondary 写入协议允许的 literal `label` fallback（`主文案 · secondary`）；Web navigation projection 拆回独立 secondary。未注册节点不生成链接/副字符，Examples authored group 保持独立，legacy 未分组 Dashboard 仍保持向后兼容的顶层顺序。
- `apps/web/src/app/App.tsx` 的 desktop sidebar 与 mobile drawer 改为紧凑树形层级、组标题 secondary、叶子右侧 secondary、active/hover/focus semantic surface；没有新增闪烁选中光点。
- `recordView` 重构为通用 Drawer/Sheet presentation：固定右侧 `460px`、blur backdrop、h-14 header、p-6 scroll body、Properties 字段卡片与 divide 分隔、长对象/数组可读格式化、通用完成按钮；selection 模式增加 Esc、Tab trap、焦点恢复和 body scroll lock，static record 保持非 modal。未引入用户/ROOT/MFA/Quick Ops 等参考页业务内容。
- 更新两种 dogfood Manifest fixture、导航/抽屉测试、Playwright 分组展开辅助与模块贡献 Playbook；未改 pinned upstream schema/provenance。

## 证据路径

- `apps/api/kernel/contribution.go`
- `apps/api/internal/manifest/manifest.go`
- `apps/api/modules/dashboard/provider.go`
- `apps/web/src/app/navigation.ts`
- `apps/web/src/app/App.tsx`
- `apps/web/src/renderer/render.tsx`
- `docs/architecture/module-contribution-playbook.md`
- `apps/api/internal/manifest/manifest_test.go`
- `apps/web/src/app/nav-groups.test.tsx`
- `apps/web/src/renderer/visual-fidelity.test.tsx`


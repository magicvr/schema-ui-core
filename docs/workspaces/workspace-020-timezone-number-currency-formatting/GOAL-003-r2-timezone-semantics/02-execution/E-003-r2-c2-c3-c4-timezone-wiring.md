---
id: GOAL-003-r2-timezone-semantics
doc: execution-entry
record_id: E-003
status: recorded
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
---

# E-003 · C2/C3/C4 实施（头部时区选择 + 站点默认接入 + 统一语义接线）

## 2026-08-26

### 已发生事实

1. **C2 用户级覆盖 UI**：`apps/web/src/components/timezone-switcher.tsx`（仿 `locale-switcher`：图标按钮 + 下拉 menu；选项 = auto + 常用 IANA 集 `Asia/Shanghai / Asia/Tokyo / America/New_York / Europe/London / UTC`；选中打勾；outline click / Escape 关闭）。挂载于 `apps/web/src/app/App.tsx` 头部（`<LocaleSwitcher />` 右侧，同 locale 通道）。catalog 新增 `timezone.switcher.label` / `timezone.switcher.auto`（zh-CN + en-US 双语）。
2. **C3 站点默认接入**：`apps/web/src/i18n/runtime.tsx` 的 `systemDefaultUrl`（`/api/branding`）fetch 增读 `siteTimezone` → `fetchedSiteTimezone`；`I18nProvider` 新增 `siteTimezone` / `storedTimezone` / `detectTimezone` 输入（测试缝），生效时区解析输入 `siteDefault = siteTimezone ?? fetchedSiteTimezone`（合同 §2 L3）。Localization tab 字段未改动（`siteTimezone` 已有，schema 未变）。
3. **C4 统一语义接线**：`I18nState` 暴露 `timezone`（生效时区）/ `timezonePreference` / `setTimezonePreference`（单通道持久化，auto = 移除 key）；`formatDate` 默认注入生效时区（调用方显式 `options.timeZone` 优先，合同「同一语义双向」）。
4. **快测**：`src/i18n/runtime-timezone.test.tsx`（6 用例：L1/L3/L4 接线、持久化 round-trip + 格式翻转、显式 timeZone 覆盖、C3 fetch 接线）；`src/components/timezone-switcher.test.tsx`（4 用例：菜单/选项/勾选、持久化、auto 移除、外部点击关闭）。合计 **10/10 pass**。

### 证据

| 主张 | 路径 / 命令 / commit |
|------|----------------------|
| C2 组件与挂载 | `apps/web/src/components/timezone-switcher.tsx`；`apps/web/src/app/App.tsx` |
| C3 站点默认读取 | `apps/web/src/i18n/runtime.tsx`（`fetchedSiteTimezone` / `siteTimezone` prop） |
| C4 统一语义 | `apps/web/src/i18n/runtime.tsx`（`formatDate` 默认生效时区；`timezone` 状态） |
| 快测 10/10 | `npx vitest run src/i18n/runtime-timezone.test.tsx src/components/timezone-switcher.test.tsx`（exit 0） |
| 双语文案 | `apps/web/src/i18n/messages/zh-CN.json` / `en-US.json`（`timezone.switcher.*`） |
| 合同条款 | `GOAL-002-r1-contract-freeze/01-decision/D-001` §2 / §4.2 |
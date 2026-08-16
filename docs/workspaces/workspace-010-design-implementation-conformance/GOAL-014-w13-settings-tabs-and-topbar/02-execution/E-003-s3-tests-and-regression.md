---
id: E-003
doc: execution
status: recorded
parent: GOAL-014-w13-settings-tabs-and-topbar
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# E-003 · S3 测试更新与回归

## 测试更新（与 T-01～T-04 配套）

- `apps/web/src/app/startup-config.test.tsx`：设置页用例改为 Tabs 形态——断言五个功能单元 tab 标签 + 恢复默认常驻；逐 tab 点击验证表单预填（站点标题/默认语种 zh-CN/默认主题 dark/验证码）；权限门禁用例改为逐 tab 断言保存按钮禁用 + 恢复默认禁用；测试挂具补 `/api/captcha/settings` 返回。
- `apps/web/src/app/shell.test.ts`：新增静态守卫——移动端品牌条（`data-shell-region="mobile-brandbar"` + `lg:hidden`）与桌面品牌链（`hidden … lg:flex`）；ThemeToggle 位于 LocaleSwitcher 之前。
- `apps/web/src/renderer/search-form-filters.test.tsx`：配对断言扩展——按钮含 `-ml-px` + `rounded-l-none`、输入含 `rounded-r-none`（贴合成单组件）。
- `apps/web/e2e/localization.spec.ts`：设置页四类断言改为点击 tab 后断言；站点标题投影断言改 `.last()`（移动端品牌条与桌面品牌链并存）。
- `apps/web/e2e/shell.spec.ts`：新增 W13 断言——亮暗按键在语种按键左侧（boundingBox x 比较）；390×844 视口下品牌条可见、汉堡在功能区；站点标题可见性断言改 `.last()`。

## e2e 陈旧断言修复（W11/W12 遗留，与本波并存记录）

首次全量 e2e 发现 3 条陈旧断言（非本波引入，W12 关门时未跑 e2e）：

| 规格 | 陈旧断言 | 修复 | 根因 |
|------|----------|------|------|
| shell.spec / localization.spec | `getByRole("link", … 设置/Settings)` | 打开顶栏用户下拉后断言 `menuitem` | W12 T-01：设置移入用户下拉（button，非 sidebar link） |
| schema-crud.spec | `getByLabel("Read users")` | `getByLabel("users.read")`；角色授权对话框改 `getByLabel("E2E Support").check()` | W11 U-01/U-02：`/api/permissions` 返回 key 标签；roles 授权改为 checkboxGroup（optionsSource /api/roles） |
| schema-crud.spec | 行内 `Password`/`Delete` 按钮 | 超过前两个的行操作从 `⋯ More actions` 菜单取 | W11 U-05：行操作收纳（roles 行仅 2 个操作仍行内） |
| shell.spec | `Open menu` | `Open navigation menu`；`getByText(siteTitle).first()` → `.last()` | 汉堡 aria-label 实值；W13 品牌条使标题双现 |

## 回归证据（全绿）

- **vitest**：63 文件 **1029/1029** 通过（W12 基线 1027 + 新增 2 静态守卫）。
- **tsc -b**：0 错误。
- **Go**：`go test ./...` 全量 0 FAIL（settings 模块含 schema 服务回归）。
- **e2e（Playwright + 真实 Go API + SQLite 新鲜库）**：admin profile 8/8 通过（1 跳过 mvp 用例）；mvp profile 8/8 通过（1 跳过 admin 用例）。W13 交互均有真实浏览器证据：设置页 Tabs 切换与保存、顶栏按键顺序、移动端品牌条、用户下拉菜单。

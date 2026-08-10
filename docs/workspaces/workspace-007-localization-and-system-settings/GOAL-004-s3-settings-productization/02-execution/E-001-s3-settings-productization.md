---
id: E-001
doc: execution
title: S3 · 四类系统设置产品化 + 公开启动配置 + 权限/审计/刷新闭环
status: recorded
parent: GOAL-004-s3-settings-productization
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

# E-001 · S3 系统设置产品化（2026-08-09）

## 事实

- **API（Go）**：
  - 迁移 `site_settings_v2`（全局 version 10，模块自有 0008→0010 重编号避免与 operation_log_settings 冲突）：新增 6 列（logo_url_light/logo_url_dark/favicon_url/default_locale/site_timezone/default_theme，TEXT NOT NULL DEFAULT ''）。
  - 仓库：`SiteSettings` 扩展 8 字段；`PatchSiteSettings` 8 指针字段级合并（未提交字段不覆盖）、空串清空语义、`ResetSiteSettings` 恢复冻结默认（title=Schema UI Core、locale/theme/timezone=auto、品牌 URL 清空）；校验：locale/theme 枚举、IANA 时区（`time.LoadLocation`，无效 → `ErrInvalidSiteTimezone` 且原子拒绝不清空）。
  - 配置：`settings.branding` namespace Validate 扩展（additive 字段 + 时区/枚举校验）。
  - Handler：`GET /api/branding` 扩展为公开启动配置（siteTitle/logoUrl/logoUrlLight/logoUrlDark/faviconUrl/defaultLocale/supportedLocales/siteTimezone/defaultTheme，旧字段兼容）；PATCH 扩展 8 字段；新增 `POST /api/settings/{id}/reset`（settings.write）；新错误码 `INVALID_DEFAULT_LOCALE`/`INVALID_DEFAULT_THEME`/`INVALID_TIMEZONE`（列入 D-002 附录 A 枚举族）；操作日志 `settings.update` 保留并扩展 action 字段。
  - Kernel/profile：`admin.settings` 声明路由 + reset 贡献键（provider 与 plan 描述符同步）。
- **Web**：
  - `theme.ts`：`resolveTheme` 优先级扩展（显式 → 系统默认非 auto → OS）；`applySystemDefaultTheme`。
  - `branding.ts`：`Branding` 扩展 8 字段；`fetchBranding` 解析；`applyDocumentBranding` faviconUrl（回退 logoUrl）；`defaultBranding`。
  - `runtime.tsx`：`I18nProvider` 新增 `systemDefaultUrl`（拉取公开启动配置注入系统默认语种）；`main.tsx` 接入 `/api/branding`。
  - Shell：浅/深色 Logo 变体按主题 CSS 切换（`dark:hidden`/`dark:block`）；品牌加载后应用系统默认主题 + favicon；登录页同。
  - Renderer：`actionButton` 在默认应用路径接通 Schema CRUD 执行器（gate→confirm→request），label 经 `resolveTextProp` 双语解析。
  - Settings 页 schema 重构：四类工具条动作（General/Branding/Localization/Appearance 编辑 modal + Restore defaults 恢复默认 request）+ 8 字段只读表 + 说明；动作 URL 用具体路径 `/api/settings/default`（pageTrigger/form 构造器拒绝 `{id}` 模板，S3 核对发现并修正）；全部文本带 `*Key`（catalog 新增 ~30 键/语种）。
- **预览语义（D-001 冻结，v1）**：表单值即时回显 + 保存后立即投影（无刷新）——Shell/登录页/公开启动配置经配置刷新事件重拉生效；恢复默认走独立 request 端点。
- **验证**：Go 全包通过（新增 `TestBrandingVp007StartupFieldsAndPatch`/`TestSettingsValidationAndReset`/仓库 VP-007 字段合并/清空/重置/校验测试；store 迁移测试更新至 10 条）；vitest **706/706**（38 文件，新增 S3 测试 16：startup 解析/默认回退/favicon、系统默认主题优先级与显式优先、Shell 浅深 Logo+favicon、provider systemDefaultUrl 语种注入与显式优先、四类设置面 zh 渲染、General 保存真实 PATCH、权限失败禁用、恢复默认确认→reset 请求）；`npm run build` exit 0；输出捕获 `{SCRATCH}/unit-s3-web.log`、`{SCRATCH}/unit-s3-api.log`。

## 产物

| 路径 | 说明 |
|------|------|
| `apps/api/internal/modules/settings/{migration,repository,configuration,provider,schema}/*` | 四类字段持久化/校验/路由/schema 页（C1/C4/C5） |
| `apps/api/internal/handler/settings.go` + tests | 公开启动配置 + PATCH/reset + 错误码 + 操作日志（C2/C5） |
| `apps/api/internal/kernel/profile.go` | admin.settings 贡献键（reset 路由） |
| `apps/web/src/{theme,branding,runtime,main,App,LoginPage}*` | 系统默认主题/语种/品牌资产投影（C2/C3） |
| `apps/web/src/renderer/render.tsx` | actionButton 接通执行器 + 双语 label |
| `apps/web/src/app/startup-config.test.tsx` | S3 测试 13 项 + theme/runtime 扩展测试 |

## 里程碑 checkpoint

- commit：`cf1e78e`（2026-08-09，S3；owned paths = 上述全部路径 + 本目标治理文档，显式 `git add` 无 `-A`）。

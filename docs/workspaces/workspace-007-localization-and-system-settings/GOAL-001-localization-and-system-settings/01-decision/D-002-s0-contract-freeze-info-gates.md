---
id: D-002
doc: decision
title: S0 · 差距盘点与契约冻结：关闭 I-L10N-001～005、冻结 F-V029 覆盖表
status: accepted
parent: GOAL-001-localization-and-system-settings
created: 2026-08-09
updated: 2026-08-18
version: 0.1.1
---

# D-002 · S0 契约冻结（用户书面裁决 2026-08-09）

## 触发

S0 阶段启动：基线全绿（`go test ./...` 全包通过；vitest 629/629），完成现状盘点后，`I-L10N-001`～`005`（required）到达「最晚需要阶段」。按 P-004/P-005，五条门禁结论全部经用户书面选择（2026-08-09 会话裁决，逐条留痕如下），本决策冻结契约并关闭门禁。

## 决定（每项 = 用户选定方案 + 冻结语义）

### I-L10N-001 · Schema 驱动文本本地化 = 前端 key 解析（用户选定）

- **兼容盘点（事实）**：`docs/schemas/app-manifest.schema.json`（v2.7.0 vendored）已在 pages/navigation 声明 `titleKey`/`labelKey`；本地 `docs/schemas/component-registry.json`（VP-006 冻结的机器可读组件契约）已为组件声明 `labelKey`/`titleKey`/`contentKey`；上游 `node.schema.json` 的 `props` 为开放对象（仅禁 CSS 名），本地 registry 按组件枚举 props。因此 **key 字段是既有兼容面，不是私有扩展**。
- **冻结策略**：前端翻译 catalog（`zh-CN`/`en-US` 纯数据文件）解析 `*Key` 字段；Renderer 与导航在渲染时解析；缺失 key **可观察**（`schema-ui:missing-translation` 事件 + 测试断言）且**安全回退**（当前语种 → `en-US` catalog → 字面文本 → key 本身），不渲染为空、不抛异常、不阻断操作。
- **边界**：不引入服务端 locale overlay；不改写上游 `page.schema.json`/`node.schema.json` 语义；协议字面文本（`label`/`title`/`content` 等）保持 en-US 规范原文。**S2 修正（A-002 补充核对）**：`component-registry.json` 为上游 pin 制品（I-PROTO-004 sha256 校验），**不改写**；其已声明的 `labelKey`/`titleKey`/`contentKey`/`options.labelKey` 直接使用；registry 未声明的缺口字段（`submitLabelKey`/`confirmKey`/`textKey`/`placeholderKey`）作为**本地页面文档约定**在 S2 D-001 登记（上游 `props` 开放，文档级合法），Renderer 解析并遵循冻结回退链。

### I-L10N-002 · 用户语种持久化 = localStorage 单通道（用户选定）

- **冻结**：显式语种持久化在 `localStorage["schema-ui:locale"]`（`zh-CN` | `en-US` | 移除 = `auto`），与既有 `theme` 机制同模式；**匿名与登录后共用同一选择**，登出/登录不清除、不合并、不新增账号资料字段、不做跨设备同步（未来波次候选）。
- **优先级（登录前后一致）**：用户显式选择 → 系统指定默认语种（Settings Localization `defaultLocale`，非 `auto` 时）→ 浏览器语言偏好（`auto`）→ `en-US` 安全回退。
- **入口**：Shell / 用户菜单内的语种切换器，匿名登录页同样可达（不需要任何设置权限）。

### I-L10N-003 · 公开启动配置 = 兼容扩展 `/api/branding`（用户选定）

- **冻结**：`GET /api/branding`（public，no-cache）additive 扩展字段：`defaultLocale`（`auto`|`zh-CN`|`en-US`）、`supportedLocales`（`["zh-CN","en-US"]`）、`defaultTheme`（`auto`|`light`|`dark`）、`siteTimezone`（IANA 名 | `UTC` | `auto`）；保留既有 `siteTitle`/`logoUrl` 与全部缓存/配置变更语义（`X-Schema-UI-Config-Changed` 沿用 `settings.branding` 命名空间，配置刷新信号不变）。
- **边界**：不新建公开 bootstrap 端点；无缓存头变更；旧消费端（仅读 siteTitle/logoUrl）不受影响。

### I-L10N-004 · 错误 envelope = 路径 (a) 有界服务端 locale 协商（用户选定）

- **冻结**：错误码保持稳定且可机读；`writeError` 兼容扩展 envelope `{error, message, messageKey?, params?}`（`message` 保留英文原文，旧客户端兼容）；已编目错误（认证/验证/设置/资源/上传，见 S0 盘点约 35 码）按 `Accept-Language` 首支持语种返回对应语种 `message` 并声明 `Content-Language` 头；未编目/`INTERNAL` 错误保持英文通用文案（安全回退、不泄露诊断）；前端仍按码/key/参数本地化保底（不可降级）。
- **证据要求**：错误码契约测试（码集合不变）+ 认证/验证/设置错误的协商、`Content-Language`、失败回退证据（S4 交付）。

### I-L10N-005 · 时区语义 = UTC 存储 + 显示转换（用户选定）

- **冻结**：存储/API 全部 UTC（现状 `updatedAt` unix + ISO8601 UTC 不变）；Settings Localization 新增 `siteTimezone`（IANA 名 | `auto`，默认 `auto` = 浏览器本地时区显示）；前端按有效 locale + 生效时区显示日期/数字（`Intl.DateTimeFormat`/`NumberFormat`），首版不暴露任意格式模板；无效 IANA 时区 → `400 INVALID_TIMEZONE` 校验错误且**不清空原值**。

## 冻结的其他契约（随本决策）

1. **支持语种**：首发 `zh-CN` + `en-US`；系统默认 `auto`。
2. **F-V029 覆盖表**（固定 UI 面 + 双 Profile Runtime Manifest pageId/schemaUrl 并集 + M1～M4 + 缺失翻译纪律）落 Root `attachments/`（见 [F-V029-coverage-table-s0-freeze.md](../attachments/F-V029-coverage-table-s0-freeze.md)），本波次证据矩阵复用同一分母；覆盖表外用户可见文案只能经用户书面 `accepted-residual` 排除。
3. **Settings 四类字段**：General（`siteTitle`）、Branding（`logoUrl`、`logoUrlLight`、`logoUrlDark`、`faviconUrl`——均为同源路径或 HTTPS URL，不上传；S3 实现预览/校验/清空/恢复默认；Shell 按主题/`prefers-color-scheme` 应用浅/深色 Logo，`faviconUrl` 联动文档 favicon、缺省回退 `logoUrl`）、Localization（`defaultLocale`、`siteTimezone`）、Appearance（`defaultTheme`）；`admin.settings` 编辑面仅 `admin` Profile 暴露；权限沿用 `settings.read`/`settings.write`。（2026-08-09 A-001 F-001 用户书面裁决 **fixed**：补全 Branding 字段，与 VP-007 交付范围对齐。）
4. **审计模式**：S0 契约冻结与 S5 关门 = `independent`（grok CLI，`-m grok-4.5 --effort high` 执行 `/audit`）；常规阶段 `self` 兜底。

## 附录 A · 稳定错误码枚举（A-001 F-005 → fixed，2026-08-09 钉死）

> S4 错误码契约测试以此清单为回归基线：码集合**不得**变更（可新增，不可改语义/复用）。来源：`apps/api` 全仓 `writeError` 字面量与 domain rejection 枚举。

**字面量码（31）**：`EMPTY_SELECTION`、`FILE_NOT_FOUND`、`FILE_TOO_LARGE`、`FORBIDDEN`、`INTERNAL`、`INVALID_BODY`、`INVALID_CREATE_BODY`、`INVALID_CREATE_FIELD`、`INVALID_FILE`、`INVALID_LOGIN_BODY`、`INVALID_LOGO_URL`、`INVALID_LOGOUT_BODY`、`INVALID_PAGE`、`INVALID_PAGE_SIZE`、`INVALID_PATCH_BODY`、`INVALID_PATCH_FIELD`、`INVALID_REFRESH_BODY`、`INVALID_SELECTION_KEY`、`INVALID_SITE_TITLE`、`INVALID_SORT_FIELD`、`INVALID_SORT_ORDER`、`INVALID_UPLOAD`、`LOGIN_FAILED`、`LOGOUT_FAILED`、`REFRESH_FAILED`、`SCHEMA_NOT_FOUND`、`SETTINGS_NOT_FOUND`、`STORAGE_UNAVAILABLE`、`UNAUTHENTICATED`、`UNAUTHORIZED`、`UNSUPPORTED_FILE_TYPE`。

**动态/域码族（8）**：`{RESOURCE}_NOT_FOUND`（如 `ROLE_NOT_FOUND`/`USER_NOT_FOUND`/`CATALOG_NOT_FOUND`）、`USERNAME_TAKEN`、`ROLE_KEY_TAKEN`、`ROLE_IN_USE`、`ROLE_SYSTEM`、`INVALID_ROLE_KEY`、`INVALID_PERMISSION_REF`、`ROLE_GRANT_FORBIDDEN`。

**增量稳定码（2026-08-18，VP-012 R4）**：按本附录“可新增、不可改语义/复用”的规则，追加异步 Job HTTP 码 `JOB_NOT_FOUND`、`JOB_NOT_CANCELLABLE`、`JOB_NOT_RETRYABLE`、`JOB_RESULT_NOT_READY`、`JOB_RESULT_EXPIRED`，以及仅持久化于 Job 终态 representation 的 `JOB_ATTEMPTS_EXHAUSTED`、`JOB_HANDLER_FAILED`。精确 HTTP/状态转换语义见 `docs/workspaces/workspace-012-shared-cross-module-contracts/GOAL-005-r4-async-job-contract/01-decision/D-002-r4-precise-contract.md` §3/§7；原始枚举及既有语义不变。

**编目策略（S4 实施）**：上列全部码为稳定机读码；`INTERNAL` 及未编目码永不携带具体诊断（英文通用文案），编目码按 `Accept-Language` 返回语种化 `message` + `messageKey`/`params`。

## 未选方案

- 服务端 locale overlay（I-L10N-001 未选：schema 缓存/一致性成本高且超出 v2.7.0 语义）。
- 账号资料语种字段 / localStorage+登录合并（I-L10N-002 未选：v1 范围扩张，跨设备同步留后续波次）。
- 新 `/api/bootstrap` 契约（I-L10N-003 未选：第二公开端点 + 消费端迁移）。
- 路径 (b) 仅前端 localize + accepted-residual（I-L10N-004 未选：用户选定有界服务端协商）。
- 服务端本地时间存储（I-L10N-005 未选：时区混淆）。

## 影响

- S1：按 I-L10N-001/002 实施 locale resolver/provider、catalog、切换与格式化。
- S2：按 F-V029 分母双语化固定 UI + 双 Profile page/schema 并集；`titleKey`/`labelKey` 真解析。
- S3：按 I-L10N-003/005 扩展 `/api/branding` 与四类设置。
- S4：按 I-L10N-004 实施有界服务端协商 + 错误码契约测试。
- S5：证据矩阵复用 F-V029 分母；关门审计 + 用户书面确认。

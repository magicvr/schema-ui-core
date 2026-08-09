# F-V029 · 最小可枚举证据面（S0 冻结，2026-08-09）

> **状态**：`frozen`（D-002，用户书面裁决 2026-08-09）。本表是 VP-007 退出判据 2 与 6 **唯一**共用分母；波次内任何用户可见文案面未列入本表 = 必须经用户书面 `accepted-residual` 排除，不得以「代表性」静默收缩。
> **行语义**：分母 = `zh-CN` / `en-US` × `mvp` / `admin` × 匿名 / 已认证。`N/A` 仅限 Profile 不可达单元格并写明模块边界，不算 pass。
> **证据路径**：S2/S5 阶段按单元格填写（测试名、日志、截图路径）；未填前不得宣称该单元格已覆盖。

## 1. 固定 UI 面（双 Profile 共有，代码内文案）

| # | 面 | 内容要点 | 翻译来源 | fallback | 证据路径（S2/S5 填） |
|---|----|----------|----------|----------|----------------------|
| U1 | 匿名登录页 | 标题、用户名/密码标签与占位、登录/登录中按钮、错误提示、开发 seed 提示、品牌标题/Logo | catalog keys | 字面文本 → en-US → key | S2：`i18n/ui-bilingual.test.tsx`（zh/en 登录面、错误码映射、`lang`）；`LoginPage.test.tsx` |
| U2 | Shell header | 应用名、Admin console 副标、移动菜单按钮 aria、登出按钮、用户名展示 | catalog keys | 同上 | S2：`ui-bilingual.test.tsx`（Shell 渲染文本）；`shell.test.ts`（aria 结构性）；`App.integration.test.tsx` |
| U3 | Shell sidebar/top/user 导航 | 分组标签（Workspace、Examples、Admin）、导航项 label（来自 manifest label/labelKey） | catalog + manifest key 解析 | 同上 | S2：`ui-bilingual.test.tsx`（titleKey/labelKey 双语解析、字面回退）；`schema-keys.structural.test.ts`（manifest key 齐全性） |
| U4 | 语种切换器 | 切换器标签、语种名（简体中文/English）、auto 项 | catalog keys | 同上 | S1：`components/locale-switcher.test.tsx`（双语标签、持久化）；`i18n/runtime.test.tsx` |
| U5 | 通用反馈 | loading / empty / error / permission denied / success 文案（含 Schema 页面加载失败面、Route fallback、session 失败提示） | catalog keys | 同上 | S2：`ui-bilingual.test.tsx`（路由回退面、列表空态、分页）；`App.integration.test.tsx`；`visual-fidelity.test.tsx` |
| U6 | 通用组件文案 | table（列头、排序、空态）、search（占位、按钮）、form（校验、必填、提交/取消）、modal（标题、确认/取消）、validation 提示 | catalog keys + schema key 解析 | 同上 | S2：`ui-bilingual.test.tsx`（列头/工具栏/提交按钮/确认框）；`schema-crud.test.tsx`；`data-table.test.tsx` |
| U7 | Admin Settings 四类设置面 | General/Branding/Localization/Appearance 的字段标签、说明、按钮（保存/预览/恢复默认）、校验错误、成功/失败反馈。Branding 字段 = `logoUrl` + `logoUrlLight` + `logoUrlDark` + `faviconUrl`（同源路径或 HTTPS URL，不上传；浅/深色 Logo 按主题应用，favicon 缺省回退 `logoUrl`） | catalog keys + schema key 解析 | 同上 | S2：`schema-keys.structural.test.ts`（settings.json 全文本有 key + 双语 catalog）；settings 页面视觉与四类字段闭环 = **S3**（届时回填） |

## 2. Manifest / Schema 面（S0 冻结 pageId 并集，双 Profile Runtime Manifest）

> 并集 = `mvp`（10）+ `admin`（12）共 **12** 个 pageId；`settings`/`activity` 仅 `admin` 可达（模块边界：`admin.settings`/`admin.activity` 不在 mvp profile 模块集内）。schemaUrl 全部为 `/api/schema/<pageId>`；manifest `title`/`label` 与 schema 内 `label`/`title`/`content`/`fallbackText` 等为协议字面文本，`*Key` 为翻译键。

| pageId | schemaUrl | mvp | admin | 翻译来源 | fallback | 证据路径（S2/S5 填） |
|--------|-----------|-----|-------|----------|----------|----------------------|
| overview | /api/schema/overview | ✓ | ✓ | schema key 解析 | 字面文本 → en-US → key | S2：`schema-keys.structural.test.ts`（全文本有 key + 双语 catalog）；渲染证据 S5 矩阵 |
| admin-list-batch | /api/schema/admin-list-batch | ✓ | ✓ | 同上 | 同上 | S2：同上 + `representative-pages.integration.test.tsx`（批量按钮/确认经 en-US catalog 解析） |
| data-display | /api/schema/data-display | ✓ | ✓ | 同上 | 同上 | S2：`schema-keys.structural.test.ts` |
| data-table | /api/schema/data-table | ✓ | ✓ | 同上 | 同上 | S2：同上 |
| search-form-table | /api/schema/search-form-table | ✓ | ✓ | 同上 | 同上 | S2：同上 |
| form-controls | /api/schema/form-controls | ✓ | ✓ | 同上 | 同上 | S2：同上 |
| form-with-reactions | /api/schema/form-with-reactions | ✓ | ✓ | 同上 | 同上 | S2：同上 |
| form-with-upload | /api/schema/form-with-upload | ✓ | ✓ | 同上 | 同上 | S2：同上 |
| users | /api/schema/users | ✓ | ✓ | 同上 | 同上 | S2：`ui-bilingual.test.tsx`（列头/工具栏/表单/确认框双语、M4 缺失 key）；`schema-crud.test.tsx` |
| roles | /api/schema/roles | ✓ | ✓ | 同上 | 同上 | S2：`schema-keys.structural.test.ts`（44 处 key 齐全）；渲染证据 S5 矩阵 |
| settings | /api/schema/settings | N/A（模块边界：admin.settings 不在 mvp） | ✓ | 同上 | 同上 | S2：`schema-keys.structural.test.ts`；四类字段闭环 S3 |
| activity | /api/schema/activity | N/A（模块边界：admin.activity 不在 mvp） | ✓ | 同上 | 同上 | S2：`schema-keys.structural.test.ts`（16 处 key 齐全） |

Manifest 导航项（sidebar/top/user 的 label/labelKey）与页面 title/titleKey 随上述并集一并覆盖（U3/U6）。

## 3. 固定主流程（M1～M4）

| # | 流程 | 关键断言 | 证据路径（S2/S5 填） |
|---|------|----------|----------------------|
| M1 | 匿名启动 → `auto`/显式 locale 生效 → 登录成功或失败反馈 | locale 解析/回退优先级在匿名态生效；登录失败反馈按当前语种呈现且错误码稳定 | S2：`ui-bilingual.test.tsx`（登录面双语 + 错误码映射）；S1：`locale.test.ts`（优先级矩阵）；M1 端到端 S5 |
| M2 | 登录 → Shell 导航 → overview/users/roles 可达读路径 → 有权限账号完成至少一次用户或角色写表单 → 验证/权限失败反馈 | 双语下写表单完整可完成；无权限/校验失败反馈本地化且不阻断 | S2：`ui-bilingual.test.tsx`（users 写表单 zh 提交按钮/字段）；`schema-crud.test.tsx`；M2 端到端 S5 |
| M3 | （Admin）settings → 修改 General/Branding/Localization/Appearance 一项 → Shell/登录页/公开 bootstrap 对应投影可观察生效 | `/api/branding` 扩展字段生效；配置刷新事件触发重新拉取；品牌/语种/主题投影一致 | **S3**（四类设置 + `/api/branding` 扩展后回填） |
| M4 | 制造缺失翻译 key → 记录 locale、key 与 UI/schema 路径 → 安全文本 fallback 且主流程仍可完成 | 缺失 key 可观察（事件/测试）且不渲染为空、不抛异常、不阻断操作 | S2：`ui-bilingual.test.tsx` M4 用例（事件 locale/key + 字面回退 + 列表/工具栏仍可用）；`catalog.test.ts`（去重/事件） |

## 4. 缺失翻译纪律

- 缺失 key 判定：catalog（当前语种）无条目且 en-US catalog 无条目。
- 可观察：`schema-ui:missing-translation` 事件（window，detail 含 locale/key/path）+ vitest 断言；生产不阻塞、不渲染空、不抛未处理异常。
- 安全回退链：当前语种 catalog → en-US catalog → 协议字面文本（label/title/content 等）→ key 本身。
- 任何未列入本表却出现在用户可见面的文案：只能经用户书面 `accepted-residual`（范围、缓解、责任人、复审触发）排除。

## 冻结快照

- 支持语种：`zh-CN`、`en-US`；系统默认 `auto`。
- 优先级：用户显式 → 系统默认（非 auto）→ 浏览器偏好（auto）→ `en-US`。
- 冻结依据：D-002（2026-08-09 用户书面裁决）；pageId 并集来源：`apps/api/internal/manifest/app-manifest.json` + 各模块 `manifest/fragment.json`（users/roles/settings/activity）+ `kernel/profile.go`（mvp/admin 模块集）。

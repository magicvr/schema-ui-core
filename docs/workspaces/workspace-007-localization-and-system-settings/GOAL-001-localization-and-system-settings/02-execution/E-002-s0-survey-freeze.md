---
id: E-002
doc: execution
title: S0 · 现状盘点、I-L10N-001～005 用户裁决关闭、F-V029 契约冻结
status: recorded
parent: GOAL-001-localization-and-system-settings
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

# E-002 · S0 盘点与契约冻结（2026-08-09）

## 事实

- **基线验证**：`go test ./...` 全部包通过；vitest 629/629（30 文件）全绿；`git status` 干净（最近提交 `1bb5d05` scaffold）。
- **现状盘点（I-L10N 证据）**：
  - 前端无 i18n 运行时：`src/` 无 locale/catalog 代码；`App.tsx`/`LoginPage.tsx`/`navigation.ts` 全硬编码英文；`app-manifest.ts` 已解析 `titleKey`/`labelKey` 但作为字面量使用（`navigation.ts` 直接返回 key 字符串，`App.tsx` `pageTitle = title ?? titleKey ?? pageId`）。
  - 协议兼容面：`docs/schemas/app-manifest.schema.json` 声明 `titleKey`/`labelKey`；`docs/schemas/component-registry.json`（VP-006 冻结）已为组件声明 `labelKey`/`titleKey`/`contentKey`；`node.schema.json` 的 `props` 为开放对象（仅禁 CSS 名）→ key 字段属既有兼容面。
  - 主题机制：`src/theme/theme.ts` localStorage `theme` 单通道（I-L10N-002 参照模式）。
  - 公开启动配置：`GET /api/branding`（public、no-cache）仅 `{siteTitle, logoUrl}`；配置变更事件 `X-Schema-UI-Config-Changed` + `settings.branding` 命名空间已存在（`config-events.ts`、`settings.go`）。
  - 错误 envelope：`writeError` 直写 `{error, message}`；全仓约 35 个稳定错误码（auth/account/schema/settings/resources/upload），`AuthError` 前端已按码分支。
  - 时区：全部时间 UTC 存储/返回（unix + ISO8601 UTC，`settings.go` `updatedAt`）。
  - Settings 现状：`site_settings` 表仅 `site_title`/`logo_url`；schema 页仅两字段表单。
  - Profile 模块集（`kernel/profile.go`）：mvp = core×6 + users + roles；admin = mvp + settings + activity；pageId 并集 12（core 8 + users/roles/settings/activity）。
- **门禁关闭（P-004/P-005，2026-08-09 用户逐条书面裁决）**：`I-L10N-001` 前端 key 解析；`I-L10N-002` localStorage 单通道；`I-L10N-003` 兼容扩展 `/api/branding`；`I-L10N-004` 路径 (a) 有界服务端协商；`I-L10N-005` UTC 存储 + 显示转换。详见 D-002。
- **契约冻结**：F-V029 覆盖表落 Root `attachments/F-V029-coverage-table-s0-freeze.md`（固定 UI 7 面 + 12 pageId/schemaUrl 并集 + M1～M4 + 缺失翻译纪律）。

## 产物

| 路径 | 说明 |
|------|------|
| `01-decision/D-002-s0-contract-freeze-info-gates.md` | 契约冻结决策（含五门禁关闭留痕） |
| `attachments/F-V029-coverage-table-s0-freeze.md` | F-V029 覆盖表（frozen） |
| `01-decision.md` / `02-execution.md` / `00-meta.md` / `goal-tree.md` | 索引与路线图同步（S0 → done，progress 1/6） |

## 里程碑 checkpoint

- commit：`32fcb1e`（2026-08-09，S0 契约冻结；owned paths = 本阶段全部治理文档 + S1 子目标 scaffold，显式 `git add` 无 `-A`）。

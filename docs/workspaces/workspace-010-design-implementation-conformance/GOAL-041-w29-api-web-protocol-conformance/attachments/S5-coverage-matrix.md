---
title: S5 · 运行时符合性覆盖矩阵（35 页 × 验收面）
status: active
created: 2026-09-06
updated: 2026-09-06
parent: GOAL-041-w29-api-web-protocol-conformance
version: 0.1.0
---

# S5 · 运行时符合性覆盖矩阵

> D-002 §3 的 S5 验收契约执行记录。列含义：**D-VAL** = 结构验证（`all-module-schemas-dval.test.ts`）；**Load+Negotiate** = `loadPageDocument`（含 F-001 页面级版本+能力协商）；**Render** = `RenderPage` 全分母渲染（`denominator-render.test.tsx` 35/35）；**行为触达** = 代表性页真实 schema + fixture fetcher（`representative-pages.test.tsx`）或定向行为测试。全部 35 页在 S5 均通过 D-VAL + Load + Render；「行为触达」为行为级子集，S5 后仍由 e2e 双 profile 兜底。

## 全分母矩阵

| pageId | D-VAL | Load+Negotiate | Render | 行为触达 | 备注 |
|--------|:---:|:---:|:---:|:---:|------|
| account | ✓ | ✓ | ✓ | representative + App | custom: account-session-toolbar/email-identity/mfa-manager |
| activity | ✓ | ✓ | ✓ | representative + App | custom: activity-export |
| admin-list-batch | ✓ | ✓ | ✓ | App | dev.examples |
| dashboard | ✓ | ✓ | ✓ | App | |
| data-dictionary | ✓ | ✓ | ✓ | App | |
| data-display | ✓ | ✓ | ✓ | App | dev.examples |
| data-permission | ✓ | ✓ | ✓ | representative + App | custom: data-permission-scopes |
| data-table | ✓ | ✓ | ✓ | representative + App | |
| dictionary-entries | ✓ | ✓ | ✓ | 定向（schema-dictionary-entries / row-action-bindings / App 深链） | v2.9 · route-binding + readOnly |
| digitaloffer-entitlements | ✓ | ✓ | ✓ | — | custom-only 模块；能力声明已清理 |
| digitaloffer-offers | ✓ | ✓ | ✓ | — | custom-only 模块；能力声明已清理 |
| digitaloffer-purchases | ✓ | ✓ | ✓ | — | custom-only 模块 |
| file-library | ✓ | ✓ | ✓ | App + 定向（download-behavior custom handlers） | |
| form-controls | ✓ | ✓ | ✓ | representative + App | |
| form-with-reactions | ✓ | ✓ | ✓ | representative + App | |
| form-with-upload | ✓ | ✓ | ✓ | App | |
| mail | ✓ | ✓ | ✓ | — | custom: mail-admin-tab |
| mail-outbox | ✓ | ✓ | ✓ | — | |
| my-wallet | ✓ | ✓ | ✓ | — | v2.9 · custom: wallet-ensure（探活后写） |
| notifications | ✓ | ✓ | ✓ | App + 定向（notification-center） | custom: notification-center；Host 铃铛入口 |
| overview | ✓ | ✓ | ✓ | App | dev.examples |
| recycle-bin | ✓ | ✓ | ✓ | App | |
| roles | ✓ | ✓ | ✓ | representative + App | |
| scheduled-tasks | ✓ | ✓ | ✓ | App + 定向（cron-preview） | custom: cron-preview |
| search-form-table | ✓ | ✓ | ✓ | representative + App | |
| settings | ✓ | ✓ | ✓ | representative + App | custom: password-policy-tab |
| system-monitoring | ✓ | ✓ | ✓ | App + 定向（monitoring） | custom: monitoring-auto-refresh |
| task-runs | ✓ | ✓ | ✓ | App | |
| telegram-operator | ✓ | ✓ | ✓ | — | custom-only；custom: telegram-admin-tab |
| telegram-settings | ✓ | ✓ | ✓ | — | custom-only；custom: telegram-admin-tab |
| users | ✓ | ✓ | ✓ | representative + App | custom: import-template-download |
| users-invites | ✓ | ✓ | ✓ | representative + App | custom: invite-issue-card / invite-resend-dialog |
| wallet | ✓ | ✓ | ✓ | — | v2.9 · readOnly |
| wallet-entries | ✓ | ✓ | ✓ | 定向（wallet-navigate） | v2.9 · route-binding |
| wallet-vouchers | ✓ | ✓ | ✓ | — | v2.9 |

**合计**：D-VAL 35/35 · Load+Negotiate 35/35 · Render 35/35 · 行为触达 25/35（其余 10 页为 custom-only 模块页 mail/mail-outbox 等，由 S5 全分母渲染 + e2e 双 profile 兜底）。

## 失败路径（fail-closed 验证）

| 路径 | 证据 |
|------|------|
| Manifest 加载失败 / 结构非法 / 版本不支持 / 缺能力 | `app-manifest.test.ts` / `upstream-fixtures.test.ts`（上游 fixtures 全绿） |
| 页面 D-VAL 失败 | `load-page.test.ts` PAGE_SCHEMA_INVALID |
| 页面版本不支持 | `load-page.test.ts` UNSUPPORTED_PROTOCOL_VERSION（2.10 负例） |
| 页面缺能力 | `load-page.test.ts` MISSING_REQUIRED_CAPABILITY（负例 + issues） |
| 未知标准 node | `render.test.ts` RENDER_UNKNOWN_NODE_TYPE |
| 未知 custom 控件 | `render.test.tsx` C-010 明显占位 + console.error |
| 未知 custom handler | `download-behavior.test.tsx` CUSTOM_HANDLER_NOT_FOUND |
| pageId 不匹配 / 404 / 5xx / 非 JSON | `load-page.test.ts` PAGE_ID_MISMATCH / PAGE_NOT_FOUND / PAGE_LOAD_FAILED / PAGE_PARSE_FAILED |

## HTTP Manifest 快照（profile 验收矩阵）

由 `apps/api/internal/composition/s5_manifest_snapshot_test.go` 生成（真实组合根 newMux → httptest GET `/.well-known/schema-ui/app-manifest.json`）：

| 组合 | pages | 快照 | 说明 |
|------|:---:|------|------|
| mvp | 6 | `S5-manifest-snapshots/mvp/app-manifest.json` | 默认最小 |
| admin | 22 | `S5-manifest-snapshots/admin/app-manifest.json` | 生产默认 |
| demo | 14 | `S5-manifest-snapshots/demo/app-manifest.json` | 演示 |
| admin + biz.digital-offer | 25 | `S5-manifest-snapshots/admin+digitaloffer/app-manifest.json` | 显式 custom |
| admin + channel.telegram | 24 | `S5-manifest-snapshots/admin+telegram/app-manifest.json` | 显式 custom（真实 TelegramRuntime 装配） |

全部快照：`protocolVersion 2.7`、`requiredCapabilities = [app.manifest, app.navigation]`（envelope）、页面数为冻结投影；证明 served manifest 是真实 profile 投影而非 35 页静态宇宙。

---
title: S1 · API/Web 页面与控件完整目录
status: active
created: 2026-09-06
updated: 2026-09-06
parent: GOAL-041-w29-api-web-protocol-conformance
version: 0.1.0
---

# S1 · API/Web 页面与控件完整目录

> 冻结当前源码快照的完整静态页面/控件分母，并以真实 Manifest 聚合函数核对 profile 投影；不替代 S5 的启动服务、权限、真实数据与浏览器验证。完整字段见 `S1-api-web-page-control-inventory.json`。

## 1. 分母与生成方法

- 被审 commit：`ce66abab0897b644c695b8365f8324598bfe2bc8`。
- API fragments：**17**；Manifest pages：**35**；schema documents：**35**。
- 重复 pageId：0；缺失 schema：0；orphan schema：0。
- 组合链：profile/config → Provider contributions → `RegisterSchemas` → `ForModulesWithFragments` → `RegisterManifest`。
- 临时 Go test 调用真实 `ForModulesWithFragments` 并按 profile 过滤 fragment；测试文件已删除，输出写入 E-002。

## 2. Runtime Manifest profile 投影

| profile | pages | page IDs |
|---|---:|---|
| `mvp` | 6 | `account`, `dashboard`, `notifications`, `roles`, `users`, `users-invites` |
| `admin` | 22 | `account`, `activity`, `dashboard`, `data-dictionary`, `data-permission`, `dictionary-entries`, `file-library`, `mail`, `mail-outbox`, `my-wallet`, `notifications`, `recycle-bin`, `roles`, `scheduled-tasks`, `settings`, `system-monitoring`, `task-runs`, `users`, `users-invites`, `wallet`, `wallet-entries`, `wallet-vouchers` |
| `demo` | 14 | `account`, `admin-list-batch`, `dashboard`, `data-display`, `data-table`, `form-controls`, `form-with-reactions`, `form-with-upload`, `notifications`, `overview`, `roles`, `search-form-table`, `users`, `users-invites` |

- 全部已注册 fragment 宇宙：**35 页**。
- `channel.telegram`、`admin.digital-offer` 等不在默认 MVP/Admin/Demo profile 内，只能由显式 custom 模块配置进入。

## 3. 完整页面目录（35）

| module / profile | page / route | schema / protocol | nodes | controls | actions | custom / v2.9 | nav |
|---|---|---|---|---|---|---|---|
| `admin.account` / mvp/admin/demo | `account`<br>`/account` | `apps/api/modules/account/schema/account.json:3`<br>2.7 | custom:3, form:3, section:4, table:1, tabs:1 | input:1, password:2, select:1, switch:1, upload:1 | saveProfile:request, uploadAvatar:upload, changePassword:request, revokeSession:request, saveNotificationSettings:request | body account-session-toolbar,email-identity,mfa-manager | user:account |
| `admin.activity` / admin | `activity`<br>`/activity` | `apps/api/modules/activity/schema/activity.json:3`<br>2.7 | custom:1, form:2, recordView:1, section:1, table:1, text:1 | datePicker:2, input:9, textarea:1 | openDetail:modal | body activity-export | sidebar:activity |
| `dev.examples` / demo | `admin-list-batch`<br>`/admin-list-batch` | `apps/api/modules/dev/examples/schema/admin-list-batch.json:3`<br>2.7 | section:1, table:1, text:1 | — | deleteUsers:request | — | sidebar:admin-list-batch |
| `admin.dashboard` / mvp/admin/demo | `dashboard`<br>`/dashboard` | `apps/api/modules/dashboard/schema/dashboard.json:3`<br>2.7 | grid:1, section:1, statCard:2, text:1 | — | — | — | sidebar:dashboard |
| `admin.data-dictionary` / admin | `data-dictionary`<br>`/data-dictionary` | `apps/api/modules/datadictionary/schema/data-dictionary.json:3`<br>2.7 | form:3, section:1, table:1 | input:4, inputNumber:2, switch:2, textarea:2 | createType:request, updateType:request, deleteType:request, openCreate:modal, openEdit:modal, openEntries:navigate | — | sidebar:data-dictionary |
| `dev.examples` / demo | `data-display`<br>`/data-display` | `apps/api/modules/dev/examples/schema/data-display.json:3`<br>2.7 | chart:1, grid:1, section:1, statCard:2, text:1 | — | — | — | sidebar:data-display |
| `admin.data-permission` / admin | `data-permission`<br>`/data-permission` | `apps/api/modules/datapermission/schema/data-permission.json:3`<br>2.7 | custom:1, form:1, section:1, table:1 | input:2, select:1, switch:1 | registerPolicy:request, updateScopes:request, openRegister:modal | body data-permission-scopes | sidebar:data-permission |
| `dev.examples` / demo | `data-table`<br>`/data-table` | `apps/api/modules/dev/examples/schema/data-table.json:3`<br>2.7 | section:1, table:1, text:1 | — | — | — | sidebar:data-table |
| `admin.data-dictionary` / admin | `dictionary-entries`<br>`/dictionary-entries/{dictKey}` | `apps/api/modules/datadictionary/schema/dictionary-entries.json:3`<br>2.9 | form:3, section:1, table:1 | input:8, inputNumber:2, select:2, switch:2, textarea:2 | createEntry:request, updateEntry:request, deleteEntry:request, openCreate:modal, openEdit:modal | route-bind:1, readOnly:4 | inner/no direct nav |
| `admin.digital-offer` / custom-only | `digitaloffer-entitlements`<br>`/digitaloffer-entitlements` | `apps/api/modules/digitaloffer/schema/digitaloffer-entitlements.json:3`<br>2.9 | form:1, section:1, table:1 | input:2, select:1 | voidEntitlement:request | — | sidebar:digitaloffer-entitlements |
| `admin.digital-offer` / custom-only | `digitaloffer-offers`<br>`/digitaloffer-offers` | `apps/api/modules/digitaloffer/schema/digitaloffer-offers.json:3`<br>2.9 | form:4, section:1, table:1 | input:4, inputNumber:5, select:3, switch:1, textarea:2 | createOffer:request, updateOffer:request, changeStatus:request, openCreate:modal, openEdit:modal, openStatus:modal | — | sidebar:digitaloffer-offers |
| `admin.digital-offer` / custom-only | `digitaloffer-purchases`<br>`/digitaloffer-purchases` | `apps/api/modules/digitaloffer/schema/digitaloffer-purchases.json:3`<br>2.9 | form:1, section:1, table:1 | input:3 | — | — | sidebar:digitaloffer-purchases |
| `admin.file-library` / admin | `file-library`<br>`/file-library` | `apps/api/modules/filelibrary/schema/file-library.json:3`<br>2.7 | form:2, section:1, table:1 | input:1, upload:1 | uploadFile:upload, submitUpload:request, openUpload:modal, downloadFile:custom, previewFile:custom, copyLinkFile:custom, deleteFile:request | — | sidebar:file-library |
| `dev.examples` / demo | `form-controls`<br>`/form-controls` | `apps/api/modules/dev/examples/schema/form-controls.json:3`<br>2.7 | form:1, section:1, text:1 | checkboxGroup:1, datePicker:1, dateRangePicker:1, input:1, inputNumber:1, radio:1, select:1, switch:1, textarea:1 | — | — | sidebar:form-controls |
| `dev.examples` / demo | `form-with-reactions`<br>`/form-with-reactions` | `apps/api/modules/dev/examples/schema/form-with-reactions.json:3`<br>2.7 | form:1, section:1, text:1 | input:1, select:1, switch:1, textarea:1 | — | — | sidebar:form-with-reactions |
| `dev.examples` / demo | `form-with-upload`<br>`/form-with-upload` | `apps/api/modules/dev/examples/schema/form-with-upload.json:3`<br>2.7 | form:1, section:1, text:1 | input:1, upload:1 | uploadAttachment:upload | — | sidebar:form-with-upload |
| `admin.settings` / admin | `mail`<br>`/mail` | `apps/api/modules/settings/schema/mail.json:3`<br>2.7 | custom:1, section:1 | — | — | body mail-admin-tab | sidebar:mail |
| `admin.settings` / admin | `mail-outbox`<br>`/mail-outbox` | `apps/api/modules/settings/schema/mail-outbox.json:3`<br>2.7 | form:1, recordView:1, section:1, table:1 | input:1, select:2 | — | — | sidebar:mail-outbox |
| `admin.wallet` / admin | `my-wallet`<br>`/my-wallet` | `apps/api/modules/wallet/schema/my-wallet.json:3`<br>2.9 | custom:1, form:1, grid:1, section:1, statCard:3, table:1, text:1 | input:1 | openWallet:request, redeemVoucher:request, openRedeem:modal | body wallet-ensure | user:my-wallet |
| `admin.notifications` / mvp/admin/demo | `notifications`<br>`/notifications` | `apps/api/modules/notifications/schema/notifications.json:3`<br>2.7 | custom:1, form:1, section:1 | input:1, select:1 | markAllRead:request | body notification-center | inner/no direct nav |
| `dev.examples` / demo | `overview`<br>`/overview` | `apps/api/modules/dev/examples/schema/overview.json:3`<br>2.7 | section:1, text:1 | — | — | — | top:overview |
| `admin.recycle-bin` / admin | `recycle-bin`<br>`/recycle-bin` | `apps/api/modules/recyclebin/schema/recycle-bin.json:3`<br>2.7 | form:1, section:1, table:1 | input:1, select:1 | restore:request, purge:request, purgeAll:request | — | sidebar:recycle-bin |
| `admin.roles` / mvp/admin/demo | `roles`<br>`/roles` | `apps/api/modules/roles/schema/roles.json:3`<br>2.7 | form:3, recordView:1, section:1, table:1 | checkboxGroup:4, input:4, select:1 | createRole:request, updateRole:request, deleteRole:request, exportRoles:custom, openCreate:modal, openEdit:modal | — | sidebar:roles |
| `admin.scheduled-tasks` / admin | `scheduled-tasks`<br>`/scheduled-tasks` | `apps/api/modules/scheduledtasks/schema/scheduled-tasks.json:3`<br>2.7 | form:3, section:1, table:1 | input:6, select:3, switch:2, textarea:2 | createTask:request, updateTask:request, deleteTask:request, runTask:request, openCreate:modal, openEdit:modal, openRuns:navigate | after cron-preview | sidebar:scheduled-tasks |
| `dev.examples` / demo | `search-form-table`<br>`/search-form-table` | `apps/api/modules/dev/examples/schema/search-form-table.json:3`<br>2.7 | form:1, section:1, table:1, text:1 | input:1 | — | — | sidebar:search-form-table |
| `admin.settings` / admin | `settings`<br>`/settings` | `apps/api/modules/settings/schema/settings.json:3`<br>2.7 | actionButton:1, custom:1, form:6, section:8, tabs:1, text:1 | input:5, inputNumber:1, select:4, upload:4 | updateGeneral:request, updateBranding:request, uploadBrandingLogo:upload, uploadBrandingFavicon:upload, updateLocalization:request, updateAppearance:request, updateAudit:request, resetSettings:request, updateCaptcha:request | body password-policy-tab | user:settings |
| `admin.system-monitoring` / admin | `system-monitoring`<br>`/system-monitoring` | `apps/api/modules/systemmonitoring/schema/system-monitoring.json:3`<br>2.7 | custom:1, grid:1, section:1, statCard:6, table:1, text:1 | — | — | body monitoring-auto-refresh | sidebar:system-monitoring |
| `admin.scheduled-tasks` / admin | `task-runs`<br>`/task-runs` | `apps/api/modules/scheduledtasks/schema/task-runs.json:3`<br>2.7 | form:1, section:1, table:1 | input:1, select:1 | — | — | inner/no direct nav |
| `channel.telegram` / custom-only | `telegram-operator`<br>`/telegram-settings/operator` | `apps/api/modules/channel/telegram/schema/telegram-operator.json:3`<br>2.7 | custom:1 | — | — | body telegram-admin-tab | inner/no direct nav |
| `channel.telegram` / custom-only | `telegram-settings`<br>`/telegram-settings` | `apps/api/modules/channel/telegram/schema/telegram-settings.json:3`<br>2.7 | actionButton:1, custom:1, section:2, text:1 | — | openTelegramOperator:navigate | body telegram-admin-tab | sidebar:telegram-settings |
| `admin.users` / mvp/admin/demo | `users`<br>`/users` | `apps/api/modules/users/schema/users.json:3`<br>2.7 | form:6, recordView:1, section:1, table:1 | checkboxGroup:2, input:4, password:2, select:2, upload:1 | openInvites:navigate, createUser:request, updateUser:request, updateUserRoles:request, changeUserPassword:request, deleteUser:request, enableUser:request, disableUser:request, unlockUser:request, resetUserMfa:request, uploadCsv:upload, exportUsers:custom, submitImport:request, openImport:modal, openCreate:modal, openEdit:modal, openRoles:modal, openPassword:modal | after import-template-download | sidebar:users |
| `admin.users` / mvp/admin/demo | `users-invites`<br>`/users-invites` | `apps/api/modules/users/schema/users-invites.json:3`<br>2.7 | custom:2, form:1, recordView:1, section:1, table:1 | input:1, select:1 | revokeInvite:request, resendInvite:modal | body invite-issue-card, action-content invite-resend-dialog | inner/no direct nav |
| `admin.wallet` / admin | `wallet`<br>`/wallet` | `apps/api/modules/wallet/schema/wallet.json:3`<br>2.9 | form:7, section:1, table:1 | input:7, inputNumber:5, select:3, textarea:4 | createAccount:request, updateStatus:request, adjust:request, freeze:request, unfreeze:request, deductFrozen:request, runReconcile:request, openCreate:modal, openAdjust:modal, openFreeze:modal, openUnfreeze:modal, openDeductFrozen:modal, openStatus:modal, openEntries:navigate | readOnly:1 | sidebar:wallet |
| `admin.wallet` / admin | `wallet-entries`<br>`/wallet-entries/{id}` | `apps/api/modules/wallet/schema/wallet-entries.json:3`<br>2.9 | form:1, section:1, table:1 | input:1, select:1 | — | route-bind:1 | inner/no direct nav |
| `admin.wallet` / admin | `wallet-vouchers`<br>`/wallet-vouchers` | `apps/api/modules/wallet/schema/wallet-vouchers.json:3`<br>2.9 | form:2, section:1, table:1 | datePicker:1, input:2, inputNumber:2, select:1 | generateBatch:request, openGenerate:modal, voidVoucher:request | — | sidebar:wallet-vouchers |

## 4. Web Renderer 与扩展面

- Renderer node types（11）：`form`, `section`, `table`, `grid`, `tabs`, `text`, `recordView`, `actionButton`, `statCard`, `chart`, `custom`。
- Form controls（14）：`input`, `select`, `inputNumber`, `datePicker`, `dateRangePicker`, `textarea`, `switch`, `checkbox`, `radio`, `cascader`, `checkboxGroup`, `richText`, `password`, `upload`。
- 未知标准 node：`RENDER_UNKNOWN_NODE_TYPE`，fail closed。
- 未知 custom component：当前显示安全 inline fallback；是否符合上游 failure contract 待 S2。
- custom registrations：**15**；其中 12 个用于页面 `body` custom node、1 个用于 action `content` custom node、2 个用于 `afterComponent`；当前使用键均有注册。

### 4.1 custom component registry

| key | 使用形态 | 注册证据 |
|---|---|---|
| `account-session-toolbar` | body custom node ×1 | `apps/web/src/components/account-session-toolbar.tsx:84` |
| `activity-export` | body custom node ×1 | `apps/web/src/components/activity-export.tsx:84` |
| `cron-preview` | afterComponent ×2 | `apps/web/src/components/cron-preview.tsx:134` |
| `data-permission-scopes` | body custom node ×1 | `apps/web/src/components/data-permission-scopes.tsx:282` |
| `email-identity` | body custom node ×1 | `apps/web/src/components/email-identity.tsx:195` |
| `import-template-download` | afterComponent ×1 | `apps/web/src/components/import-template-download.tsx:70` |
| `invite-issue-card` | body custom node ×1 | `apps/web/src/components/invite-issue-card.tsx:248` |
| `invite-resend-dialog` | action content custom ×1 | `apps/web/src/components/invite-resend-dialog.tsx:133` |
| `mail-admin-tab` | body custom node ×1 | `apps/web/src/components/mail-admin-tab.tsx:350` |
| `mfa-manager` | body custom node ×1 | `apps/web/src/components/mfa-manager.tsx:379` |
| `monitoring-auto-refresh` | body custom node ×1 | `apps/web/src/components/monitoring-auto-refresh.tsx:62` |
| `notification-center` | body custom node ×1 | `apps/web/src/components/notification-center.tsx:250` |
| `password-policy-tab` | body custom node ×1 | `apps/web/src/components/password-policy-tab.tsx:137` |
| `telegram-admin-tab` | body custom node ×2 | `apps/web/src/components/telegram-admin-tab.tsx:1089` |
| `wallet-ensure` | body custom node ×1 | `apps/web/src/components/wallet-ensure.tsx:112` |

### 4.2 custom action handlers

| handler | API target | schema uses |
|---|---|---:|
| `export.users` | `/api/export/users` | 1 |
| `export.roles` | `/api/export/roles` | 1 |
| `library.download` | `/api/library/files/{id}/download` | 1 |
| `library.preview` | `/api/library/files/{id}/download` | 1 |
| `library.copyLink` | `/api/library/files/{id}/download` | 1 |

## 5. Host/App 壳层单列

| surface | 证据 | S1 边界 |
|---|---|---|
| login | `apps/web/src/app/LoginPage.tsx:60-91` | Host/App 候选；不自动判为协议绕过 |
| auth/session gate | `apps/web/src/app/AuthGate.tsx:48-78` | Host/App 候选；不自动判为协议绕过 |
| branding | `apps/web/src/app/branding.ts; apps/web/src/app/App.tsx:786` | Host/App 候选；不自动判为协议绕过 |
| host boot/recovery | `apps/web/src/host/boot.ts:269-320` | Host/App 候选；不自动判为协议绕过 |
| host failure | `apps/web/src/host/failure.ts:43-180; apps/web/src/app/HostFailureScreen.tsx:42-56` | Host/App 候选；不自动判为协议绕过 |
| manifest failure | `apps/web/src/app/ManifestFailure.tsx:6` | Host/App 候选；不自动判为协议绕过 |
| shell/navigation/user menu | `apps/web/src/app/App.tsx:291-446` | Host/App 候选；不自动判为协议绕过 |

## 6. 测试与触达现状

- representative pages：**10/35**：`data-table`, `search-form-table`, `form-controls`, `form-with-reactions`, `users`, `users-invites`, `roles`, `settings`, `activity`, `data-permission`。
- representative 测试直接读取 API schema 并注入 fixture fetcher，不是实际 API handler/HTTP Manifest。
- `mvp-dogfood`：13 页；`admin-dogfood`：26 页；均为协议 2.7，且与当前默认 profile 投影不同。
- S5 仍需真实服务、权限、数据和浏览器链路；S1 只冻结目录。

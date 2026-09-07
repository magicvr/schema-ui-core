---
doc_type: goal-attachment
id: r1-navigation-profile-slot-matrix
parent: GOAL-001-nav-group-collapsible
status: recorded
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
---

# R1 · 当前导航 / Profile / slot 盘点

> 这是 R1 的代码盘点记录，不是运行时 e2e 验收。Profile 默认成员来自 `apps/api/kernel/profile.go`；NodeID、PageID、slot 与可选模块来自各 Provider/Manifest fragment。运行时最终 Manifest、权限过滤与直接 URL 行为仍在 R3/R4 验证。

## 1. 当前 sidebar 分组分母

| group key | frozen order | current NodeID | page | module | baseline profile |
|---|---:|---|---|---|---|
| identity-access | 10 | `menu_users` | `users` | `admin.users` | mvp/admin/demo/custom |
| identity-access | 10 | `menu_roles` | `roles` | `admin.roles` | mvp/admin/demo/custom |
| identity-access | 10 | `menu_data_permission` | `data-permission` | `admin.data-permission` | admin/custom |
| content-data | 20 | `menu_files` | `file-library` | `admin.file-library` | admin/custom |
| content-data | 20 | `menu_dictionary` | `data-dictionary` | `admin.data-dictionary` | admin/custom |
| operations | 30 | `menu_activity` | `activity` | `admin.activity` | admin/custom |
| operations | 30 | `menu_monitoring` | `system-monitoring` | `admin.system-monitoring` | admin/custom |
| operations | 30 | `menu_scheduled_tasks` | `scheduled-tasks` | `admin.scheduled-tasks` | admin/custom |
| operations | 30 | `menu_recycle_bin` | `recycle-bin` | `admin.recycle-bin` | admin/custom |
| communications | 40 | `menu_mail` | `mail` | `admin.settings` | admin/custom |
| communications | 40 | `menu_mail_outbox` | `mail-outbox` | `admin.settings` | admin/custom |
| communications | 40 | `menu_telegram` | `telegram-settings` | `channel.telegram` | custom when enabled |
| commerce | 50 | `menu_wallet` | `wallet` | `admin.wallet` | admin/custom |
| commerce | 50 | `menu_wallet_vouchers` | `wallet-vouchers` | `admin.wallet` | admin/custom |
| commerce | 50 | `menu_digitaloffer_offers` | `digitaloffer-offers` | `biz.digital-offer` | custom when enabled |
| commerce | 50 | `menu_digitaloffer_entitlements` | `digitaloffer-entitlements` | `biz.digital-offer` | custom when enabled |
| commerce | 50 | `menu_digitaloffer_purchases` | `digitaloffer-purchases` | `biz.digital-offer` | custom when enabled |

## 2. Intentional non-grouped or non-sidebar surfaces

| surface | current slot / behavior | verification boundary |
|---|---|---|
| `menu_dashboard` | sidebar top-level singleton | position, active state, direct URL in R4 |
| `dev.examples` | existing `Examples` Manifest group; demo-only | preserve group and run demo regression in R4 |
| `menu_settings` | user slot from `admin.settings` | remain user slot; do not move to sidebar |
| `menu_account` | user slot from `admin.account` | remain user slot; do not move to sidebar |
| `menu_wallet_self` | user slot from `admin.wallet` | remain user slot; do not move to sidebar |
| notifications | shell notification bell / no sidebar NodeID | bell and notification route regression; not grouped |
| `admin.data-transfer`, `admin.login-captcha`, `admin.mfa` | no independent navigation contribution | no sidebar omission: no current registered sidebar node |

## 3. Profile coverage

| profile | current relevant modules | expected navigation coverage |
|---|---|---|
| `mvp` | dashboard, users, roles, account, notifications | Dashboard singleton + partial identity-access; account remains user slot; notifications remain bell |
| `admin` | mvp surface + settings, activity, data-transfer, file-library, data-dictionary, system-monitoring, scheduled-tasks, login-captcha, recycle-bin, data-permission, mfa, wallet | Dashboard + identity-access + content-data + operations + communications(mail/outbox) + commerce(wallet/vouchers); settings/account/my-wallet remain user slot |
| `demo` | mvp surface + `dev.examples` | mvp navigation + top `Overview` + existing `Examples` group |
| `custom` | explicit compiled module list | R4 must exercise a custom set containing `channel.telegram` and `biz.digital-offer`; a demo/custom combination must verify cross-module `commerce`/`communications` aggregation and no loss |

## 4. Evidence and limits

- Static sources: `apps/api/kernel/profile.go`; `apps/api/kernel/contribution.go`; the module `provider.go` files; module `manifest/fragment.json` files; `apps/api/internal/manifest/manifest.go`; `apps/web/src/app/navigation.ts`; `apps/web/src/app/App.tsx`.
- Baseline API validation: `go test ./... -count=1` passed on 2026-09-07.
- This matrix does not claim that group metadata, grouped Manifest output, collapsible UI, keyboard behavior, session state, or direct URL auto-expansion is implemented. Those are R2/R3/R4 facts to be produced later.

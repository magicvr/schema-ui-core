---
doc_type: goal-attachment
id: r4-navigation-route-profile-matrix
parent: GOAL-001-nav-group-collapsible
status: recorded
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
---

# R4 · 当前 sidebar 全量迁移与 Profile/route 矩阵

> 本附件把当前已注册 sidebar 分母、其页面 route、深链父级和运行时 Profile 组合固定为 R4 验收输入。代码事实来自各模块 Manifest fragment；运行时证据来自 `apps/api/internal/composition/nav_group_r4_test.go`，Web projection/深链证据来自 `apps/web/src/app/nav-groups-r4.test.ts`。

## 1. Sidebar NodeID → group → route

| NodeID | group | pageRef | direct route | deep-link parent coverage |
|---|---|---|---|---|
| `menu_dashboard` | top-level singleton | `dashboard` | `/dashboard` | singleton remains outside groups |
| `menu_users` | identity-access | `users` | `/users` | `users-invites` → users |
| `menu_roles` | identity-access | `roles` | `/roles` | — |
| `menu_data_permission` | identity-access | `data-permission` | `/data-permission` | — |
| `menu_files` | content-data | `file-library` | `/file-library` | — |
| `menu_dictionary` | content-data | `data-dictionary` | `/data-dictionary` | `dictionary-entries` → data-dictionary |
| `menu_activity` | operations | `activity` | `/activity` | — |
| `menu_monitoring` | operations | `system-monitoring` | `/system-monitoring` | — |
| `menu_scheduled_tasks` | operations | `scheduled-tasks` | `/scheduled-tasks` | `task-runs` → scheduled-tasks |
| `menu_recycle_bin` | operations | `recycle-bin` | `/recycle-bin` | — |
| `menu_mail` | communications | `mail` | `/mail` | — |
| `menu_mail_outbox` | communications | `mail-outbox` | `/mail-outbox` | — |
| `menu_telegram` | communications | `telegram-settings` | `/telegram-settings` | `telegram-operator` → telegram-settings |
| `menu_wallet` | commerce | `wallet` | `/wallet` | `wallet-entries/{id}` → wallet |
| `menu_wallet_vouchers` | commerce | `wallet-vouchers` | `/wallet-vouchers` | — |
| `menu_digitaloffer_offers` | commerce | `digitaloffer-offers` | `/digitaloffer-offers` | — |
| `menu_digitaloffer_entitlements` | commerce | `digitaloffer-entitlements` | `/digitaloffer-entitlements` | — |
| `menu_digitaloffer_purchases` | commerce | `digitaloffer-purchases` | `/digitaloffer-purchases` | — |

## 2. Slot boundaries

| slot/surface | expected behavior |
|---|---|
| sidebar | Dashboard first; five structured groups by GroupOrder 10/20/30/40/50; current registered nodes above; no missing/duplicate pageRef |
| top | only profile-specific top contributions (demo `overview`); never normalized into sidebar groups |
| user | `account`, `my-wallet`, `settings` remain user-slot links when their modules are enabled; never normalized into sidebar groups |
| notification | shell NotificationBell; `menu_notifications` system-data row is not a Manifest sidebar link |
| `dev.examples` | demo-only authored `Examples` NavGroup remains independent from Admin five groups |

## 3. Runtime Profile matrix

| profile/combo | identity-access | content-data | operations | communications | commerce | top | user | Examples |
|---|---|---|---|---|---|---|---|---|
| `mvp` | users, roles | — | — | — | — | — | account | — |
| `admin` | users, roles, data-permission | file-library, data-dictionary | activity, system-monitoring, scheduled-tasks, recycle-bin | mail, mail-outbox | wallet, wallet-vouchers | — | account, my-wallet, settings | — |
| `demo` | users, roles | — | — | — | — | overview | account | present |
| custom = admin + `channel.telegram` | admin set | admin set | admin set | mail, mail-outbox, telegram-settings | wallet, wallet-vouchers | — | account, my-wallet, settings | — |
| custom = admin + `channel.telegram` + `biz.digital-offer` | admin set | admin set | admin set | mail, mail-outbox, telegram-settings | wallet, wallet-vouchers, digitaloffer-offers, digitaloffer-entitlements, digitaloffer-purchases | — | account, my-wallet, settings | — |

## 4. Verification evidence

- API runtime Manifest/Profile matrix: `go test ./internal/composition -run TestR4NavigationGroupProfileMatrix -count=1` passed.
- Web direct/deep route projection matrix: `nav-groups-r4.test.ts` 2/2 passed; R3 interaction suite `nav-groups.test.tsx` 3/3 and `navigation.test.ts` 7/7 passed.
- Web type check: direct `tsc -b` passed.
- The matrix does not by itself claim VP/R5 close-out; R5 still requires full regression, playbook update, evidence synthesis and close-out audit.

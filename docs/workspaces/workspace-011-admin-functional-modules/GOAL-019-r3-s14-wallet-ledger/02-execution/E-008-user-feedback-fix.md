---
id: E-008
goal: GOAL-019-r3-s14-wallet-ledger
title: 用户反馈修复：PAGE_SCHEMA_INVALID + 菜单图标 + 导航排序
date: 2026-08-16
status: recorded
parent: GOAL-019-r3-s14-wallet-ledger
created: 2026-08-16
updated: 2026-08-16
version: 1.0.0
---

# E-008 · 用户反馈修复（2026-08-16）

## 事实

- **PAGE_SCHEMA_INVALID（根因）**：wallet.json 违反 pin 的 page/action/node schema——(a) request action 携带 `requestMapping`（action.schema.json 不允许，additionalProperties: false；行级 `{id}` 由 renderer 按 row 自动填充，scheduled-tasks 先例）；(b) `permissionCascade.keys` 含非枚举值 `read`/`adjust`（仅允许 edit/delete）；(c) 行操作 `permissionIntent: "adjust"` 非枚举（仅 edit/delete）；(d) 对账按钮 `permissionIntent: "read"` 非枚举。运行时 load-page 强制结构验证（D-VAL）→ fail closed。
- **修复**：移除 4 处 requestMapping；表格 permissionCascade keys=["edit"]，permissions.edit = wallet.adjust（ADR-0023 意图仅 edit/delete，UI 显隐按最敏感写键收紧；API 仍按端点三键强制——A-007 已核对）；对账按钮移除 permissionIntent（API 为 wallet.read，按钮始终显示）；行操作 intent 统一 edit。wallet-entries.json 无此问题（验证通过）。
- **图标**：iconRegistry（apps/web/src/app/App.tsx）无 `wallet` 键 → fragment.json `"icon": "wallet"` 无渲染映射。注册 `wallet: Wallet`（lucide-react）。
- **导航排序（用户指令）**：menu_wallet 从尾部移至 **menu_roles 之后（角色下面）、menu_account 之前**（即 Roles → Wallet → Account → Activity；Account 在 user slot，Activity 为「操作日志」）。同步 DefaultNavigationOrder + navigation_order_test 快照 + composition_test sidebar want + testsupport Order + admin fixture（Admin 组 roles 之后）+ SHA 重钉。
- **验证**：wallet.json/wallet-entries.json 通过 validatePageDocument（D-VAL）；go 全量 + web 全量回归全绿（见 pwsh-394/395）。

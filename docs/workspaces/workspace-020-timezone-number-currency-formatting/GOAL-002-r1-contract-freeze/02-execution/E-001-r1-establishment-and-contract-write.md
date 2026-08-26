---
id: GOAL-002-r1-contract-freeze
doc: execution-entry
record_id: E-001
status: recorded
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
---

# E-001 · 立项与合同正文落盘

## 2026-08-26

### 已发生事实

1. GOAL-002 五件套 + 三个 ledger 目录建立（`01` 原语；模板源 `docs/templates/goal-folder/`）。
2. 合同正文 `01-decision/D-001-r1-contract-freeze.md` 落盘（§0～§6：效力 / 范围与落点 / 时区来源 / 数字货币落点 / 设置归属与字段 / 越界 / 消费指引）。
3. 代码基现状核对（C2）完成：
   - `SiteSettings.siteTimezone` / `defaultLocale` 存在（`apps/api/internal/handler/settings.go`；`/api/branding` 公开投影）→ §2 L3 / §4.1「已有」断言成立。
   - 金额 int64 JSON（`apps/api/internal/handler/wallet.go`）、时间 RFC3339 毫秒 UTC（`apps/api/internal/handler/rfc3339.go`）→ §3.3 机器格式断言成立。
   - `apps/web/src/i18n/format.ts`（Intl.* fail-safe）、`runtime.tsx`（`schema-ui:locale` 单通道）→ §2 L1 通道与 §3.1 展示既有基础成立。
   - Localization tab 已含 locale + timezone 字段（`apps/web/src/app/startup-config.test.tsx` 快测）→ §4.1 归属落位成立。
4. 越界守卫：本轮全部改动限于 `docs/workspaces/workspace-020-timezone-number-currency-formatting/**`；未改 DDL / 迁移台账 / Profile 默认集 / 模块矩阵 / Manifest / `docs/contracts/`。
5. 状态留痕：C1 done、C2 done（meta 检查点表）；合同正文状态 = accepted（lead 提案 + 前置门禁已用户裁决；正文全文待用户审阅，审阅意见进入整改环或 A-001 自审）。

### 证据

| 主张 | 路径 / 命令 / commit |
|------|----------------------|
| 五件套与 ledger 目录齐全 | `docs/workspaces/workspace-020-timezone-number-currency-formatting/GOAL-002-r1-contract-freeze/`（00-meta / 01-decision / 02-execution / 03-audit / attachments + 三个 ledger 目录） |
| 合同正文落盘 | `GOAL-002-r1-contract-freeze/01-decision/D-001-r1-contract-freeze.md` |
| 现状核对（C2）证据 | 见 Root `GOAL-001-.../02-execution/E-002` 证据表（settings.go / wallet.go / rfc3339.go / format.ts / runtime.tsx / startup-config.test.tsx） |
| 改动范围 | `git status` / 提交记录（仅 workspace-020 文档路径） |
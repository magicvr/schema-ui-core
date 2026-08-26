---
id: GOAL-003-r2-timezone-semantics
doc: execution-entry
record_id: E-002
status: recorded
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
---

# E-002 · C1 生效时区解析器实施

## 2026-08-26

### 已发生事实

1. **新模块** `apps/web/src/i18n/timezone.ts`（与 `locale.ts` 同构）：`TIMEZONE_STORAGE_KEY`、`isValidIanaTimeZone`（Intl 探针校验）、`normalizeTimezonePreference`、`readStoredTimezone` / `writeStoredTimezone`（单通道；`auto` = 移除 key）、`detectBrowserTimezone`（会话探测）、`resolveEffectiveTimezone`（L1 → L2 → L3 → L4）。
2. **快测** `apps/web/src/i18n/timezone.test.ts`（jsdom）：15 用例覆盖 L1 覆盖优先 / L1 auto·无效值穿透 / L2 探测 / L2 空·无效跳过 / L3 站点默认 / L3 auto·空跳过 / L4 兜底 / 敌意输入不抛错 / 探测合法性 / 存储 round-trip。
3. **验证**：`npx vitest run src/i18n/timezone.test.ts` → **15/15 pass**；`npx vitest run src/i18n` → **79/79 pass**（无回归，含 runtime/ui-bilingual/s5-denominator）。
4. 越界守卫：本轮实施仅新增 web i18n 纯函数模块 + 快测；未改 API / DB / DDL / Profile 默认集 / `docs/contracts/`。

### 证据

| 主张 | 路径 / 命令 / commit |
|------|----------------------|
| 解析器实现（L1～L4） | `apps/web/src/i18n/timezone.ts` |
| 快测 15/15 | `npx vitest run src/i18n/timezone.test.ts`（exit 0） |
| i18n 全量 79/79 无回归 | `npx vitest run src/i18n`（exit 0） |
| 合同条款 | `GOAL-002-r1-contract-freeze/01-decision/D-001` §2（L1～L4）/ §4.2（单通道） |
| 改动范围 | `git status`（仅 `apps/web/src/i18n/timezone.ts` + `timezone.test.ts` + 本区 docs） |
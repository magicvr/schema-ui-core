---
id: GOAL-005-r4-evidence-closeout
doc: execution-entry
record_id: E-002
status: recorded
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
---

# E-002 · C1/C2/C3 证据矩阵与核账项处置

## 2026-08-26

### 已发生事实

1. **C1 证据矩阵落盘**：`attachments/r4-evidence-matrix.md`——合同 §1～§6 逐条证据映射（含 commit/测试文件/用例数）；双 locale 范例表（zh-CN/en-US 各时区展示、货币展示、输入解析、round-trip 场景）；全量回归基线（Go 全绿 + web 88 files/1181 tests）。
2. **C2 越界核账**：§5 清单 + Root 边界逐项核对（汇率/计费、RT-T03、Profile 默认集、模块矩阵、Manifest、`docs/contracts/`、热加载、钱包演示面重开）全部成立；API 机器合同不变量复证。
3. **C3 核账项处置**：
   - GOAL-002 F-001/F-002 → closed（留痕/履约证据见矩阵）；
   - GOAL-003 F-001/F-002 → closed（无 epoch 控件引入；TIMEZONE_OPTIONS 留痕登记）；
   - GOAL-004 F-002/F-005（grouping 位序）与 F-006（币种目录）→ final residual（延续用户 2026-08-26 书面接受范围）；
   - **GOAL-004 F-007（安全整数）→ fixed**：`money.ts` `formatMoney` 对非安全整数值返回 `""`、`parseLocalizedMoney` 对超出安全范围的结果返回 `null`（守卫 + 快测 4 断言；`money.test.ts` 24/24）。
   - GOAL-004 F-008 → accepted（业务演示面不重开）。
4. 越界守卫：本轮代码改动仅 `apps/web/src/i18n/money.ts` + `money.test.ts`。

### 证据

| 主张 | 路径 / 命令 / commit |
|------|----------------------|
| 证据矩阵 | `GOAL-005-r4-evidence-closeout/attachments/r4-evidence-matrix.md` |
| F-007 加严 | `apps/web/src/i18n/money.ts`（`Number.isSafeInteger` 守卫）+ `money.test.ts` |
| 快测 | `npx vitest run src/i18n/money.test.ts`（24/24 · exit 0） |
| 全量回归 | `go test ./...`（全绿）；`npx vitest run`（88 files / 1181 tests · exit 0） |
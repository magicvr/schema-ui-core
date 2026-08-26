---
id: GOAL-004-r3-number-currency-semantics
doc: audit-entry
record_id: A-003
status: recorded
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
---

## A-003 · 响应 A-002（independent · grok-build）：required 闭合记录（2026-08-26）

- **source**：self（编排器 `/govern` 响应记录；**不是** independent）
- **auditor**：govern 编排器（本会话）
- **类型 / scope**：response · 对 A-002（独立审 grok-4.6 · high）全部 findings 的响应与闭合

### 响应 A-002 的 findings

| F | 主张 | 处置 | 关闭证据 |
|---|------|------|----------|
| **F-001**（high · required） | `updateLocalization.bodyMapping` 缺 `defaultCurrency`，Settings 保存路径丢弃该字段 | **fixed** | `settings.json` bodyMapping 补 `defaultCurrency`；Go 守卫测试 `schema/schema_test.go`（递归核对 Localization 表单字段 ⊆ bodyMapping，defaultCurrency 显式断言）；web 保存快测（`startup-config.test.tsx`「saves the Localization form incl. defaultCurrency …」断言 PATCH body 含 `defaultCurrency: "USD"`） |
| **F-002**（med · required） | 前端 branding 不消费 `defaultCurrency`，站点默认到不了运行时 | **fixed** | `branding.ts` 类型/解析/默认值补 `defaultCurrency`；`money.ts` 新增 `resolveEffectiveCurrency` + `formatMoney`/`parseLocalizedMoney` 的 `siteDefaultCurrency` 通道（优先级：显式 > 站点 > §4.3 映射）；`runtime.tsx` `/api/branding` 增读 `defaultCurrency` → `I18nState.defaultCurrency`（prop/fetch 双通道）；快测全覆盖（money/startup-config/runtime-timezone） |
| **F-003**（med · recommended） | Settings 页与 branding 快测缺失导致拦不住 F-001/F-002 | **fixed** | 新增：Go handler（PATCH/branding 投影/INVALID_DEFAULT_CURRENCY/重置断言）；Go schema 守卫；web 保存快测 + fetchBranding 解析断言 + 预填断言 |
| **F-004**（low · recommended） | 方案写 `"" \| "auto" \| ISO 4217`，实现与合同均无 `"auto"` | **fixed** | `D-001-r3-number-currency-plan.md` 回写：货币无 `"auto"` 语义（空 = 未配置），注明初稿误写已更正 |
| **F-005**（low · recommended） | 分组位序容差（`12,34.5`） | accepted-residual（范围 = R3 校验不校验位序；R4 核账时评估加严；代码注释 + 快测已文档化） | `money.ts` 注释 + `money.test.ts` 容差用例（用户确认 R4 处理） |
| **F-006**（low · recommended） | ISO 4217 为句法三字母非币种目录 | accepted-residual（本波按「句法近似」处理；如需完整目录由 R4/后续工作区评估） | E-003/E-004 记录；本表留痕 |
| **F-007**（low · recommended） | JS `number` 无安全整数守卫 | accepted-residual（常规金额安全；R4 核账评估 BigInt 或 `Number.isSafeInteger` 拒绝） | 本表留痕 |
| **F-008**（low · recommended） | 金额展示未接业务面 | accepted（不属于本波分母，VP-020 冻结「钱包金额语义化不重开演示面」；与 F-002 站点默认通道不同） | A-001 F-003 同款结论 + VP-020 首波冻结表 |
| **F-009**（low · recommended） | `migrate0052` 符号残留 | **fixed** | `migration.go`：`migrate0062` + 注释 `(0062 · …)` |
| **F-010**（low · recommended） | 台账 progress 不一致（frontmatter 0/6、正文 5/6、goal-tree 0/6；02-execution 索引缺 E-004） | **fixed** | `00-meta.md` frontmatter → `5/6`（version 0.2.0）；`goal-tree.md` → `5/6`；`02-execution.md` 索引补 E-004 |

### 闭合状态

- **开放 required = 0**（F-001/F-002 均 `fixed`，证据如上，可核对）。
- F-005/F-006/F-007：`accepted-residual`（范围与复审触发 = R4 核账），等待用户书面接受后结案（P-003 三路径；如用户不同意则将加严列为 R4 必做）。
- 复审路径：`go test ./...`（apps/api 全量）与 `npx vitest run`（apps/web 全量）在修复后复跑全绿（证据见 E-005）。

### 建议下一步

用户确认闭合记录（含 F-005/F-006/F-007 residual 接受）→ 可选：重新调用 grok build 对 F-001/F-002 修复做独立复审（对齐惯例）→ GOAL-004 `done`（R3 关门）→ R4 立项（GOAL-005 证据与关门，承接 F-002/F-005/F-006/F-007 核账）。
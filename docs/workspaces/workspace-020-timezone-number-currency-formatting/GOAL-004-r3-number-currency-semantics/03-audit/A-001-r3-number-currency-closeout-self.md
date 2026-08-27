---
id: GOAL-004-r3-number-currency-semantics
doc: audit-entry
record_id: A-001
status: recorded
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
---

## A-001 · R3 数字/货币语义关门自审（2026-08-26 · self leg）

- **source**：self
- **auditor**：govern 编排器（本会话）
- **类型 / scope**：close-out（self leg）· GOAL-004（R3）全量：合同 `GOAL-002 D-001` §3 / §4.1 / §4.3 一致性、money 工具、defaultCurrency 设置字段端到端（migration v62 / repository / handler / schema / 文案）、pin 同步、双全量回归、越界守卫
- **verdict**：**pass**（条件：independent 腿（grok build）无新增必改后关门；本条不含 independent 结果）

### 范围与区间

2026-08-26 立项至 C6 时点（C1～C5 done · progress 5/6）。审计模式 frozen `independent`（migration/API data 类）：self 先行（本条），随后本地 grok build（grok-4.6 · high）独立腿（A-002），编排器合并响应后由用户确认关门。共享资料引用：`shared_materials_catalog: none`。

### 成果（有证据）

1. **C1/C2/C3/C5（前端）**：`apps/web/src/i18n/money.ts`——`formatMoney`（Intl `style:currency`，无自建模板；机器值输入）、`defaultCurrencyFor`（zh-CN→CNY / en-US→USD / 兜底 USD，合同 §4.3，兑现 GOAL-002 F-002）、`parseLocalizedMoney`/`parseLocalizedNumber`（归一化到机器值；`null` 错误语义）、双向 round-trip；`money.test.ts` 20/20。
2. **C4（设置面 end-to-end）**：
   - migration **v62** `site_default_currency`（`ALTER TABLE site_settings ADD COLUMN default_currency TEXT NOT NULL DEFAULT ''`；v52 初始号与 core.persistence 冲突已更正，checksum 冻结 `74ede127…`）；
   - repository：字段/校验（空或大写三字母 ISO 4217）/PATCH merge/reset/SELECT+Scan/ErrNoRows 缺省；
   - handler：PATCH body + settings 行 + **公开 branding 投影** + `INVALID_DEFAULT_CURRENCY`；errorcatalog 双语；
   - settings.json Localization tab 字段 + responseMapping；web catalog 双语（field/error 两条目）。
   - **pin 同步**：store 迁移目录（applied 61→62、checksum 表、`completeFingerprintCatalogHead` 61→62、`lockedHeadExtraTables[62]={}`）、error_contract `frozenLiteralCodes`。
3. **回归证据**：`go test ./...` 全量绿（cmd/server、handler、store、settings 等）；web `vitest run` **88 files / 1175 tests** 全绿。
4. **越界守卫**：defaultCurrency 为站点设置列；Profile 默认集 / 模块矩阵 / Manifest / `docs/contracts/`（stage 门禁）未触碰；API 机器合同不变量（§3.3）未改。

### 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 1. 货币展示合同（ISO 4217 符号/位置/小数位，zh-CN/en-US 快测） | ✅ | `money.ts` + `money.test.ts`（formatMoney 双 locale + 显式 currency） |
| 2. 默认货币映射表（§4.3，缺省不抛错） | ✅ | `defaultCurrencyFor` + 快测（zh-CN→CNY、en-US→USD、fr-FR→USD）；GOAL-002 F-002 兑现 |
| 3. 输入解析归一化（int64 最小单位 / number；错误语义明确；不提交原文） | ✅ | `parseLocalizedMoney`/`parseLocalizedNumber`（null 语义 + 快测） |
| 4. defaultCurrency 端到端（API 列/PATCH/投影/Localization tab；不改 Profile 默认集） | ✅ | migration v62 + repository + handler + schema + 双语文案 + Go 单测 |
| 5. 展示↔输入双向一致（同一合同反向） | ✅ | round-trip 快测（en-US/zh-CN，含 0/JPY 场景） |
| 6. 无越界；机器合同不变量；`docs/contracts/` 未触碰；required=0（独立腿合并后） | ✅（待 A-002 合并后终判） | 提交范围 + §3.3 不变式 + pin 同步证据 |

### Findings

- **F-001 · migration 版本号冲突与更正留痕**（low · recommended · closed）
  - 描述：`default_currency` 迁移初定 Version 52，与 core.persistence（`mail_config`）全局版本冲突，测试报 `migration version 52 conflicts`；已改为 **62**（当前最大 61 + 1）并冻结 checksum。E-003 已留痕；本条目补充审计视角确认无历史库受影响（此前未发布）。
  - 证据：`settings/migration/migration.go`（62）；E-003「版本号冲突修正」；`go test ./internal/store/...` 全绿。
- **F-002 · 分组位序容差（`12,34.5` 类输入被接受）**（low · recommended · open → 移交 R4 核账）
  - 描述：`parseLocalizedNumberCore` 只做分隔符剥离，不验证分组位序；契约 §3.2 只要求「归一化 + 明确错误语义」，位序严谨性超出 R3 校验范围（代码注释 + 快测已文档化）。R4 关门核对时可评估是否加严。
  - 证据：`money.ts` 注释；`money.test.ts`「分组位序容差」用例。
- **F-003 · 金额展示尚未接线到业务展示面**（low · recommended · open → 不属于本波分母）
  - 描述：`formatMoney`/`defaultCurrencyFor` 为合同工具层；钱包等业务展示面（演示面）按 VP-020 边界**不进本波**（`不支持钱包金额语义化演示面重开`）。合同工具可用即可关门；接线属后续消费方。
  - 证据：VP-020 首波冻结表（货币行）；合同 §6 消费指引。

### 必改项汇总（required 列表）

无（0 条）。

### 结论 + 建议下一步

R3 交付完整、证据可核对；scope 内无 required 必改项，无到期 required 信息项（I-002/I-005 已 verified）；self leg **pass**。F-002/F-003 为 recommended，不阻断关门，随 R4 核账。下一步：按审计模式调用本地 grok build（grok-4.6 · high）执行 independent 腿（A-002）→ 编排器合并响应 → 用户确认关门 → R4 立项（GOAL-005）。
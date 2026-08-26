---
id: GOAL-004-r3-number-currency-semantics
doc: audit-entry
record_id: A-002
status: recorded
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
---

## A-002 · R3 数字/货币语义关门独立交叉审计（2026-08-26）

> 誊入说明：本条由编排器自本地 grok build（grok-4.6 · reasoning high）headless 会话原样誊入（`C:\Users\magicvr\.grok\bin\grok.exe -p -m grok-4.6 --effort high`；会话提示词见 E-004）。grok 按指令只出报告文本、未写入任何文件——落盘与索引由编排器完成，`source: independent` 保持不变。独立审当场复跑了 `go test ./...` 与 `npx vitest run`。

- **source**：independent
- **auditor**：grok-build grok-4.6 reasoning high
- **类型 / scope**：close-out · GOAL-004（R3）全量：合同 `GOAL-002/01-decision/D-001-r1-contract-freeze.md` §3 / §4.1 / §4.3；`money.ts` 工具；`defaultCurrency` 设置字段端到端（migration v62 / repository / handler / schema / errorcatalog）；pin 同步（store / error_contract）；越界守卫；回归证据（`go test ./...` 与 web vitest）
- **verdict**：**fail**

### 范围与区间

2026-08-26 立项至本独立审时点。C1～C5 在目标内被标为 done；C6 待关门。审计模式 frozen `independent`（migration / API data 类）。对照权威合同 = `GOAL-002-r1-contract-freeze/01-decision/D-001-r1-contract-freeze.md` §3 / §4.1 / §4.3；实施方案 = GOAL-004 `D-001-r3-number-currency-plan.md`。P-005：无新增 required 信息项；I-002 / I-005 已 verified（Root D-002）。本意见不改 `status` / `progress` / goal-tree。

### 成果（有证据）

1. **C1/C2/C3/C5 前端工具层可核对且独立复测通过**：`money.ts`（`formatMoney` Intl `style:currency` 无模板；`defaultCurrencyFor` §4.3；`parseLocalizedMoney`/`parseLocalizedNumber` 归一化 + `null`）`money.test.ts` **20/20**（独立审当场复跑）。GOAL-002 A-001 F-002（映射表）由 C2 兑现。
2. **C4 API 持久化面真实存在**：migration v62（`default_currency` 列）；repository（字段/校验/PATCH merge/reset/ErrNoRows 缺省 `""`；`TestRepositoryDefaultCurrencyPatch`）；handler（PATCH body、`settingsRow`、公开 `brandingResponse.DefaultCurrency`、`INVALID_DEFAULT_CURRENCY`）；errorcatalog + web catalog 双语。
3. **pin 同步到位**：`completeFingerprintCatalogHead = 62`；`lockedHeadExtraTables[62] = {}`；checksum `74ede1278137b3ce454255c87283315dbff69fed4e64cf18950bf9b8bb104391`；`operations_test`/`restart_test` applied 头 = 62；`frozenLiteralCodes` 含 `INVALID_DEFAULT_CURRENCY`。
4. **回归证据（独立审当场复跑）**：`apps/api` `go test ./...` 全绿；`apps/web` `npx vitest run` **88 files / 1175 tests** 全绿。
5. **越界守卫成立**：提交 `b2e9849f`/`c38e76d9` 未触碰 `docs/contracts/`、Profile 默认集、模块矩阵、Manifest、时区字段语义；钱包金额仍为 `int64` JSON；§3.3 API 机器合同不变量保持。

### 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 1. 货币展示（ISO 4217 符号/位置/小数位，zh-CN/en-US） | 达成 | `formatMoney` + `money.test.ts`（独立复跑 20/20） |
| 2. 默认货币映射表 §4.3（缺省不抛错） | 达成 | `defaultCurrencyFor` + 快测 |
| 3. 输入解析归一化（机器值；失败不提交原文） | 达成（有容差，见 F-005） | `parseLocalizedMoney`/`parseLocalizedNumber` → `null`；round-trip |
| 4. 站点级 `defaultCurrency` 端到端（API 列/PATCH/投影/Localization tab 字段；不改 Profile） | **未达成** | API 列/PATCH/JSON 投影有；**Settings 保存路径未接线**（F-001）；前端 branding 运行时丢弃该字段（F-002） |
| 5. 展示↔输入双向一致 | 达成（工具层） | round-trip en-US/zh-CN（含 0 / JPY `minorUnits: 0`） |
| 6. 无越界；§3.3 不变量；`docs/contracts/` 未触碰；required=0 | 越界/不变量达成；**required≠0** | 提交范围 + 本条 F-001/F-002 开放 |

### Findings

**F-001 · Localization 保存动作未映射 `defaultCurrency`（Settings 写路径断开）** · high · **required** · open
- 描述：`settings.json` Localization 表单加了 `defaultCurrency`，`responseMapping` 也加了（GET 能回填），但 `actions.updateLocalization.bodyMapping` **仍只有** `defaultLocale` / `siteTimezone`。渲染器 `request-construction.ts` 在存在 `bodyMapping` 时按白名单投影 PATCH body，未列入的字段不会发送。结果：管理员在「本地化」页填写并保存时，**站点默认货币不会被持久化**。自审 A-001 将标准 4 标为 ✅、C4 标为 done 的声明不成立。
- 证据：`apps/api/internal/modules/settings/schema/settings.json`（`updateLocalization.bodyMapping` vs 字段 + `responseMapping`）；`apps/web/src/protocol/conformance/request-construction.ts`（`bodyMapping` 白名单）；`git show c38e76d9 -- …/settings.json`。

**F-002 · 公开 branding 投影未被前端运行时消费** · med · **required** · open
- 描述：合同 §4.1 把 `defaultCurrency` 定义为站点级货币展示默认（空 = 走 §4.3 映射）。API `brandingGET` 已把字段加入公开投影，但 `apps/web/src/app/branding.ts` 的 `Branding` 类型、`fetchBranding`、`defaultBranding` 均无 `defaultCurrency`——该键在前端启动配置边界被丢弃；`formatMoney` 未传 `currency` 时只走 `defaultCurrencyFor(locale)`，从不读站点设置。对照 R2 已把 `siteTimezone` 接到同一 branding 通道；R3 未做对称接线。
- 证据：`settings.go`（API 投影有）；`branding.ts`（前端类型/解析无）；`money.ts`（缺省只走 locale 映射）；`startup-config.test.tsx`（断言不含 `defaultCurrency`）。

**F-003 · 方案要求的 Settings 页快测未落地，故未拦住 F-001/F-002** · med · recommended · open
- 描述：GOAL-004 D-001 明确要求「Settings 页快测断言（沿用 `startup-config.test.tsx` 模式）」。现有 `startup-config.test.tsx` 本地化 tab 只断言 locale + 时区，没有 `#field-defaultCurrency`、保存 body 含 `defaultCurrency`、或 `fetchBranding` 解析该字段；handler `settings_test.go` 的 branding/PATCH 用例同样未覆盖 `defaultCurrency` / `INVALID_DEFAULT_CURRENCY` HTTP。仓库级与 pin 测试全绿，无法发现 schema 写路径与 branding 消费缺口。
- 证据：`D-001-r3-number-currency-plan.md` §决定.3；`startup-config.test.tsx`；`settings_test.go`。

**F-004 · 冻结方案仍写 PATCH 接受 `"auto"`，实现与合同均不接受** · low · recommended · open
- 描述：D-001 方案写校验 `"" | "auto" | 有效 ISO 4217 三字母`；合同 §4.1 对货币是「ISO 4217 三字母；空 = 未配置」，**没有** `"auto"`（与 locale/timezone 不同）。`validateCurrency` 拒绝 `"auto"`。实现与合同一致、与方案不一致；应回写方案。
- 证据：`D-001-r3-number-currency-plan.md`；合同 §4.1；`repository.go` `validateCurrency`。

**F-005 · 分组位序容差（同意 A-001 F-002）** · low · recommended · open（随 R4 核账）

**F-006 · ISO 4217 仅为句法三字母，非币种目录** · low · recommended · open（合同写「ISO 4217」/方案写「有效 ISO 4217」；R3 未要求全目录——句法近似需措辞留痕）

**F-007 · 金额机器值用 JS `number`，无 int64/安全整数越界守卫** · low · recommended · open（超过 `Number.MAX_SAFE_INTEGER` 丢精度且无拒绝；R4 可评估 BigInt 或安全整数拒绝）

**F-008 · 金额展示未接到业务面（同意 A-001 F-003）** · low · recommended · open（不属于本波分母；注意 F-002 是站点默认通道、非业务展示消费）

**F-009 · migration 符号仍留 v52 旧名** · low · recommended · open（目录 Version 已是 62；`migrate0052`、`siteCurrencyDDL (0052 · …)` 残留，不影响 Apply 绑定，增加审计噪音）

**F-010 · 台账派生 progress 不一致（不作为代码门禁）** · low · recommended · open（`00-meta.md` frontmatter `progress: 0/6` 而正文 5/6；`goal-tree.md` 仍 0/6；`02-execution.md` 索引未列入 E-004——降低可核对性；progress 不得用于闭合 finding）

### 必改项汇总（required）

1. **F-001（high）**：把 `defaultCurrency` 加入 `updateLocalization.bodyMapping`（并补 Localization 保存快测）。未修不得把 C4 / 成功标准 4 视为达成。
2. **F-002（med）**：`branding.ts`（及启动配置测试）解析/默认化 `defaultCurrency`；运行时在站点值非空时把它作为 `formatMoney`/`parseLocalizedMoney` 的显式 `currency`，空则回退 §4.3 映射。

F-003～F-010 不阻断修完 F-001/F-002 后的复审，但 F-003 的快测应与 F-001/F-002 一并做，避免再漏。

### 与既有意见的异同（A-001 self）

| 项 | A-001 self | 本独立审 |
|----|------------|----------|
| C1/C2/C3/C5 工具 + round-trip | pass | 同意（当场复测 20/20） |
| C4 API migration/repository/handler/pin | pass | 同意（API 面真实；pin 同步真实） |
| C4「端到端 / Localization tab」 | ✅ | **不同意**：写路径未接线（F-001） |
| branding「运行时取站点货币默认」 | 作为 C4 成果陈述 | **不同意**：前端未消费（F-002） |
| 分组容差 / 未接业务展示面 | F-002/F-003 recommended | 同意（本条 F-005/F-008） |
| v52→v62 冲突 | F-001 closed | 同意闭合；残留 `migrate0052` 命名（F-009） |
| required | 0 | **2**（F-001 high，F-002 med） |
| verdict | pass（条件：independent 无新必改） | **fail** |

无与 self 相反的「应放行」意见；冲突点是 C4 完成声明，按 P-004 由编排器展示并由用户裁决是否 residual / overruled（独立审建议：**修 F-001/F-002，不要 overruled**）。

### 结论 + 建议给编排器 / 用户的下一步

工具层（§3.1/§3.2/§4.3）与 API 列/校验/JSON 投影/pin/回归都经得起核对；**关门阻断在 C4 最后一公里**：Settings Localization **不能保存** `defaultCurrency`，公开投影也**到不了**前端运行时。因此 close-out **不能**无条件放行，verdict **fail**。

建议 `/govern` 下一步：

1. 响应本 A-002；F-001/F-002 走 `fixed`（补 `bodyMapping` + branding 解析/回退 + 快测），不要 overruled。
2. 复跑：`go test ./...`（`apps/api`）与 `npx vitest run`（`apps/web`），至少覆盖 settings schema 保存、`fetchBranding`、`INVALID_DEFAULT_CURRENCY`。
3. 独立审复审 F-001/F-002 闭合后，再谈 C6 关门与 R4 立项。
4. F-004 方案 `"auto"` 措辞随决策回写；F-005/F-006/F-007 可进 R4 核账。

### 声明

本意见 **source: independent**，不修改目标 `status` / `progress` / 方案正文 / goal-tree。落盘与索引由编排器完成（本条文件 + `03-audit.md` 索引 A-002），响应由 `/govern` 处理。
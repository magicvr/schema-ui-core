---
id: GOAL-004-r3-number-currency-semantics
doc: audit-entry
record_id: A-004
status: recorded
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
---

## A-004 · A-002 F-001/F-002 闭合独立复审（2026-08-26）

- **source**：independent
- **auditor**：grok-build grok-4.6 reasoning high
- **类型 / scope**：finding-closure · A-002 required F-001（Localization `updateLocalization.bodyMapping` + `schema_test.go` 守卫 + `startup-config` 保存快测）与 F-002（`branding.ts` 解析 `defaultCurrency` / `money.ts` `resolveEffectiveCurrency`+`siteDefaultCurrency` 通道 / runtime `I18nState.defaultCurrency` prop+fetch）。**不含**全量关门、C6、A-002 F-003～F-010。
- **verdict**：**pass**

### 范围与区间

工作区：`workspace-020-timezone-number-currency-formatting`（Root `GOAL-001-timezone-number-currency-formatting`；canonical `docs/workspaces/workspace-020-timezone-number-currency-formatting/`；`shared_materials_catalog: none`；`primary_plan` = VP-020）。对照合同：同区 `GOAL-002-r1-contract-freeze/01-decision/D-001-r1-contract-freeze.md` §4.1 / §4.3。对照前意见：A-002（independent · fail · 2 required）与 A-003（self 响应 · 声称 F-001/F-002 `fixed`）。P-005：本目标无新增 required 信息项；I-002 / I-005 已 verified（Root D-002），不影响本闭合门禁。本意见不改 `status` / `progress` / 方案正文 / goal-tree。

### 成果（有证据）

独立复审当场核对源码并复跑相关测试（非编排器转述）：

| 包 / 文件 | 结果 |
|-----------|------|
| `apps/api` `go test ./internal/modules/settings/schema/ ./internal/handler/ -count=1` | **ok**（schema 0.498s；handler 33.614s；含 `TestLocalizationBodyMappingCoversFormFields` 与 PATCH/branding `defaultCurrency` 断言） |
| `apps/web` `npx vitest run src/app/startup-config.test.tsx src/i18n/money.test.ts src/i18n/runtime-timezone.test.tsx` | **45/45**（3 files；`money.test.ts` 23；`runtime-timezone.test.tsx` 7；`startup-config.test.tsx` 15） |

### 对照 A-002 关闭证据

| A-002 finding | 关闭主张 | 本复审状态 | 可核对证据 |
|---------------|----------|------------|------------|
| **F-001** high · required · Settings 保存路径未映射 `defaultCurrency` | A-003 `fixed` | **closed（证据充分、可重复）** | 见下「F-001 核验」 |
| **F-002** med · required · 公开 branding 投影未被前端运行时消费 | A-003 `fixed` | **closed（证据充分、可重复）** | 见下「F-002 核验」 |

#### F-001 核验

1. **`settings.json` `updateLocalization.bodyMapping` 含 `defaultCurrency`**  
   `apps/api/internal/modules/settings/schema/settings.json`：`updateLocalization.bodyMapping` 现为 `defaultLocale` / `siteTimezone` / **`defaultCurrency: "defaultCurrency"`**。Localization 表单字段同三键；`recordSource.responseMapping` 亦含 `defaultCurrency`（GET 回填）。渲染器 `apps/web/src/protocol/conformance/request-construction.ts` 在存在 `bodyMapping` 时仍按白名单投影 PATCH body——A-002 描述的丢弃机制仍在；现白名单已覆盖该字段。
2. **`schema_test.go` 守卫**  
   `apps/api/internal/modules/settings/schema/schema_test.go` `TestLocalizationBodyMappingCoversFormFields`：从 `SchemaDocuments()["settings"]` 解码真实文档；递归 `body` 树找到 `submitAction == "updateLocalization"` 的唯一表单（`checked == 1`）；断言每个 field `id` ⊆ `bodyMapping`；并显式断言 `bodyMapping["defaultCurrency"] == "defaultCurrency"`。本复审 `go test ./internal/modules/settings/schema/` **ok**。
3. **startup-config 保存快测**  
   `apps/web/src/app/startup-config.test.tsx`：真实 `settings.json` 路径加载；本地化 tab 预填断言 `#field-defaultCurrency` = `CNY`；用例「saves the Localization form incl. defaultCurrency through the real PATCH action (R3 F-001)」改值为 `USD` 后提交，断言 PATCH `/api/settings/default` body `defaultCurrency === "USD"`（及 `defaultLocale`）。本复审该文件 **15/15**。

原缺陷（bodyMapping 仅 locale+timezone → 保存静默丢弃货币）已不成立。

#### F-002 核验

1. **`branding.ts` 解析 / 默认化**  
   `apps/web/src/app/branding.ts`：`Branding.defaultCurrency: string`；`fetchBranding` 读 `defaultCurrency` 并 `toUpperCase()`（缺省 `""`）；`defaultBranding()` 为 `""`（合同 §4.1 空 = 未配置 → §4.3 映射）。`startup-config.test.tsx`：完整 payload 断言 `CNY`；稀疏 payload 断言 `""`。
2. **`money.ts` 站点默认通道**  
   `resolveEffectiveCurrency(locale, siteDefault)`：有效 ISO 4217 站点值优先，否则 `defaultCurrencyFor`（§4.3）。`formatMoney` / `parseLocalizedMoney` 的 `siteDefaultCurrency` 经私有 `resolveCurrency`：显式 `currency` > 站点默认 > 内嵌映射；非法站点值回退映射、不抛错。`money.test.ts` 覆盖站点优先 / 显式覆盖站点 / 空与非法回退 / 解析符号剥离。本复审 **23/23**。
3. **runtime `I18nState.defaultCurrency` prop + fetch**  
   `apps/web/src/i18n/runtime.tsx`：`I18nState.defaultCurrency`；测试缝 `siteDefaultCurrency`；`systemDefaultUrl` fetch 读 `record.defaultCurrency`（trim + uppercase；空 → `null`）；`effectiveSiteDefaultCurrency = siteDefaultCurrency ?? fetchedSiteDefaultCurrency`，写入 state（空串 = 未配置）。生产装配 `apps/web/src/main.tsx` 两处均为 `<I18nProvider systemDefaultUrl="/api/branding">`，与 R2 `siteTimezone` 对称，fetch 通道即生产路径。`runtime-timezone.test.tsx` 用例「R3 F-002: site default currency reaches the runtime (prop + fetch)」断言 prop=`CNY`、fetch=`USD`。本复审该文件 **7/7**。

原缺陷（前端启动配置丢弃该键；`formatMoney` 无站点通道）已不成立。业务面金额展示未接 `formatMoney` 仍属 A-002 **F-008**（recommended / 已 accepted；**不在本复审 scope**），不是 F-002 未闭合。

### 对照成功标准（若适用）

本条只审 A-002 F-001/F-002 闭合，不重判 GOAL-004 全量成功标准。与这两条直接相关的标准 4（站点级 `defaultCurrency` 端到端含 Localization 保存与公开投影消费）在 **写路径 + 前端消费通道** 上现可核对为达成；业务展示接线仍按 F-008 排除本波。

### Findings

本复审 **无新增** required / recommended finding。

A-002 F-001 / F-002 在本条视为 **closed**（闭合路径 = `fixed`，证据上表；不改 A-002 原文历史状态——由 `/govern` 响应本条）。

### 必改项汇总

**无。** 本 scope 开放 required = 0。

### 与既有意见的异同

| 项 | A-002 independent | A-003 self 响应 | 本条 A-004 |
|----|-------------------|-----------------|------------|
| F-001 bodyMapping | required · open（保存丢弃） | `fixed` | **同意闭合**（源码 + schema 守卫 + 真实表单 PATCH 快测 + 当场复跑） |
| F-002 branding/runtime | required · open（前端丢弃；money 不读站点） | `fixed`（解析 + money 通道 + I18nState prop/fetch） | **同意闭合**（与用户指定核验点一致；生产 `systemDefaultUrl` 已接） |
| F-003～F-010 | recommended | 部分 fixed / residual / accepted | **未审**（超出本 scope） |
| verdict | fail | —（response） | **pass**（仅 F-001/F-002 闭合） |

无与 A-003 相反的「不得闭合」意见；无 P-004 冲突需用户裁决。

### 结论 + 建议给编排器 / 用户的下一步

A-002 两条 required 的修复声明经独立复核对账成立，verdict **pass**。本条 **不** 放行 GOAL-004 `status: done` / C6 关门（关门仍属编排器 + 用户确认；A-002 F-005/F-006/F-007 residual 是否书面接受亦不在本 scope）。

建议 `/govern`：

1. 响应本 A-004：将 A-002 F-001 / F-002 合法闭合为 `fixed`（本条为独立核对证据）。
2. 再处理 C6：用户确认关门；F-005/F-006/F-007 若仍待书面 residual，按 P-003/P-004 留痕后再标 `done`。
3. 不要用本条 pass 覆盖未审的 recommended 项或 R4 核账。

### 声明

本意见 **source: independent**，auditor = grok-build grok-4.6 reasoning high。不修改目标 `status` / `progress` / 方案正文 / goal-tree。响应由 `/govern` 处理。

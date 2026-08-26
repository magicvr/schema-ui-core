# R4 证据矩阵 · workspace-020（GOAL-005 · 2026-08-26）

> 权威合同 = `GOAL-002-r1-contract-freeze/01-decision/D-001-r1-contract-freeze.md` §1～§6。本矩阵把 R1～R3 交付逐条映射到可核对证据。全量回归基线：`apps/api` `go test ./...` 全绿（含 settings/handler/store pin）；`apps/web` `npx vitest run` 88 files / **1181** tests 全绿（F-007 加严后 1180 → 1181）。

## §1 范围与落点

| 条款 | 证据 |
|------|------|
| 展示/输入语义 = 前端；API = 机器合同 | `apps/web/src/i18n/format.ts`、`money.ts`（Intl 展示 + 解析归一化）；API 无 locale 序列化（§3.3 复证） |
| 内嵌默认 zh-CN + auto 可运行 | `locale.ts` resolveLocale 兜底 en-US；`timezone.ts` L4 auto；`money.ts` §4.3 映射；无配置启动无硬依赖（全量回归通过） |

## §2 时区来源（R2）

| 条款 | 证据（快测） |
|------|--------------|
| L1 用户覆盖（schema-ui:timezone 单通道） | `timezone.test.ts`（15 例）L1 优先/auto 穿透/无效穿透；`runtime-timezone.test.tsx` 持久化翻转 + auto 移除 key |
| L2 会话探测（Intl） | `timezone.test.ts` L2 探测/空·无效跳过；`detectBrowserTimezone` 合法性 |
| L3 站点默认（siteTimezone） | `runtime-timezone.test.tsx` C3 fetch（/api/branding siteTimezone → Europe/London） |
| L4 auto 兜底 / 无效 IANA 降级 / 不抛错 | `timezone.test.ts` L4 + hostile inputs；`format.test.ts` 无效时区降级 |
| 展示统一语义（formatDate 默认生效时区） | `runtime-timezone.test.tsx` 格式翻转 + 显式覆盖；（输入侧结论：renderer 无 epoch 控件，date-only 本地日语义——E-003 留痕） |

## §3 数字/货币语义（R3）

| 条款 | 证据（快测） |
|------|--------------|
| 货币展示（Intl style:currency 无模板） | `money.test.ts`（24 例）formatMoney zh-CN/en-US/显式 currency/JPY minorUnits=0 |
| 默认货币映射（§4.3） | `defaultCurrencyFor` zh-CN→CNY / en-US→USD / fr-FR→USD |
| 站点默认通道（§4.1） | `runtime-timezone.test.tsx` prop=CNY / fetch=USD；`money.test.ts` siteDefaultCurrency 优先级 |
| 输入解析归一化（机器值 / null 语义） | `money.test.ts` parseLocalizedMoney/Number 全用例 + 双向 round-trip（C5） |
| API 机器合同（§3.3）不变量 | `go test ./...`（handler wallet/settings 断言 int64 JSON + RFC3339 UTC）；R3 未改序列化（E-003） |
| 安全整数守卫（R4 F-007） | `money.test.ts` 超出 MAX_SAFE_INTEGER → format "" / parse null |

## §4 设置归属与字段（R3）

| 条款 | 证据 |
|------|------|
| Localization tab：defaultLocale/siteTimezone/defaultCurrency | `settings.json`（字段 + responseMapping + bodyMapping）；`startup-config.test.tsx` 预填断言（#field-defaultCurrency=CNY） |
| 保存路径（F-001 守卫） | `startup-config.test.tsx` 真实 PATCH body 断言（defaultCurrency=USD）；`schema_test.go` bodyMapping 覆盖守卫；Go handler PATCH/branding 断言 |
| 用户级覆盖通道（时区，头部） | `timezone-switcher.test.tsx`（4 例）+ App.tsx 挂载 |
| 公开投影消费（F-002 守卫） | `startup-config.test.tsx` fetchBranding defaultCurrency=CNY / 稀疏=空 |

## §5 越界核账（逐项）

| 越界项 | 核账结论 | 证据 |
|--------|----------|------|
| 汇率/换算/计费/结算/发票 | 未引入（无语义/代码） | §5 无实现；wallet 金额仍 int64 机器值 |
| DB `timestamptz` 持久化合同（RT-T03） | 未引入 DDL | 本波新增迁移仅 `site_settings.default_currency`（v62，TEXT 列） |
| 多时区排程/日历/外部时区服务 | 未引入 | 无相关代码 |
| 翻译/文案中心（VP-007 重开） | 未重开 | 仅消费 runtime |
| 热加载 | 未引入 | — |
| 改 Profile 默认集/模块矩阵/Manifest | 未改 | 提交范围核对（R1～R3 commits） |
| `docs/contracts/` stage 门禁 | 未触碰 | 提交范围核对 |
| 钱包演示面重开（VP-011） | 未重开 | formatMoney 未接业务面（F-008 accepted） |
| API 机器合同不变量 | 保持 | §3.3；Go handler 断言 |

## 双 locale 范例（合同 §6 · zh-CN / en-US 各至少一场景）

| 场景 | zh-CN | en-US |
|------|-------|-------|
| 时区展示（L1 覆盖生效） | `runtime-timezone.test.tsx`（set-shanghai → Asia/Shanghai 格式翻转） | 同文件（PROBE America/New_York） |
| 货币展示（站点默认） | `startup-config.test.tsx` #field-defaultCurrency=CNY + fetchBranding CNY | `money.test.ts` formatMoney en-US USD `$123.45` / site CNY 覆盖 |
| 货币输入解析（round-trip） | `money.test.ts` ¥123.45 → 12345；round-trip zh-CN | `money.test.ts` $1,234.56 → 123456；round-trip en-US |
| 无效输入错误语义 | parse null（zh/en 通用断言） | 同左 |

## 核账项处置汇总（C3）

| 项 | 处置 |
|----|------|
| GOAL-002 F-001（数字字段解释） | closed（R1 D-001 §4.1 已留痕：无独立数字字段） |
| GOAL-002 F-002（默认货币映射表） | closed（GOAL-004 C2 履约 + 快测） |
| GOAL-003 F-001（epoch 输入控件按 §2.3） | closed（R4 核对 renderer 无含时间控件；date-only 本地日语义） |
| GOAL-003 F-002（TIMEZONE_OPTIONS 扩展留痕） | closed（switcher 常量可核对；扩展须留痕——本账登记） |
| GOAL-004 F-002/F-005（grouping 位序） | final residual（用户 2026-08-26 书面接受；文档 + 快测留痕；R4 不加严） |
| GOAL-004 F-006（币种目录） | final residual（句法三字母 + §4.3 映射表权威；全目录属后续可选项） |
| GOAL-004 F-007（安全整数） | **fixed**（format "" / parse null 守卫 + 快测，2026-08-26） |
| GOAL-004 F-008（业务展示接线） | accepted（VP-020 冻结「钱包金额语义化不重开演示面」） |
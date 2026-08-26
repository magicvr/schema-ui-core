---
id: GOAL-002-r1-contract-freeze
doc: decision-entry
record_id: D-001
status: accepted
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
---

## D-001 · R1 格式语义合同冻结正文

> 本文是 VP-020 首波（R1 合同冻结）的**权威合同正文**，供 R2（时区语义）、R3（数字/货币语义）直接消费；R4 按其核对（快测 + `zh-CN`/`en-US` 双 locale 范例）。前置信息门禁 I-001 / I-002 / I-005 已由用户裁决 accepted（Root `GOAL-001-.../01-decision/D-002-r1-info-adjudication.md`），本文是对裁决的合同化落笔。

### §0 效力与引文

- 归属：workspace-020（`workspace_id = workspace-020-timezone-number-currency-formatting`）· Root `GOAL-001-timezone-number-currency-formatting` · 子目标 `GOAL-002-r1-contract-freeze`。
- 上游：VP-020（`active` · `vision_ref = schema-ui-core-admin-foundation@0.2.0`）；消费 VP-007 locale 运行时、VP-005 设计系统。
- 生效时机：`status: accepted`；正文与 Root `D-002` 冲突时以本条为准并回写修订。
- **审计模式**：`self`（低风险、可逆、文档型合同）。不直接改时区 / 格式相关 DDL 或迁移台账；`docs/contracts/`（skills consumer 镜像，stage 门禁）不在本波交付范围。

### §1 范围与落点（总则）

1. 本波交付**展示 / 输入 / 序列化语义面**（应用层），不交付持久化时区合同（架构 RT-T03 仍 `registered`）、不交付汇率 / 换算 / 计费 / 结算（业务域）。
2. **落点划分**（I-002 裁决）：
   - **展示与输入解析的 locale 语义 = 前端**（`Intl.*` 运行时 + 输入归一化）。
   - **API = 机器合同**：与 locale 无关；时间一律 UTC、金额一律整数最小单位、数字一律无格式 JSON number。API 机器格式在本合同 §3.3 文档化，R3 实现时不得改变。
3. **内嵌默认**（无任何配置可运行）：有效 locale = `zh-CN`（站点 `defaultLocale` 未配置时按 VP-007 `auto` 解析 → 浏览器/系统 → 站点）；时间语义 = `auto`（站点 `siteTimezone` 未配置时按会话探测）；货币展示 = 无配置时按有效 locale 映射默认货币（§4.3）。任何配置缺失**不得**成为 mvp/dev 启动硬依赖。

### §2 时区来源合同（I-001 裁决）

**解析优先级（从高到低，取第一个可用的确定性来源）：**

| 级 | 来源 | 载体 | 说明 |
|----|------|------|------|
| L1 | 用户级覆盖 | `localStorage["schema-ui:timezone"]`（IANA 名） | 与 locale 同款单通道（先例 `schema-ui:locale`）；`"auto"` 或缺失 = 不覆盖；登录/登出不清除 |
| L2 | 会话级 auto 探测 | 客户端探测（`Intl.DateTimeFormat().resolvedOptions().timeZone`） | 仅当 L1 未覆盖时生效；探测结果可能为空串 |
| L3 | 站点默认 | `SiteSettings.siteTimezone`（`/api/branding` 投影 `siteTimezone`） | 站点管理员配置；`"auto"` / 空 = 未配置 |
| L4 | 内嵌默认 | `auto` | 上述均不可用时的最终兜底（不抛错、不白屏） |

**语义规则：**

1. 「timezone 生效值」= IANA 名（如 `Asia/Shanghai`、`America/New_York`）或 `auto`。
2. **展示**：时间展示统一使用生效时区（`Intl.DateTimeFormat(locale, { timeZone })`）；无效 IANA 名降级到 locale 默认时区（沿用 `apps/web/src/i18n/format.ts` 既有 fail-safe 行为），**不得**抛错。
3. **输入**：时间输入解析（R2 实现）与展示使用同一生效时区语义——「统一语义」= 同一时刻在展示面与输入面呈现/解释一致；输入面解析出的时间点（epoch）与展示面互换不产生偏移。
4. **偏移展示**：IANA 语义展示可附 UTC 偏移（R2 实现细节）；偏移值由 Intl 派生，不手工维护。
5. 本波**不**引入 DB 用户时区字段、不引入服务端按请求时区格式化（API 机器合同见 §3.3）。

### §3 数字 / 货币语义落点合同（I-002 裁决）

#### §3.1 展示合同（前端）

- 有效 locale 驱动（`Intl.NumberFormat` / `Intl.DateTimeFormat`）：千分位分隔、小数位、百分比、货币符号 / 位置 / 小数位均由 locale + 选项决定，**不**自建格式模板（延用 VP-007「不暴露任意日期/数字格式模板」原则）。
- 展示输入源：API 机器值（§3.3）→ 前端格式化 → 展示面。

#### §3.2 输入解析合同（前端）

- 用户输入的 locale 化数字（千分位 / 小数点 / 货币符号）在前端解析归一化为**机器值**（§3.3 格式）后提交 API。
- 解析失败 = 明确的输入错误反馈（错误文案走 S4 前台本地化通道）；**不得**向 API 发送未归一化的 locale 字符串。
- R3 交付具体解析器与错误语义；本波冻结「归一化后才提交」这一落点。

#### §3.3 API 机器合同（文档化；R3 实现不得改变）

| 类型 | 机器格式 | 现状依据 |
|------|----------|----------|
| 时间戳 | RFC3339 毫秒 **UTC** 字符串（`2006-01-02T15:04:05.000Z`，形如 `2026-08-26T07:00:00.000Z`）；输入侧接受 RFC3339（可含 offset，解析后归一 UTC） | `apps/api/internal/handler/rfc3339.go` + 各 handler `UTC().Format(...)`；`parseOperationTime` 接受 YYYY-MM-DD / RFC3339 |
| 金额 | **int64 JSON 数字**（最小货币单位；如 `amountDelta`、`balanceTotal`） | `apps/api/internal/handler/wallet.go`（`json:"amountDelta"` 等） |
| 普通数字 | 无格式 JSON number | 全局现状 |
| 日期查询参数 | `YYYY-MM-DD`（含当天边界的本地日语义由 handler 定义）或 RFC3339 | `INVALID_DATE_FILTER` 现状 |

- **不变式**：API 不随 locale / 时区改变输出格式；时区/格式语义完全由前端消费端决定。
- R3 不得为「展示方便」改造 API 序列化；如需新增数字字段，沿用上表机器格式。

### §4 设置归属与字段（I-005 裁决）

#### §4.1 站点默认（Settings → Localization tab）

| 字段 | 键 / 现状 | 说明 |
|------|-----------|------|
| 默认语言 | `defaultLocale`（已有） | 站点默认 locale；`auto` = 解析交给客户端 |
| 默认时区 | `siteTimezone`（已有） | 站点默认时区（L3）；`auto` / 空 = 未配置 |
| 默认货币 | `defaultCurrency`（**新增**，ISO 4217 三字母代码） | 站点级货币展示默认；空 = 未配置（§4.3） |
| 数字格式 | （无独立字段） | 跟随有效 locale，按 §3.1 展示 |

- 新增 `defaultCurrency` 属设置面扩展（站点级），不是 Profile 默认集 / 模块矩阵 / Manifest 变更。
- R2 实现时若需改 settings schema / Settings 页 schema，须走本波内契约（不触碰 `docs/contracts/` 镜像门禁；表单行为见 VP-007 既有 Settings 机制）。

#### §4.2 用户级覆盖（头部 locale 通道）

- 时区选择与头部 locale 切换并列（同 UI 区），持久化于 `localStorage["schema-ui:timezone"]`（L1；`auto` = 移除 key）。
- 用户级覆盖**不含**货币/数字格式（无用户级货币偏好字段；货币默认仅站点级，避免偏好面膨胀——如需可在后续工作区扩展）。

#### §4.3 内嵌默认货币映射（无配置时）

无显式 `defaultCurrency` 时按有效 locale 映射（R3 可扩展表）：`zh-CN → CNY`；`en-US → USD`。此表属合同；R3 交付时不要求覆盖全部 ISO 4217，仅要求映射表可核对、缺省不抛错。

### §5 越界与退出分母

未进本波（越界即合同违约）：

1. DB `timestamptz` 持久化时区合同（架构 RT-T03，仍 `registered`）。
2. 汇率 / 换算 / 计费 / 结算 / 发票（业务域）。
3. 多时区排程 / 日历 / 提醒；外部时区服务依赖；强制全站统一 UTC / 固定时区启动。
4. 翻译与文案中心（VP-007）；热加载；改 Profile 默认集 / 模块矩阵 / Manifest 装配语义；重开 VP-007 / VP-012。
5. 每个页面/模块自选格式引擎（自定义格式模板）；格式相关安全（时区炸弹类输入）与符合性 gap 归持续程序（VP-009/010 边界）。

### §6 消费指引与核对

- **R2（GOAL-003）消费**：§2 时区来源合同（L1～L4 优先级 + 统一语义 + 降级）+ §4.2 用户覆盖通道实现。
- **R3（GOAL-004）消费**：§3 数字/货币展示与输入解析合同 + §4.3 默认货币映射表。
- **R4（GOAL-005）核对**：快测 + `zh-CN` / `en-US` 至少各一场景，逐条对照 §2～§4；无越界（§5）；open required = 0。
- 合同修订：任何语义变更须回写本文并留痕（版本号递增）；未经用户裁决的修订不生效。
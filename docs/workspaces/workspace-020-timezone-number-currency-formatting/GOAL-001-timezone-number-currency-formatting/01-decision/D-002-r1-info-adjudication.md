---
id: GOAL-001-timezone-number-currency-formatting
doc: decision-entry
record_id: D-002
status: accepted
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
---

## D-002 · R1 信息门禁裁决（I-001 / I-002 / I-005）

### 触发

2026-08-26，R1（合同冻结）立项前，P-005 required 信息项 I-001、I-002 到期（最晚阶段 = R1 方案冻结），I-005（non-blocking，最晚 R2）一并由 lead 提案。编排器按 P-004.4 向用户提出裁决问题（附代码基证据：`apps/web/src/i18n/format.ts` + `runtime.tsx` 的 Intl.* 展示与 localStorage locale 通道；`apps/api/internal/handler/settings.go` 的 `siteTimezone`/`defaultLocale`；`apps/api/internal/handler/wallet.go` 的 int64 最小单位金额；各 handler 的 RFC3339 毫秒 UTC 时间序列化）。

### 决定（用户 2026-08-26 书面采纳，均选推荐项）

1. **I-001 时区来源**：`会话级 auto 探测 + 站点兜底 + 用户级 localStorage 覆盖`。
   - 解析优先级（从高到低）：用户级覆盖（localStorage 单通道，与 locale 同模式 `schema-ui:*`）→ 会话级 auto 探测（`Intl.DateTimeFormat().resolvedOptions().timeZone` 等客户端探测，无效 IANA 名降级）→ 站点默认 `siteTimezone` → 内嵌默认（`auto`）。
   - **零 DB schema 变更**；不引入用户表/设置表新字段；不改 Profile 默认集。
2. **I-002 数字 / 货币语义落点**：`前端落点`。
   - 展示与输入解析的 locale 语义全部在前端（`Intl.*` + 输入归一化）；API 保持**机器合同**：时间 = UTC RFC3339 毫秒字符串、金额 = int64 最小单位 JSON 数字、普通数字 = 无格式 JSON number。
   - API 机器合同在 R1 合同正文（GOAL-002 D-001）文档化，R3（数字/货币语义）按此落地并核对。
3. **I-005 设置归属**：`站点默认进 Settings → Localization tab；用户级覆盖并入头部 locale 通道`。
   - 站点默认字段：`defaultLocale`（已有）、`siteTimezone`（已有）、`defaultCurrency`（新增，ISO 4217）；数字格式无独立字段（跟随有效 locale）。
   - 用户级覆盖：头部 locale 切换旁新增时区选择，同 `schema-ui:*` localStorage 单通道持久化；`auto` = 移除 key。

### 为什么

- 与已交付模式一致：locale 已用「站点默认 + 浏览器探测 + localStorage 用户覆盖」三级，时区复用同一通道，学习成本与实现成本最低；不越「DB timestamptz 持久化合同不进本波」边界。
- API 现状（RFC3339 UTC / int64 最小单位 / JSON number）即机器合同；前端落点零破坏性变更、可逆，满足「内嵌默认无配置可运行」。
- Localization tab 已承载 locale + timezone，货币默认字段同 tab 内聚；用户级覆盖不新增设置 tab 语义。

### 未选方案

- 用户级时区存库（DB 用户设置字段）：多端一致更强，但引入本波用户级 schema 变更、放大范围，且与「不改 Profile 默认集」边界相邻；作为后续架构议题（如需）另行立项。
- API 序列化携带 locale 语义（decimal 字符串、按请求 locale 输出）：破坏性变更、跨 handler 迁移成本高，超出 Admin 展示/输入合同需要。
- 数字格式独立站点字段：与「locale 驱动格式」意图重复，弃用；格式随有效 locale，仅货币默认需显式字段。

### 影响

- 解锁 R1 方案冻结：GOAL-002（合同冻结）可立项并可写合同正文。
- Root 信息台账：I-001 / I-002 → **verified**（用户裁决 accepted，证据 = 本决策 + 用户 2026-08-26 问答答复）；I-005 → **verified**（lead 提案 + 用户确认）。
- I-003 / I-004 保持 VP 冻结投影（registered，不进）。
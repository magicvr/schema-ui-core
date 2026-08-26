---
doc_type: vision-plan
id: VP-020-timezone-number-currency-formatting
title: 时区 / 数字 / 货币格式语义
status: closed
vision_ref: schema-ui-core-admin-foundation@0.2.0
lead_workspace: workspace-020-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-27
version: 0.3.0
parent: null
---

# VP-020 · 时区 / 数字 / 货币格式语义

## 状态与门闩（2026-08-26 · **active**）

| 项 | 值 |
|----|-----|
| status | **`active`**（2026-08-26 激活 · VRev-044 self `pass`；Admin 类 freshness PASS `66f5fd1f` → `c6fda691`） |
| **lead_workspace** | `workspace-020-timezone-number-currency-formatting`（Root `GOAL-001-timezone-number-currency-formatting` **active** · 0/4；唯一 delivery） |
| **Vision required** | VRev-044 self `pass`（V-F079/V-F080 → 激活事务内 fixed）仍成立 |
| **推进门闩** | 已解除（2026-08-26 激活并开区）；既定新鲜度锚点 = `c6fda691`，供后续 VP 消费 |
| **组合位置** | Admin 功能分支 · 基架能力剩余 #5；消费 VP-007 locale 运行时与 VP-005 设计系统 |
| **完整 ≠ 架构持久化合同** | 本 VP 只交付**展示 / 输入 / 序列化语义面**。DB `timestamptz` 持久化时区合同仍归架构分支 RT-T03（`registered`，不进本波） |

## 意图

在 VP-007 已交付的多语种运行时（`zh-CN` / `en-US` + `auto` 解析）与 VP-005 设计系统之上，把「时间 / 时区 / 数字 / 货币」的**格式语义**收成可核对的 Admin 合同：

1. **时区**：会话/用户级时区解析与展示（IANA 名称 + UTC 偏移 + `auto` 探测）；时间展示与输入统一走同一时区语义。
2. **数字**：locale 驱动的千分位、小数位、百分比等格式合同；展示面与输入解析面一致（同一定义，双向可用）。
3. **货币**：ISO 4217 货币代码 + 符号 / 位置 / 小数位合同。**不是**汇率、换算或计费。
4. **内嵌默认**：无显式配置时 `zh-CN` 语义可运行；不得把「必须配置时区/货币」做成 mvp/dev 启动硬依赖。

本 VP 属 **Admin 功能分支**。持久化层的时区合同（DB `timestamptz`）属架构分支 RT-T03（仍 `registered`），本 VP 只消费其语义边界、不越界实现。不重开 VP-007 / VP-012；不引入业务域计费。

## 配置面与模块归属

- 走既有 **locale / Settings** 面（VP-007 交付的 Localization 设置与 `auto` 机制），**不是**新模块、不改 Profile 默认集。
- **缺省**：`zh-CN` 语义 + `auto` 时区；无任何配置仍能开发与快测。
- **生产 / 本 VP 验收**：显式配置（时区、数字、货币默认）后，格式合同可核对（快测 + UI 范例；`zh-CN` / `en-US` 至少各一）。
- **生效方式**：设置保存后随会话/进程生效；热加载不进退出分母。

## 首波冻结（退出分母 = Admin 格式语义）

| 能力 | 本 VP 交付 | 不进本 VP |
|------|------------|-----------|
| 时区 | 会话/用户级时区解析与展示（IANA / offset / `auto`）；时间输入与展示统一语义 | DB `timestamptz` 持久化合同（架构 RT-T03）；多时区排程 / 日历；热加载 |
| 数字 | locale 驱动千分位 / 小数位 / 百分比格式与输入解析合同 | 自定义每页格式引擎；业务域计量 |
| 货币 | ISO 4217 代码 / 符号 / 位置 / 小数位展示合同 | 汇率 / 换算 / 结算 / 计费（业务域）；钱包金额语义化（不重开 VP-011 演示面） |
| 设置面 | Settings Localization 下时区 / 数字 / 货币默认（归属与字段名由 lead Root R1 冻结，I-020-005） | 翻译中心 / 文案（VP-007，不重开）；改 Profile 默认集 |
| 内嵌默认 | 无配置 `zh-CN` 可运行；快测与 UI 范例可核对 | 强制全站统一 UTC / 固定时区启动；外部时区服务依赖 |

## 非目标

- **DB 时区持久化合同**（架构 RT-T03 仍 `registered`；如需可另立架构 VP，不在本波）
- **汇率 / 换算 / 计费 / 结算 / 发票**（业务域分支）
- **多时区排程 / 日历 / 提醒**；翻译与文案中心（VP-007 已交付）
- **热加载**；改 Profile 默认集 / 模块矩阵 / Manifest 装配语义
- 重开 VP-012；替代 VP-009 / VP-010；改变 Charter 边界；业务域页面

## 与相邻 VP 的边界

| VP / 分支 | 关系 |
|-----------|------|
| **VP-003** | 遵守薄内核。格式语义是 locale 工具层能力，不建平行认证或中央注册路径 |
| **VP-005** | 展示组件消费设计 token；语义与样式分离 |
| **VP-007** | 本 VP 消费其语言解析与 `auto` 运行时；不重开 locale / 设置骨架 |
| **VP-008 `go`** | Admin 类能力；激活前做 Admin 类 freshness。若实现改变 Profile 默认集 / Manifest 装配，按消费有效性暂挂 |
| **VP-009 / VP-010** | 格式相关安全（如时区炸弹类输入）与符合性 gap 归持续程序 |
| **架构 RT-T03** | DB `timestamptz` 仍 `registered`；本 VP 只做应用语义面，不越界 |
| **业务域** | 金额 / 计费语义不进本波；业务域成立后可消费本语义合同 |

## 方向级退出判据

1. 时区 / 数字 / 货币格式语义合同落盘并可核对（快测 + UI 范例；`zh-CN` / `en-US` 至少各一场景）。
2. `auto` 时区解析可用；显式配置后展示与输入语义一致（同一合同双向）。
3. 未引入汇率 / 计费 / DB 持久化时区合同；未改 Charter；未改 Profile 默认集作为本波成功条件。
4. 开放 required finding = 0（或已合法闭合）。

详细纲领阶段由 lead Root `GOAL-001-timezone-number-currency-formatting`（P-001）书写：R1 合同冻结（时区来源 / 设置归属 / 序列化落点）→ R2 时区语义 → R3 数字 / 货币语义 → R4 证据（快测 + 双 locale 范例）。本 VP 不写 Goal 五件套。

## 信息需求（P-005）

允许带未知立项。下列不影响「本 VP 意图已冻结」，但必须在对应阶段前关闭或经用户接受残余。

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 状态 |
|----|--------------|------|----------|----------|------|
| I-020-001 | 时区来源：会话级（客户端 / 请求探测）vs 用户级（存库）vs 两者；影响设置归属与 schema。 | required | 方案冻结 | R1 合同冻结 | **verified**（用户 2026-08-26 裁决：会话级 auto + 站点兜底 + 用户级 localStorage 覆盖 · Root D-002；合同 §2 L1～L4 实现） |
| I-020-002 | 数字 / 货币语义落点：仅前端展示 vs 后端 API 序列化也携带语义（如 decimal 字符串合同）。 | required | 方案冻结 | R1 合同冻结 | **verified**（用户 2026-08-26 裁决：前端落点；API 机器合同文档化 · Root D-002；§3 实现） |
| I-020-003 | 持久化时区合同（DB `timestamptz`）是否进本波？**本 VP 冻结为不进**（架构 RT-T03 仍 `registered`）。本行只作台账投影。 | required | 退出分母 | R1 | **registered**（VP 已冻结不进；Root D-001 投影） |
| I-020-004 | 汇率 / 换算是否进本波？**本 VP 冻结为不进**（业务域）。本行只作台账投影。 | required | 退出分母 | R1 | **registered**（VP 已冻结不进；Root D-001 投影） |
| I-020-005 | 默认时区 / 数字 / 货币的配置归属：Settings 哪一 tab、哪些字段。 | non-blocking | 方案冻结 | R2 | **verified**（用户 2026-08-26 确认：Localization tab 站点默认 + 头部 locale 通道 · Root D-002；实现于 GOAL-003/004） |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-020-timezone-number-currency-formatting | GOAL-001-timezone-number-currency-formatting | lead | 2026-08-26 | 唯一 delivery；Root active 0/4（R1～R4 待立项；激活审查 VRev-044 self `pass`） |

## 关门记录

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| 2026-08-27 | **closed v0.3.0**（用户书面确认「Root done 4/4 + VP-020 收尾」） | lead workspace-020 结项：Root done 4/4（R1 合同冻结 → R2 时区语义 → R3 数字/货币语义 → R4 证据与关门）；关门审计双腿 pass（Root A-001 self + A-002 grok build independent，0 required）；全量回归（Go 全绿 / web 1181）；退出判据 1–4 全部满足；VRev-045（self）关门审查 pass | 工作区 goal-tree / Root 03-audit A-001/A-002 / GOAL-005 `attachments/r4-evidence-matrix.md` | 无开放必改；书面接受残余：分组位序容差、`defaultCurrency` 句法三字母（非完整 ISO 目录）、业务金额展示不接线（钱包演示面不重开）；RT-T03 保持 registered |

## 规划修订短史

| date | change |
|------|--------|
| 2026-08-26 | 初创 `planned`：用户确认立项（Admin 功能分支基架能力剩余 #5 · 时区/数字/货币格式语义）；退出分母 = Admin 格式语义面；DB 持久化时区合同（RT-T03）、汇率/计费不进分母。roadmap 索引原子同步 |
| 2026-08-26 | v0.2.0 **激活**：VRev-044 self `pass`（V-F079/V-F080 → 激活事务内 fixed）；Admin 类 freshness PASS（`66f5fd1f` → `c6fda691`，不暂挂 `go`）；lead = `workspace-020-timezone-number-currency-formatting`；Root 五件套与 ledger 目录建立（跨入口调用 govern 原语，用户 2026-08-26 明确指令） |
| 2026-08-26→27 | **交付与关门**：R1 合同冻结（用户裁决 I-001/I-002/I-005，D-002）→ R2 时区语义（GOAL-003）→ R3 数字/货币语义（GOAL-004 · grok 独立审 fail→fixed→复审 pass；migration v62 `default_currency`）→ R4 证据与关门（GOAL-005 · 证据矩阵/越界核账/F-007 加严；Root 关门审计 A-001 self + A-002 grok independent 双 pass）→ **v0.3.0 `closed`**（2026-08-27 用户书面确认）；VRev-045 关门审查；I-020-001/002/005 回写 verified |